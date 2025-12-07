"use client"

import { authClient } from "@shared/lib/auth-client"
import { useCallback, useEffect, useState } from "react"

export interface SessionInfo {
	id: string
	token: string
	userId: string
	expiresAt: Date
	ipAddress?: string | null
	userAgent?: string | null
	createdAt: Date
	updatedAt: Date
	isCurrent: boolean
}

export function useSessions() {
	const [sessions, setSessions] = useState<SessionInfo[]>([])
	const [loading, setLoading] = useState(true)
	const [error, setError] = useState<string | null>(null)
	const [revoking, setRevoking] = useState<string | null>(null)

	const fetchSessions = useCallback(async () => {
		setLoading(true)
		setError(null)

		try {
			const result = await authClient.listSessions()

			if (result.error) {
				setError(result.error.message || "Fehler beim Laden der Sessions")
				return
			}

			const currentSession = await authClient.getSession()
			const currentToken = currentSession.data?.session?.token

			const sessionsWithCurrent = (result.data || []).map((session) => ({
				...session,
				expiresAt: new Date(session.expiresAt),
				createdAt: new Date(session.createdAt),
				updatedAt: new Date(session.updatedAt),
				isCurrent: session.token === currentToken,
			}))

			setSessions(sessionsWithCurrent)
		} catch {
			setError("Fehler beim Laden der Sessions")
		} finally {
			setLoading(false)
		}
	}, [])

	const revokeSession = useCallback(
		async (token: string) => {
			setRevoking(token)

			try {
				const result = await authClient.revokeSession({ token })

				if (result.error) {
					setError(result.error.message || "Fehler beim Beenden der Session")
					return false
				}

				await fetchSessions()
				return true
			} catch {
				setError("Fehler beim Beenden der Session")
				return false
			} finally {
				setRevoking(null)
			}
		},
		[fetchSessions],
	)

	const revokeOtherSessions = useCallback(async () => {
		setLoading(true)

		try {
			const result = await authClient.revokeOtherSessions()

			if (result.error) {
				setError(result.error.message || "Fehler beim Beenden der Sessions")
				return false
			}

			await fetchSessions()
			return true
		} catch {
			setError("Fehler beim Beenden der Sessions")
			return false
		} finally {
			setLoading(false)
		}
	}, [fetchSessions])

	useEffect(() => {
		fetchSessions()
	}, [fetchSessions])

	return {
		sessions,
		loading,
		error,
		revoking,
		revokeSession,
		revokeOtherSessions,
		refresh: fetchSessions,
	}
}
