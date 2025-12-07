"use client"

import { authClient } from "@shared/lib/auth-client"
import { Button } from "@shared/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@shared/ui/card"
import { AlertCircle, Clock, Loader2, XCircle } from "lucide-react"
import Link from "next/link"
import { useRouter, useSearchParams } from "next/navigation"
import { useEffect, useState } from "react"

type Status = "loading" | "success" | "error" | "rate-limited" | "expired" | "invalid"

export default function MagicLinkVerifyPage() {
	const searchParams = useSearchParams()
	const router = useRouter()
	const [status, setStatus] = useState<Status>("loading")

	const token = searchParams.get("token")
	const callbackURL = searchParams.get("callbackURL") || "/dashboard"

	useEffect(() => {
		if (!token) {
			setStatus("invalid")
			return
		}

		async function verifyMagicLink() {
			try {
				const result = await authClient.magicLink.verify({
					query: { token: token! },
				})

				if (result.error) {
					const message = result.error.message?.toLowerCase() || ""
					if (message.includes("too many") || message.includes("rate")) {
						setStatus("rate-limited")
					} else if (message.includes("expired")) {
						setStatus("expired")
					} else if (message.includes("invalid")) {
						setStatus("invalid")
					} else {
						setStatus("error")
					}
					return
				}

				// Success - redirect immediately
				router.push(callbackURL)
			} catch {
				setStatus("error")
			}
		}

		verifyMagicLink()
	}, [token, router, callbackURL])

	return (
		<div className="flex min-h-screen items-center justify-center bg-background">
			{status === "loading" && <LoadingCard />}
			{status === "rate-limited" && <RateLimitedCard />}
			{status === "expired" && <ExpiredCard />}
			{status === "invalid" && <InvalidCard />}
			{status === "error" && <ErrorCard />}
		</div>
	)
}

function LoadingCard() {
	return (
		<Card className="w-full max-w-md">
			<CardHeader className="text-center">
				<div className="mx-auto mb-4 flex h-12 w-12 items-center justify-center">
					<Loader2 className="h-8 w-8 animate-spin text-primary" />
				</div>
				<CardTitle className="text-2xl">Anmeldung wird verifiziert...</CardTitle>
				<CardDescription>Bitte warte einen Moment.</CardDescription>
			</CardHeader>
		</Card>
	)
}

function RateLimitedCard() {
	return (
		<Card className="w-full max-w-md">
			<CardHeader className="text-center">
				<div className="mx-auto mb-4 flex h-12 w-12 items-center justify-center rounded-full bg-yellow-100 dark:bg-yellow-900">
					<Clock className="h-6 w-6 text-yellow-600 dark:text-yellow-400" />
				</div>
				<CardTitle className="text-2xl">Zu viele Anfragen</CardTitle>
				<CardDescription>
					Du hast zu viele Anmeldeversuche gemacht. Bitte warte einen Moment und versuche es erneut.
				</CardDescription>
			</CardHeader>
			<CardContent className="space-y-4">
				<p className="text-center text-sm text-muted-foreground">
					Aus Sicherheitsgründen ist die Anzahl der Anmeldeversuche begrenzt.
				</p>
				<Button className="w-full" asChild>
					<Link href="/login">Zurück zur Anmeldung</Link>
				</Button>
			</CardContent>
		</Card>
	)
}

function ExpiredCard() {
	return (
		<Card className="w-full max-w-md">
			<CardHeader className="text-center">
				<div className="mx-auto mb-4 flex h-12 w-12 items-center justify-center rounded-full bg-orange-100 dark:bg-orange-900">
					<AlertCircle className="h-6 w-6 text-orange-600 dark:text-orange-400" />
				</div>
				<CardTitle className="text-2xl">Link abgelaufen</CardTitle>
				<CardDescription>
					Dieser Anmelde-Link ist leider abgelaufen. Bitte fordere einen neuen Link an.
				</CardDescription>
			</CardHeader>
			<CardContent className="space-y-4">
				<p className="text-center text-sm text-muted-foreground">
					Anmelde-Links sind aus Sicherheitsgründen nur 10 Minuten gültig.
				</p>
				<Button className="w-full" asChild>
					<Link href="/login">Neuen Link anfordern</Link>
				</Button>
			</CardContent>
		</Card>
	)
}

function InvalidCard() {
	return (
		<Card className="w-full max-w-md">
			<CardHeader className="text-center">
				<div className="mx-auto mb-4 flex h-12 w-12 items-center justify-center rounded-full bg-red-100 dark:bg-red-900">
					<XCircle className="h-6 w-6 text-red-600 dark:text-red-400" />
				</div>
				<CardTitle className="text-2xl">Ungültiger Link</CardTitle>
				<CardDescription>Dieser Link ist ungültig oder wurde bereits verwendet.</CardDescription>
			</CardHeader>
			<CardContent className="space-y-4">
				<p className="text-center text-sm text-muted-foreground">
					Jeder Anmelde-Link kann nur einmal verwendet werden.
				</p>
				<Button className="w-full" asChild>
					<Link href="/login">Neuen Link anfordern</Link>
				</Button>
			</CardContent>
		</Card>
	)
}

function ErrorCard() {
	return (
		<Card className="w-full max-w-md">
			<CardHeader className="text-center">
				<div className="mx-auto mb-4 flex h-12 w-12 items-center justify-center rounded-full bg-red-100 dark:bg-red-900">
					<XCircle className="h-6 w-6 text-red-600 dark:text-red-400" />
				</div>
				<CardTitle className="text-2xl">Fehler bei der Anmeldung</CardTitle>
				<CardDescription>Etwas ist schiefgelaufen. Bitte versuche es erneut.</CardDescription>
			</CardHeader>
			<CardContent>
				<Button className="w-full" asChild>
					<Link href="/login">Zurück zur Anmeldung</Link>
				</Button>
			</CardContent>
		</Card>
	)
}
