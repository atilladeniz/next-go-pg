import { headers } from "next/headers"
import { auth } from "./auth"

/**
 * Gets the current Better Auth session.
 * Use this in Server Components for protected pages.
 *
 * Note: Go JWT token refresh is handled by middleware or API routes,
 * not in Server Components (Next.js 16 restricts cookie modification).
 */
export async function getSession() {
	const session = await auth.api.getSession({
		headers: await headers(),
	})

	return session
}
