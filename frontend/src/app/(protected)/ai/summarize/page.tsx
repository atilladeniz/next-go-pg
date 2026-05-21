import { NewRunForm, RunList } from "@features/ai-summarize"
import {
	getAiSummaries,
	getAiSummariesId,
	getGetAiSummariesIdQueryKey,
	getGetAiSummariesQueryKey,
} from "@shared/api/endpoints/ai/ai"
import { getQueryClient } from "@shared/lib"
import { getSession } from "@shared/lib/auth-server"
import { dehydrate, HydrationBoundary } from "@tanstack/react-query"
import { Header } from "@widgets/header"
import { cookies } from "next/headers"
import { redirect } from "next/navigation"

// Server-side prefetch — fills the React Query cache BEFORE the page
// hydrates, so the initial render already has the latest history AND
// each expanded run's detail (incl. step_durations). Without this, the
// browser would tick through ○○○○○ → real state after the GET roundtrip
// completed. URL `?ids=1,2,3` lists every currently-expanded run; we
// prefetch each one in parallel so multi-card open survives reload.
//
// Pattern from frontend/node_modules/next/dist/docs/01-app/01-getting-started/06-fetching-data.md
// + @tanstack/react-query SSR guide.
export default async function AISummarizePage({
	searchParams,
}: {
	searchParams: Promise<{ id?: string; ids?: string }>
}) {
	const session = await getSession()
	if (!session) redirect("/login")

	const { id, ids } = await searchParams
	const rawIds = ids ?? id ?? ""
	const focusedIds = rawIds
		.split(",")
		.map((s) => Number.parseInt(s.trim(), 10))
		.filter((n) => Number.isFinite(n) && n > 0)

	const cookieStore = await cookies()
	const cookieHeader = cookieStore
		.getAll()
		.map((c) => `${c.name}=${c.value}`)
		.join("; ")
	const fetchOpts: RequestInit = {
		headers: { Cookie: cookieHeader },
		cache: "no-store",
	}

	const queryClient = getQueryClient()
	const prefetches: Promise<unknown>[] = [
		queryClient.prefetchQuery({
			queryKey: getGetAiSummariesQueryKey(),
			queryFn: () => getAiSummaries(fetchOpts),
		}),
	]
	for (const focusedId of focusedIds) {
		prefetches.push(
			queryClient.prefetchQuery({
				queryKey: getGetAiSummariesIdQueryKey(focusedId),
				queryFn: () => getAiSummariesId(focusedId, fetchOpts),
			}),
		)
	}
	await Promise.all(prefetches)

	return (
		<HydrationBoundary state={dehydrate(queryClient)}>
			<div className="min-h-screen bg-background">
				<Header user={session.user} />
				<main className="mx-auto max-w-3xl space-y-4 px-4 py-8">
					<div>
						<h1 className="text-2xl font-semibold tracking-tight">Repository-Zusammenfassung</h1>
						<p className="text-sm text-muted-foreground">
							URL eintragen, mehrere Runs parallel — Karten unten zeigen Live-Fortschritt.
						</p>
					</div>
					<NewRunForm />
					<RunList />
				</main>
			</div>
		</HydrationBoundary>
	)
}
