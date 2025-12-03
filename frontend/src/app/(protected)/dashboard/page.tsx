"use client"

import { Header } from "@/components/header"
import { useServerSession } from "@/components/session-provider"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"

export default function DashboardPage() {
	const { user } = useServerSession()

	return (
		<div className="min-h-screen bg-background">
			<Header />

			<main className="mx-auto max-w-7xl px-4 py-8">
				<Card>
					<CardHeader>
						<CardTitle>Willkommen, {user.name || "User"}!</CardTitle>
					</CardHeader>
					<CardContent>
						<p className="text-muted-foreground">E-Mail: {user.email}</p>
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
