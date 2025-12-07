"use client"

import { Button } from "@shared/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@shared/ui/card"

interface EmailSentCardProps {
	email: string
	onReset: () => void
}

export function EmailSentCard({ email, onReset }: EmailSentCardProps) {
	return (
		<Card className="w-full max-w-md">
			<CardHeader className="text-center">
				<CardTitle className="text-2xl">E-Mail gesendet</CardTitle>
				<CardDescription>
					Wir haben eine E-Mail an <strong>{email}</strong> gesendet.
				</CardDescription>
			</CardHeader>
			<CardContent className="space-y-4">
				<p className="text-center text-sm text-muted-foreground">
					Prüfe dein Postfach und klicke auf den Link.
				</p>
				<Button variant="outline" className="w-full" onClick={onReset}>
					Andere E-Mail verwenden
				</Button>
			</CardContent>
		</Card>
	)
}
