import { betterAuth } from "better-auth"
import { magicLink } from "better-auth/plugins"
import nodemailer from "nodemailer"
import { Pool } from "pg"

const transporter = nodemailer.createTransport({
	host: process.env.SMTP_HOST || "127.0.0.1",
	port: Number(process.env.SMTP_PORT) || 1025,
	secure: false,
})

const sendMail = async (to: string, subject: string, html: string) => {
	await transporter.sendMail({
		from: process.env.SMTP_FROM || "noreply@localhost",
		to,
		subject,
		html,
	})
}

function parseUserAgent(ua: string | null): string {
	if (!ua) return "Unbekanntes Gerät"
	let browser = "Browser"
	let os = "System"
	if (ua.includes("Firefox")) browser = "Firefox"
	else if (ua.includes("Edg/")) browser = "Edge"
	else if (ua.includes("Chrome")) browser = "Chrome"
	else if (ua.includes("Safari")) browser = "Safari"
	if (ua.includes("Windows")) os = "Windows"
	else if (ua.includes("Mac OS")) os = "macOS"
	else if (ua.includes("Linux")) os = "Linux"
	else if (ua.includes("Android")) os = "Android"
	else if (ua.includes("iPhone") || ua.includes("iPad")) os = "iOS"
	return `${browser} auf ${os}`
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
			const verifyUrl = url.replace("/api/auth/verify-email", "/verify-email")
			await sendMail(
				user.email,
				"E-Mail bestätigen",
				`
				<h1>Willkommen!</h1>
				<p>Klicke auf den folgenden Link, um deine E-Mail-Adresse zu bestätigen:</p>
				<a href="${verifyUrl}">${verifyUrl}</a>
				<p>Der Link ist 24 Stunden gültig.</p>
				`,
			)
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
					const db = new Pool({ connectionString: process.env.DATABASE_URL })
					try {
						const result = await db.query('SELECT email, name FROM "user" WHERE id = $1', [
							session.userId,
						])
						const user = result.rows[0]
						if (!user) return

						const device = parseUserAgent(session.userAgent ?? null)
						const time = new Date().toLocaleString("de-DE", {
							dateStyle: "medium",
							timeStyle: "short",
						})

						const appUrl = process.env.NEXT_PUBLIC_APP_URL || "http://localhost:3000"
						const settingsUrl = `${appUrl}/settings`

						await sendMail(
							user.email,
							"Neue Anmeldung erkannt",
							`
							<h1>Neue Anmeldung in deinem Konto</h1>
							<p>Hallo ${user.name || ""},</p>
							<p>Wir haben eine neue Anmeldung in deinem Konto festgestellt:</p>
							<ul>
								<li><strong>Gerät:</strong> ${device}</li>
								<li><strong>IP-Adresse:</strong> ${session.ipAddress || "Unbekannt"}</li>
								<li><strong>Zeit:</strong> ${time}</li>
							</ul>
							<p>Wenn du das nicht warst, überprüfe bitte sofort deine aktiven Sessions und beende verdächtige Sitzungen:</p>
							<p><a href="${settingsUrl}" style="display: inline-block; padding: 12px 24px; background-color: #000; color: #fff; text-decoration: none; border-radius: 6px;">Sessions verwalten</a></p>
							<p style="margin-top: 16px; font-size: 14px; color: #666;">
								Oder kopiere diesen Link: <a href="${settingsUrl}">${settingsUrl}</a>
							</p>
							`,
						)
					} finally {
						await db.end()
					}
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
				const appUrl = process.env.NEXT_PUBLIC_APP_URL || "http://localhost:3000"

				await sendMail(
					email,
					"Dein Anmelde-Link",
					`
					<h1>Anmeldung</h1>
					<p>Klicke auf den Button, um dich anzumelden:</p>
					<p style="margin: 24px 0;">
						<a href="${verifyUrl}" style="display: inline-block; padding: 12px 24px; background-color: #000; color: #fff; text-decoration: none; border-radius: 6px;">Jetzt anmelden</a>
					</p>
					<p style="font-size: 14px; color: #666;">
						Oder kopiere diesen Link: <a href="${verifyUrl}">${verifyUrl}</a>
					</p>
					<p style="font-size: 14px; color: #666;">Der Link ist 10 Minuten gültig.</p>
					`,
				)
			},
		}),
	],
})
