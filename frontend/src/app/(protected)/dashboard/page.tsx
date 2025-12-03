import { redirect } from "next/navigation"
import { Header } from "@/components/header"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { getSession, getStats } from "@/lib/auth-server"

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

				<div className="mt-6 grid gap-6 md:grid-cols-2 lg:grid-cols-3">
					<StatsCard title="Projekte" value={stats?.projectCount ?? 0} label="Aktive Projekte" />
					<StatsCard
						title="Aktivität"
						value={stats?.activityToday ?? 0}
						label="Heute"
						valueClassName="text-green-600"
					/>
					<StatsCard
						title="Benachrichtigungen"
						value={stats?.notifications ?? 0}
						label="Ungelesen"
						valueClassName="text-orange-600"
					/>
				</div>
			</main>
		</div>
	)
}

function StatsCard({
	title,
	value,
	label,
	valueClassName = "text-primary",
}: {
	title: string
	value: number
	label: string
	valueClassName?: string
}) {
	return (
		<Card>
			<CardHeader>
				<CardTitle className="text-base">{title}</CardTitle>
			</CardHeader>
			<CardContent>
				<p className={`text-3xl font-bold ${valueClassName}`}>{value}</p>
				<p className="text-sm text-muted-foreground">{label}</p>
			</CardContent>
		</Card>
	)
}
