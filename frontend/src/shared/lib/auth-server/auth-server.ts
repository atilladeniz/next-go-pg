import { cookies, headers } from "next/headers"
import { createGoToken, GO_AUTH_COOKIE, GO_AUTH_COOKIE_OPTIONS } from "../go-jwt"
import { auth } from "./auth"

/**
 * Gets the current Better Auth session and refreshes the Go JWT token.
 * Use this in Server Components for protected pages.
 */
export async function getSession() {
	const session = await auth.api.getSession({
		headers: await headers(),
	})

	// If session exists, refresh the Go JWT token
	if (session?.user && session?.session) {
		const token = await createGoToken({
			sub: session.user.id,
			email: session.user.email,
			name: session.user.name || undefined,
			sid: session.session.id,
		})

		const cookieStore = await cookies()
		cookieStore.set(GO_AUTH_COOKIE, token, GO_AUTH_COOKIE_OPTIONS)
	}

	return session
}
