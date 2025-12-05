// Public API for shared/lib/auth-server
// Server-only: use in Server Components, API routes, middleware
// DO NOT import in "use client" components!
export { auth } from "./auth"
export { getSession } from "./auth-server"
