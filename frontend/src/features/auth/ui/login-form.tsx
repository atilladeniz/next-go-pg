"use client"

import { signIn } from "@shared/lib/auth-client"
import { Button } from "@shared/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@shared/ui/card"
import { Input } from "@shared/ui/input"
import { Label } from "@shared/ui/label"
import Link from "next/link"
import { useRouter } from "next/navigation"
import { useState } from "react"

export function LoginForm() {
	const router = useRouter()
	const [email, setEmail] = useState("")
	const [password, setPassword] = useState("")
	const [error, setError] = useState("")
	const [loading, setLoading] = useState(false)

	const handleSubmit = async (e: React.FormEvent) => {
		e.preventDefault()
		setError("")
		setLoading(true)

		const result = await signIn.email({
			email,
			password,
		})

		if (result.error) {
			setError(result.error.message || "Login fehlgeschlagen")
			setLoading(false)
			return
		}

		router.push("/dashboard")
	}

	return (
		<div className="flex min-h-screen items-center justify-center bg-background">
			<Card className="w-full max-w-md">
				<CardHeader className="text-center">
					<CardTitle className="text-2xl">Anmelden</CardTitle>
					<CardDescription>
						Noch kein Konto?{" "}
						<Link href="/register" className="text-primary hover:underline">
							Registrieren
						</Link>
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

						<div className="space-y-2">
							<Label htmlFor="password">Passwort</Label>
							<Input
								id="password"
								type="password"
								value={password}
								onChange={(e) => setPassword(e.target.value)}
								required
							/>
						</div>

						<Button type="submit" className="w-full" disabled={loading}>
							{loading ? "Wird geladen..." : "Anmelden"}
						</Button>
					</form>
				</CardContent>
			</Card>
		</div>
	)
}
