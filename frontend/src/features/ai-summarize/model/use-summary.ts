"use client"

import { useGetAiSummariesId } from "@shared/api/endpoints/ai/ai"

// useRepoSummary loads a stored summary.
//
// Always call with the REAL summaryId — never substitute `null` / `0`
// when you want to disable the fetch. The key has to stay stable so
// React Query's cache survives toggle cycles (e.g. collapsing an
// accordion card mid-animation must NOT wipe the failure reason that
// was already fetched). Use the `enabled` parameter to gate the fetch.
export function useRepoSummary(summaryId: number, enabled = true) {
	return useGetAiSummariesId(summaryId, {
		query: {
			enabled: enabled && summaryId > 0,
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
