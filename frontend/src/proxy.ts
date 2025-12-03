import { type NextRequest, NextResponse } from "next/server"

const protectedRoutes = ["/dashboard"]
const authRoutes = ["/login", "/register"]

export function proxy(request: NextRequest) {
	const { pathname } = request.nextUrl
	const sessionCookie = request.cookies.get("better-auth.session_token")

	const isProtectedRoute = protectedRoutes.some((route) => pathname.startsWith(route))
	const isAuthRoute = authRoutes.some((route) => pathname.startsWith(route))

	// Redirect to login if accessing protected route without session
	if (isProtectedRoute && !sessionCookie) {
		const loginUrl = new URL("/login", request.url)
		loginUrl.searchParams.set("callbackUrl", pathname)
		return NextResponse.redirect(loginUrl)
	}

	// Redirect to dashboard if accessing auth routes with active session
	if (isAuthRoute && sessionCookie) {
		return NextResponse.redirect(new URL("/dashboard", request.url))
	}

	return NextResponse.next()
}

export const config = {
	matcher: ["/dashboard/:path*", "/login", "/register"],
}
