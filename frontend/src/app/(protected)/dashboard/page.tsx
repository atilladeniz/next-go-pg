import { dehydrate, HydrationBoundary } from "@tanstack/react-query"
import { cookies } from "next/headers"
import { redirect } from "next/navigation"
import { getGetStatsQueryKey, getStats } from "@/api/endpoints/users/users"
import { Header } from "@/components/header"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { getSession } from "@/lib/auth-server"
import { getQueryClient } from "@/lib/get-query-client"
import { StatsGrid } from "./stats-grid"

export default async function DashboardPage() {
	const session = await getSession()
	if (!session) redirect("/login")

	const { user } = session

	// Get cookies for server-side fetch
	const cookieStore = await cookies()
	const cookieHeader = cookieStore
		.getAll()
		.map((c) => `${c.name}=${c.value}`)
		.join("; ")

	// Prefetch with React Query
	const queryClient = getQueryClient()
	await queryClient.prefetchQuery({
		queryKey: getGetStatsQueryKey(),
		queryFn: () =>
			getStats({
				headers: { Cookie: cookieHeader },
				cache: "no-store",
			}),
	})

	return (
		<HydrationBoundary state={dehydrate(queryClient)}>
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

					<StatsGrid />
				</main>
			</div>
		</HydrationBoundary>
	)
}
