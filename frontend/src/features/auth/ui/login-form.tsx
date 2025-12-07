"use client"

import { signIn } from "@shared/lib/auth-client"
import { Button } from "@shared/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@shared/ui/card"
import { Input } from "@shared/ui/input"
import { Label } from "@shared/ui/label"
import { useState } from "react"

export function LoginForm() {
	const [email, setEmail] = useState("")
	const [error, setError] = useState("")
	const [loading, setLoading] = useState(false)
	const [sent, setSent] = useState(false)

	const handleSubmit = async (e: React.FormEvent) => {
		e.preventDefault()
		setError("")
		setLoading(true)

		const result = await signIn.magicLink({
			email,
			callbackURL: "/dashboard",
		})

		if (result.error) {
			setError(result.error.message || "Anmeldung fehlgeschlagen")
			setLoading(false)
			return
		}

		setSent(true)
		setLoading(false)
	}

	if (sent) {
		return (
			<div className="flex min-h-screen items-center justify-center bg-background">
				<Card className="w-full max-w-md">
					<CardHeader className="text-center">
						<CardTitle className="text-2xl">Link gesendet</CardTitle>
						<CardDescription>
							Wir haben einen Anmelde-Link an <strong>{email}</strong> gesendet.
						</CardDescription>
					</CardHeader>
					<CardContent className="space-y-4">
						<p className="text-center text-sm text-muted-foreground">
							Prüfe dein Postfach und klicke auf den Link, um dich anzumelden.
						</p>
						<Button
							variant="outline"
							className="w-full"
							onClick={() => {
								setSent(false)
								setEmail("")
							}}
						>
							Andere E-Mail verwenden
						</Button>
					</CardContent>
				</Card>
			</div>
		)
	}

	return (
		<div className="flex min-h-screen items-center justify-center bg-background">
			<Card className="w-full max-w-md">
				<CardHeader className="text-center">
					<CardTitle className="text-2xl">Anmelden</CardTitle>
					<CardDescription>
						Gib deine E-Mail-Adresse ein, um einen Anmelde-Link zu erhalten.
					</CardDescription>
				</CardHeader>
				<CardContent>
					<form onSubmit={handleSubmit} className="space-y-4">
						{error && (
							<div className="rounded-lg bg-destructive/10 p-3 text-sm text-destructive">
								{error}
							</div>
						)}

						<div className="space-y-2">
							<Label htmlFor="email">E-Mail</Label>
							<Input
								id="email"
								type="email"
								value={email}
								onChange={(e) => setEmail(e.target.value)}
								required
								placeholder="name@example.com"
							/>
						</div>

						<Button type="submit" className="w-full" disabled={loading}>
							{loading ? "Wird gesendet..." : "Anmelde-Link senden"}
						</Button>
					</form>
				</CardContent>
			</Card>
		</div>
	)
}
