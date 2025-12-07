"use client"

import { formatLocationFromIP, parseUserAgent } from "@shared/lib/geo"
import { Button } from "@shared/ui/button"
import { Card, CardContent } from "@shared/ui/card"
import { Globe, Laptop, Monitor, Smartphone, Tablet } from "lucide-react"
import type { SessionInfo } from "../model/use-sessions"

interface SessionCardProps {
	session: SessionInfo
	onRevoke: (token: string) => void
	isRevoking: boolean
}

function DeviceIcon({ device }: { device: string }) {
	switch (device) {
		case "Mobile":
			return <Smartphone className="h-5 w-5" />
		case "Tablet":
			return <Tablet className="h-5 w-5" />
		case "Desktop":
			return <Monitor className="h-5 w-5" />
		default:
			return <Laptop className="h-5 w-5" />
	}
}

export function SessionCard({ session, onRevoke, isRevoking }: SessionCardProps) {
	const { browser, os, device } = parseUserAgent(session.userAgent ?? null)
	const locationText = formatLocationFromIP(session.ipAddress ?? null)

	const formatDate = (date: Date) => {
		return new Intl.DateTimeFormat("de-DE", {
			day: "2-digit",
			month: "2-digit",
			year: "numeric",
			hour: "2-digit",
			minute: "2-digit",
		}).format(date)
	}

	return (
		<Card className={session.isCurrent ? "border-primary" : ""}>
			<CardContent className="p-4">
				<div className="flex items-start justify-between gap-4">
					<div className="flex items-start gap-3">
						<div className="mt-1 rounded-lg bg-muted p-2">
							<DeviceIcon device={device} />
						</div>
						<div className="space-y-1">
							<div className="flex items-center gap-2">
								<span className="font-medium">
									{browser} auf {os}
								</span>
								{session.isCurrent && (
									<span className="rounded-full bg-primary px-2 py-0.5 text-xs text-primary-foreground">
										Aktuelle Session
									</span>
								)}
							</div>
							<div className="flex items-center gap-4 text-sm text-muted-foreground">
								<span className="flex items-center gap-1">
									<Globe className="h-3 w-3" />
									{locationText}
								</span>
								<span>IP: {session.ipAddress || "Unbekannt"}</span>
							</div>
							<div className="text-xs text-muted-foreground">
								Angemeldet: {formatDate(session.createdAt)}
							</div>
						</div>
					</div>
					{!session.isCurrent && (
						<Button
							variant="outline"
							size="sm"
							onClick={() => onRevoke(session.token)}
							disabled={isRevoking}
						>
							{isRevoking ? "..." : "Beenden"}
						</Button>
					)}
				</div>
			</CardContent>
		</Card>
	)
}
