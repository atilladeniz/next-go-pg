import { redirect } from "next/navigation"
import { Header } from "@/components/header"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { getSession, getStats } from "@/lib/auth-server"
import { StatsGrid } from "./stats-grid"

export default async function DashboardPage() {
	const [session, stats] = await Promise.all([getSession(), getStats()])

	if (!session) redirect("/login")

	const { user } = session

	return (
		<div className="min-h-screen bg-background">
			<Header user={user} />

			<main className="mx-auto max-w-7xl px-4 py-8">
				<Card>
					<CardHeader>
						<CardTitle>Willkommen, {user.name || "User"}!</CardTitle>
					</CardHeader>
					<CardContent>
						<p className="text-muted-foreground">E-Mail: {user.email}</p>
					</CardContent>
				</Card>

				<StatsGrid initialStats={stats} />
			</main>
		</div>
	)
}
