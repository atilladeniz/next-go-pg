"use client"

import Link from "next/link"
import { useRouter } from "next/navigation"
import { useState } from "react"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { signUp } from "@/lib/auth-client"

export default function RegisterPage() {
	const router = useRouter()
	const [name, setName] = useState("")
	const [email, setEmail] = useState("")
	const [password, setPassword] = useState("")
	const [error, setError] = useState("")
	const [loading, setLoading] = useState(false)

	const handleSubmit = async (e: React.FormEvent) => {
		e.preventDefault()
		setError("")
		setLoading(true)

		const result = await signUp.email({
			name,
			email,
			password,
		})

		if (result.error) {
			setError(result.error.message || "Registrierung fehlgeschlagen")
			setLoading(false)
			return
		}

		router.push("/dashboard")
	}

	return (
		<div className="flex min-h-screen items-center justify-center bg-background">
			<Card className="w-full max-w-md">
				<CardHeader className="text-center">
					<CardTitle className="text-2xl">Registrieren</CardTitle>
					<CardDescription>
						Bereits ein Konto?{" "}
						<Link href="/login" className="text-primary hover:underline">
							Anmelden
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
							<Label htmlFor="name">Name</Label>
							<Input
								id="name"
								type="text"
								value={name}
								onChange={(e) => setName(e.target.value)}
								required
								placeholder="Max Mustermann"
							/>
						</div>

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
								minLength={8}
							/>
							<p className="text-xs text-muted-foreground">Mindestens 8 Zeichen</p>
						</div>

						<Button type="submit" className="w-full" disabled={loading}>
							{loading ? "Wird geladen..." : "Registrieren"}
						</Button>
					</form>
				</CardContent>
			</Card>
		</div>
	)
}
