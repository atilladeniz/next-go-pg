import { betterAuth } from "better-auth"
import { magicLink } from "better-auth/plugins"
import { Pool } from "pg"

const apiUrl = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080"
const webhookSecret = process.env.WEBHOOK_SECRET || ""

// Helper to call backend webhooks
const callWebhook = async (endpoint: string, data: Record<string, string>) => {
	try {
		await fetch(`${apiUrl}/api/v1/webhooks/${endpoint}`, {
			method: "POST",
			headers: {
				"Content-Type": "application/json",
				"X-Webhook-Secret": webhookSecret,
			},
			body: JSON.stringify(data),
		})
	} catch (error) {
		console.error(`Failed to call webhook ${endpoint}:`, error)
	}
}

export const auth = betterAuth({
	database: new Pool({
		connectionString: process.env.DATABASE_URL,
	}),
	session: {
		expiresIn: 60 * 60 * 24 * 7, // 7 days
		updateAge: 60 * 60 * 24, // 1 day
		cookieCache: {
			enabled: true,
			maxAge: 60 * 5, // 5 minutes
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
		window: 60,
		max: 100,
		storage: "database",
		customRules: {
			"/sign-in/magic-link": {
				window: 60,
				max: 3,
			},
			"/send-verification-email": {
				window: 60,
				max: 3,
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
			expiresIn: 60 * 10, // 10 minutes
			sendMagicLink: async ({ email, url }) => {
				// Rewrite URL from /api/auth/magic-link/verify to /magic-link/verify
				const verifyUrl = url.replace("/api/auth/magic-link/verify", "/magic-link/verify")
				await callWebhook("send-magic-link", {
					email,
					url: verifyUrl,
				})
			},
		}),
	],
})
