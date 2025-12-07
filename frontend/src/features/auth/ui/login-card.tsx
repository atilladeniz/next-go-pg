"use client"

import { Button } from "@shared/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@shared/ui/card"
import { Input } from "@shared/ui/input"
import { Label } from "@shared/ui/label"

interface LoginCardProps {
	email: string
	error: string
	loading: boolean
	retryAfter: number
	onEmailChange: (email: string) => void
	onSubmit: (e: React.FormEvent) => void
}

export function LoginCard({
	email,
	error,
	loading,
	retryAfter,
	onEmailChange,
	onSubmit,
}: LoginCardProps) {
	const isDisabled = loading || retryAfter > 0

	const getButtonText = () => {
		if (loading) return "Wird gesendet..."
		if (retryAfter > 0) return `Warte ${retryAfter}s...`
		return "Weiter"
	}

	return (
		<Card className="w-full max-w-md">
			<CardHeader className="text-center">
				<CardTitle className="text-2xl">Anmelden</CardTitle>
				<CardDescription>Gib deine E-Mail-Adresse ein, um fortzufahren.</CardDescription>
			</CardHeader>
			<CardContent>
				<form onSubmit={onSubmit} className="space-y-4">
					{error && (
						<div className="rounded-lg bg-destructive/10 p-3 text-sm text-destructive">{error}</div>
					)}

					<div className="space-y-2">
						<Label htmlFor="email">E-Mail</Label>
						<Input
							id="email"
							type="email"
							value={email}
							onChange={(e) => onEmailChange(e.target.value)}
							required
							placeholder="name@example.com"
							disabled={isDisabled}
						/>
					</div>

					<Button type="submit" className="w-full" disabled={isDisabled}>
						{getButtonText()}
					</Button>

					<p className="text-center text-xs text-muted-foreground">
						Du erhältst einen Link per E-Mail. Neue Nutzer werden automatisch registriert.
					</p>
				</form>
			</CardContent>
		</Card>
	)
}
