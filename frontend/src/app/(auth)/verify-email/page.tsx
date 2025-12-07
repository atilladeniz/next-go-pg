import { auth } from "@shared/lib/auth-server"
import { Button } from "@shared/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@shared/ui/card"
import Link from "next/link"
import { redirect } from "next/navigation"

interface Props {
	searchParams: Promise<{ token?: string; error?: string }>
}

export default async function VerifyEmailPage({ searchParams }: Props) {
	const { token, error: errorParam } = await searchParams

	if (errorParam) {
		return (
			<div className="flex min-h-screen items-center justify-center bg-background">
				<Card className="w-full max-w-md">
					<CardHeader className="text-center">
						<CardTitle className="text-2xl">Verifizierung fehlgeschlagen</CardTitle>
						<CardDescription>Der Link ist ungültig oder abgelaufen.</CardDescription>
					</CardHeader>
					<CardContent className="space-y-4">
						<p className="text-center text-sm text-muted-foreground">
							Bitte fordere einen neuen Link an.
						</p>
						<Button className="w-full" asChild>
							<Link href="/login">Zur Anmeldung</Link>
						</Button>
					</CardContent>
				</Card>
			</div>
		)
	}

	if (!token) {
		redirect("/login")
	}

	const result = await auth.api.verifyEmail({
		query: { token },
	})

	if (!result) {
		return (
			<div className="flex min-h-screen items-center justify-center bg-background">
				<Card className="w-full max-w-md">
					<CardHeader className="text-center">
						<CardTitle className="text-2xl">Verifizierung fehlgeschlagen</CardTitle>
						<CardDescription>Der Link ist ungültig oder abgelaufen.</CardDescription>
					</CardHeader>
					<CardContent className="space-y-4">
						<p className="text-center text-sm text-muted-foreground">
							Bitte fordere einen neuen Link an.
						</p>
						<Button className="w-full" asChild>
							<Link href="/login">Zur Anmeldung</Link>
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
					<CardTitle className="text-2xl">E-Mail bestätigt!</CardTitle>
					<CardDescription>Deine E-Mail-Adresse wurde erfolgreich verifiziert.</CardDescription>
				</CardHeader>
				<CardContent className="space-y-4">
					<p className="text-center text-sm text-muted-foreground">
						Du kannst dich jetzt anmelden.
					</p>
					<Button className="w-full" asChild>
						<Link href="/login">Jetzt anmelden</Link>
					</Button>
				</CardContent>
			</Card>
		</div>
	)
}
