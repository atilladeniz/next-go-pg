"use client"

import { useGetAiSummaries } from "@shared/api/endpoints/ai/ai"
import type { AiworkflowsInterfacesHttpRepoSummaryListItem } from "@shared/api/models"

// useHistory returns the authenticated user's recent runs. We unwrap the
// customFetch envelope and re-poll while any non-terminal run is present
// so the list reflects ongoing work without manual refresh.
export function useHistory() {
	const query = useGetAiSummaries({
		query: {
			refetchInterval: (q) => {
				const envelope = q.state.data as
					| { data?: { items?: AiworkflowsInterfacesHttpRepoSummaryListItem[] } }
					| undefined
				const items = envelope?.data?.items ?? []
				const stillRunning = items.some((it) => it.status === "pending" || it.status === "running")
				return stillRunning ? 3000 : false
			},
		},
	})

	const envelope = query.data as
		| { data?: { items?: AiworkflowsInterfacesHttpRepoSummaryListItem[] } }
		| undefined
	const items = envelope?.data?.items ?? []

	return { items, isLoading: query.isLoading, isError: query.isError }
}
