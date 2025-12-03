"use client"

import { useRouter } from "next/navigation"
import { ModeToggle } from "@/components/mode-toggle"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { signOut, useSession } from "@/lib/auth-client"

export default function DashboardPage() {
	const router = useRouter()
	const { data: session, isPending } = useSession()

	const handleSignOut = async () => {
		await signOut()
		router.push("/")
	}

	if (isPending) {
		return (
			<div className="flex min-h-screen items-center justify-center bg-background">
				<div className="text-muted-foreground">Laden...</div>
			</div>
		)
	}

	return (
		<div className="min-h-screen bg-background">
			<header className="border-b">
				<div className="mx-auto flex max-w-7xl items-center justify-between px-4 py-4">
					<h1 className="text-xl font-semibold">Dashboard</h1>
					<div className="flex items-center gap-2">
						<ModeToggle />
						<Button variant="outline" onClick={handleSignOut}>
							Abmelden
						</Button>
					</div>
				</div>
			</header>

			<main className="mx-auto max-w-7xl px-4 py-8">
				<Card>
					<CardHeader>
						<CardTitle>Willkommen, {session?.user?.name || "User"}!</CardTitle>
					</CardHeader>
					<CardContent>
						<p className="text-muted-foreground">E-Mail: {session?.user?.email}</p>
					</CardContent>
				</Card>

				<div className="mt-6 grid gap-6 md:grid-cols-2 lg:grid-cols-3">
					<Card>
						<CardHeader>
							<CardTitle className="text-base">Statistiken</CardTitle>
						</CardHeader>
						<CardContent>
							<p className="text-3xl font-bold text-primary">0</p>
							<p className="text-sm text-muted-foreground">Aktive Projekte</p>
						</CardContent>
					</Card>
					<Card>
						<CardHeader>
							<CardTitle className="text-base">Aktivität</CardTitle>
						</CardHeader>
						<CardContent>
							<p className="text-3xl font-bold text-green-600">0</p>
							<p className="text-sm text-muted-foreground">Heute</p>
						</CardContent>
					</Card>
					<Card>
						<CardHeader>
							<CardTitle className="text-base">Benachrichtigungen</CardTitle>
						</CardHeader>
						<CardContent>
							<p className="text-3xl font-bold text-orange-600">0</p>
							<p className="text-sm text-muted-foreground">Ungelesen</p>
						</CardContent>
					</Card>
				</div>
			</main>
		</div>
	)
}
