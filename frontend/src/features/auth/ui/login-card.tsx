"use client"

import { signIn } from "@shared/lib/auth-client"
import { Button } from "@shared/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@shared/ui/card"
import { Input } from "@shared/ui/input"
import { Label } from "@shared/ui/label"
import { Key, Loader2 } from "lucide-react"
import { useRouter } from "next/navigation"
import { useState } from "react"

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
	const router = useRouter()
	const [passkeyLoading, setPasskeyLoading] = useState(false)
	const [passkeyError, setPasskeyError] = useState<string | null>(null)
	const isDisabled = loading || retryAfter > 0 || passkeyLoading

	const getButtonText = () => {
		if (loading) return "Wird gesendet..."
		if (retryAfter > 0) return `Warte ${retryAfter}s...`
		return "Weiter"
	}

	const handlePasskeyLogin = async () => {
		setPasskeyLoading(true)
		setPasskeyError(null)

		try {
			const result = await signIn.passkey()

			if (result?.error) {
				setPasskeyError(result.error.message || "Passkey-Anmeldung fehlgeschlagen")
				setPasskeyLoading(false)
				return
			}

			router.push("/dashboard")
		} catch (err) {
			const errorMessage = err instanceof Error ? err.message : "Passkey-Anmeldung fehlgeschlagen"
			setPasskeyError(errorMessage)
			setPasskeyLoading(false)
		}
	}

	const displayError = error || passkeyError

	return (
		<Card className="w-full max-w-md">
			<CardHeader className="text-center">
				<CardTitle className="text-2xl">Anmelden</CardTitle>
				<CardDescription>Gib deine E-Mail-Adresse ein, um fortzufahren.</CardDescription>
			</CardHeader>
			<CardContent className="space-y-4">
				{displayError && (
					<div className="rounded-lg bg-destructive/10 p-3 text-sm text-destructive">
						{displayError}
					</div>
				)}

				{/* Passkey Login Button */}
				<Button
					type="button"
					variant="outline"
					className="w-full"
					onClick={handlePasskeyLogin}
					disabled={isDisabled}
				>
					{passkeyLoading ? (
						<Loader2 className="mr-2 h-4 w-4 animate-spin" />
					) : (
						<Key className="mr-2 h-4 w-4" />
					)}
					Mit Passkey anmelden
				</Button>

				<div className="relative">
					<div className="absolute inset-0 flex items-center">
						<span className="w-full border-t" />
					</div>
					<div className="relative flex justify-center text-xs uppercase">
						<span className="bg-background px-2 text-muted-foreground">oder</span>
					</div>
				</div>

				{/* Magic Link Form */}
				<form onSubmit={onSubmit} className="space-y-4">
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
							autoComplete="username webauthn"
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
