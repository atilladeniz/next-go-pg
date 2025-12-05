import { UserInfo } from "@entities/user"
import { StatsGrid } from "@features/stats"
import { getGetStatsQueryKey, getStats } from "@shared/api/endpoints/users/users"
import { getQueryClient } from "@shared/lib"
import { getSession } from "@shared/lib/auth-server"
import { Card, CardContent, CardHeader, CardTitle } from "@shared/ui/card"
import { dehydrate, HydrationBoundary } from "@tanstack/react-query"
import { Header } from "@widgets/header"
import { cookies } from "next/headers"
import { redirect } from "next/navigation"

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
							<UserInfo user={user} showEmail />
						</CardContent>
					</Card>

					<StatsGrid />
				</main>
			</div>
		</HydrationBoundary>
	)
}
