"use server"

import { cookies } from "next/headers"
import { getSession } from "../auth-server"
import { createGoToken, GO_AUTH_COOKIE, GO_AUTH_COOKIE_OPTIONS } from "./index"

/**
 * Refreshes the Go JWT token based on the current Better Auth session.
 * This should be called:
 * 1. After successful login (in the auth flow)
 * 2. On page load for protected pages (to ensure fresh token)
 * 3. Before making direct API calls to Go backend
 *
 * Returns the token if successful, null if no session exists.
 */
export async function refreshGoToken(): Promise<string | null> {
	const session = await getSession()

	if (!session?.user || !session?.session) {
		// No session, clear the Go token cookie
		const cookieStore = await cookies()
		cookieStore.delete(GO_AUTH_COOKIE)
		return null
	}

	// Create a fresh JWT for Go
	const token = await createGoToken({
		sub: session.user.id,
		email: session.user.email,
		name: session.user.name || undefined,
		sid: session.session.id,
	})

	// Set the cookie
	const cookieStore = await cookies()
	cookieStore.set(GO_AUTH_COOKIE, token, GO_AUTH_COOKIE_OPTIONS)

	return token
}

/**
 * Clears the Go JWT token (call on logout)
 */
export async function clearGoToken(): Promise<void> {
	const cookieStore = await cookies()
	cookieStore.delete(GO_AUTH_COOKIE)
}

/**
 * Gets the current Go JWT token from cookies (if exists)
 */
export async function getGoToken(): Promise<string | null> {
	const cookieStore = await cookies()
	return cookieStore.get(GO_AUTH_COOKIE)?.value || null
}
