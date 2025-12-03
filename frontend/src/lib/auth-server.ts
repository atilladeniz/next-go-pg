import { cookies, headers } from "next/headers"
import { auth } from "./auth"

export async function getSession() {
	const session = await auth.api.getSession({
		headers: await headers(),
	})
	return session
}

export async function getStats() {
	const cookieStore = await cookies()
	const cookieHeader = cookieStore
		.getAll()
		.map((c) => `${c.name}=${c.value}`)
		.join("; ")

	try {
		const res = await fetch("http://localhost:8080/api/v1/stats", {
			headers: { Cookie: cookieHeader },
			cache: "no-store",
		})
		if (!res.ok) return null
		return res.json()
	} catch {
		return null
	}
}
