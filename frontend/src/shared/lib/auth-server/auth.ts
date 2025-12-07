import { passkey } from "@better-auth/passkey"
import { betterAuth } from "better-auth"
import { magicLink, twoFactor } from "better-auth/plugins"
import { Pool } from "pg"

// Session configuration constants
const SESSION_EXPIRY_SECONDS = 60 * 60 * 24 * 7 // 7 days
const SESSION_UPDATE_AGE_SECONDS = 60 * 60 * 24 // 1 day
const COOKIE_CACHE_MAX_AGE_SECONDS = 60 * 5 // 5 minutes
const MAGIC_LINK_EXPIRY_SECONDS = 60 * 10 // 10 minutes

// Rate limiting constants
const RATE_LIMIT_WINDOW_SECONDS = 60
const RATE_LIMIT_MAX_REQUESTS = 100
const MAGIC_LINK_RATE_LIMIT_MAX = 3
const VERIFICATION_EMAIL_RATE_LIMIT_MAX = 3

const apiUrl = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080"
const webhookSecret = process.env.WEBHOOK_SECRET || ""

// Critical webhooks that must succeed (user won't receive email otherwise)
const CRITICAL_WEBHOOKS = ["send-magic-link", "send-verification-email", "send-2fa-otp"]

// Helper to call backend webhooks
const callWebhook = async (endpoint: string, data: Record<string, string>) => {
	const isCritical = CRITICAL_WEBHOOKS.includes(endpoint)

	try {
		const response = await fetch(`${apiUrl}/api/v1/webhooks/${endpoint}`, {
			method: "POST",
			headers: {
				"Content-Type": "application/json",
				"X-Webhook-Secret": webhookSecret,
			},
			body: JSON.stringify(data),
		})

		if (!response.ok && isCritical) {
			// Critical webhooks must succeed - throw error to surface to user
			throw new Error(`Email konnte nicht gesendet werden (${response.status})`)
		}
	} catch (error) {
		if (isCritical) {
			// Re-throw critical errors so Better Auth can surface them to the user
			throw error
		}
		// Non-critical webhooks (notifications) can fail silently
		// They are logged by the backend
	}
}

export const auth = betterAuth({
	database: new Pool({
		connectionString: process.env.DATABASE_URL,
	}),
	session: {
		expiresIn: SESSION_EXPIRY_SECONDS,
		updateAge: SESSION_UPDATE_AGE_SECONDS,
		cookieCache: {
			enabled: true,
			maxAge: COOKIE_CACHE_MAX_AGE_SECONDS,
		},
	},
	trustedOrigins: [process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080"],
	emailVerification: {
		sendOnSignUp: true,
		autoSignInAfterVerification: false,
		sendVerificationEmail: async ({ user, url }) => {
			// Rewrite URL from /api/auth/verify-email to /verify-email
			const verifyUrl = url.replace("/api/auth/verify-email", "/verify-email")
			await callWebhook("send-verification-email", {
				email: user.email,
				name: user.name || "",
				url: verifyUrl,
			})
		},
	},
	rateLimit: {
		enabled: true,
		window: RATE_LIMIT_WINDOW_SECONDS,
		max: RATE_LIMIT_MAX_REQUESTS,
		storage: "database",
		customRules: {
			"/sign-in/magic-link": {
				window: RATE_LIMIT_WINDOW_SECONDS,
				max: MAGIC_LINK_RATE_LIMIT_MAX,
			},
			"/send-verification-email": {
				window: RATE_LIMIT_WINDOW_SECONDS,
				max: VERIFICATION_EMAIL_RATE_LIMIT_MAX,
			},
		},
	},
	databaseHooks: {
		session: {
			create: {
				after: async (session) => {
					// Delegate to backend - clean separation of concerns
					await callWebhook("session-created", {
						sessionId: session.id,
						userId: session.userId,
						userAgent: session.userAgent || "",
						ipAddress: session.ipAddress || "",
					})
				},
			},
		},
	},
	plugins: [
		magicLink({
			disableSignUp: false,
			expiresIn: MAGIC_LINK_EXPIRY_SECONDS,
			sendMagicLink: async ({ email, url }) => {
				// Rewrite URL from /api/auth/magic-link/verify to /magic-link/verify
				const verifyUrl = url.replace("/api/auth/magic-link/verify", "/magic-link/verify")
				await callWebhook("send-magic-link", {
					email,
					url: verifyUrl,
				})
			},
		}),
		twoFactor({
			issuer: process.env.NEXT_PUBLIC_APP_NAME || "Next-Go-PG",
			skipVerificationOnEnable: true, // No password for magic link users
			otpOptions: {
				async sendOTP({ user, otp }) {
					await callWebhook("send-2fa-otp", {
						email: user.email,
						name: user.name || "",
						otp,
					})
				},
			},
		}),
		passkey({
			rpID: process.env.PASSKEY_RP_ID || "localhost",
			rpName: process.env.NEXT_PUBLIC_APP_NAME || "Next-Go-PG",
			origin: process.env.NEXT_PUBLIC_APP_URL || "http://localhost:3000",
		}),
	],
})
