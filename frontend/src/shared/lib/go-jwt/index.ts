import { jwtVerify, SignJWT } from "jose"

// JWT configuration for Go backend communication
const GO_JWT_SECRET = new TextEncoder().encode(
	process.env.GO_JWT_SECRET || "dev-secret-change-in-production",
)
const GO_JWT_ALGORITHM = "HS256"
const GO_JWT_EXPIRATION = "15m" // Short-lived tokens, refresh on each request

export interface GoJWTPayload {
	sub: string // User ID
	email: string
	name?: string
	sid: string // Session ID
	iat?: number
	exp?: number
}

/**
 * Creates a signed JWT for Go backend authentication
 * This token is separate from Better Auth's session and allows
 * Go to validate requests without calling back to Next.js
 */
export async function createGoToken(payload: Omit<GoJWTPayload, "iat" | "exp">): Promise<string> {
	return new SignJWT({
		sub: payload.sub,
		email: payload.email,
		name: payload.name,
		sid: payload.sid,
	})
		.setProtectedHeader({ alg: GO_JWT_ALGORITHM })
		.setIssuedAt()
		.setExpirationTime(GO_JWT_EXPIRATION)
		.sign(GO_JWT_SECRET)
}

/**
 * Verifies a Go JWT token (useful for testing)
 */
export async function verifyGoToken(token: string): Promise<GoJWTPayload | null> {
	try {
		const { payload } = await jwtVerify(token, GO_JWT_SECRET)
		return payload as unknown as GoJWTPayload
	} catch {
		return null
	}
}

/**
 * Cookie name for the Go auth token
 */
export const GO_AUTH_COOKIE = "go-auth-token"

/**
 * Cookie options for the Go auth token
 */
export const GO_AUTH_COOKIE_OPTIONS = {
	httpOnly: true,
	secure: process.env.NODE_ENV === "production",
	sameSite: "lax" as const,
	path: "/",
	maxAge: 60 * 15, // 15 minutes (matches JWT expiration)
}
