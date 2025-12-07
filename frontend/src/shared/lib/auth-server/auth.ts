import { betterAuth } from "better-auth"
import { magicLink } from "better-auth/plugins"
import nodemailer from "nodemailer"
import { Pool } from "pg"

const transporter = nodemailer.createTransport({
	host: process.env.SMTP_HOST || "127.0.0.1",
	port: Number(process.env.SMTP_PORT) || 1025,
	secure: false,
})

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
	plugins: [
		magicLink({
			sendMagicLink: async ({ email, url }) => {
				await transporter.sendMail({
					from: process.env.SMTP_FROM || "noreply@localhost",
					to: email,
					subject: "Dein Anmelde-Link",
					html: `
						<h1>Anmeldung</h1>
						<p>Klicke auf den folgenden Link, um dich anzumelden:</p>
						<a href="${url}">${url}</a>
						<p>Der Link ist 15 Minuten gültig.</p>
					`,
				})
			},
		}),
	],
})
