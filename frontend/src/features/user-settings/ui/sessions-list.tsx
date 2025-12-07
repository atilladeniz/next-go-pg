"use client"

import { Button } from "@shared/ui/button"
import { useSessions } from "../model/use-sessions"
import { SessionCard } from "./session-card"

export function SessionsList() {
	const { sessions, loading, error, revoking, revokeSession, revokeOtherSessions } = useSessions()

	if (loading) {
		return (
			<div className="space-y-3">
				{[1, 2].map((i) => (
					<div key={i} className="h-24 animate-pulse rounded-lg bg-muted" />
				))}
			</div>
		)
	}

	if (error) {
		return <div className="rounded-lg bg-destructive/10 p-4 text-sm text-destructive">{error}</div>
	}

	const otherSessions = sessions.filter((s) => !s.isCurrent)

	return (
		<div className="space-y-4">
			<div className="space-y-3">
				{sessions.map((session) => (
					<SessionCard
						key={session.id}
						session={session}
						onRevoke={revokeSession}
						isRevoking={revoking === session.token}
					/>
				))}
			</div>

			{otherSessions.length > 0 && (
				<Button variant="outline" className="w-full" onClick={revokeOtherSessions}>
					Alle anderen Sessions beenden ({otherSessions.length})
				</Button>
			)}
		</div>
	)
}
