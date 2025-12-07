import { passkeyClient } from "@better-auth/passkey/client"
import { magicLinkClient, twoFactorClient } from "better-auth/client/plugins"
import { createAuthClient } from "better-auth/react"

export const authClient = createAuthClient({
	baseURL: process.env.NEXT_PUBLIC_APP_URL || "http://localhost:3000",
	plugins: [
		magicLinkClient(),
		twoFactorClient({
			onTwoFactorRedirect() {
				window.location.href = "/2fa"
			},
		}),
		passkeyClient(),
	],
})

export const { signIn, signUp, signOut, useSession, sendVerificationEmail, twoFactor, passkey } =
	authClient
