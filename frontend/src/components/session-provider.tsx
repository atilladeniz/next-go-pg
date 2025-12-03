"use client"

import { createContext, useContext } from "react"

type User = {
	id: string
	name: string
	email: string
	emailVerified: boolean
	image?: string | null
	createdAt: Date
	updatedAt: Date
}

type Session = {
	user: User
	session: {
		id: string
		userId: string
		expiresAt: Date
	}
}

const SessionContext = createContext<Session | null>(null)

export function SessionProvider({
	session,
	children,
}: {
	session: Session
	children: React.ReactNode
}) {
	return <SessionContext value={session}>{children}</SessionContext>
}

export function useServerSession() {
	const session = useContext(SessionContext)
	if (!session) {
		throw new Error("useServerSession must be used within SessionProvider")
	}
	return session
}
