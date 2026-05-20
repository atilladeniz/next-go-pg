"use client"

import { useGetAiSummariesId } from "@shared/api/endpoints/ai/ai"

// useRepoSummary loads a stored summary. Pass null to disable the query
// (e.g. before a run has been started).
export function useRepoSummary(summaryId: number | null) {
	return useGetAiSummariesId(summaryId ?? 0, {
		query: {
			enabled: summaryId !== null,
			// Poll occasionally as a safety net in case the SSE event was
			// missed (tab backgrounded, connection blip). Backstop only;
			// the primary update path is SSE invalidation.
			refetchInterval: (q) => {
				// customFetch wraps the response as { data, status, headers };
				// the persisted `status` field lives at q.state.data.data.status.
				const envelope = q.state.data as { data?: { status?: string } } | undefined
				const status = envelope?.data?.status
				if (!status) return false
				if (status === "completed" || status === "failed" || status === "cancelled") {
					return false
				}
				return 5000
			},
		},
	})
}
