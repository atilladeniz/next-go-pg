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
	plugins: [
		magicLink({
			disableSignUp: false,
			expiresIn: 60 * 10, // 10 minutes
			sendMagicLink: async ({ email, url }) => {
				await sendMail(
					email,
					"Dein Anmelde-Link",
					`
					<h1>Anmeldung</h1>
					<p>Klicke auf den folgenden Link, um dich anzumelden:</p>
					<a href="${url}">${url}</a>
					<p>Der Link ist 10 Minuten gültig.</p>
					`,
				)
			},
		}),
	],
})
