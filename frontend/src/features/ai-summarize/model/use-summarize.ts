"use client"

import { usePostAiSummarizeRepo } from "@shared/api/endpoints/ai/ai"

// useSummarizeRepo wraps the Orval-generated mutation. Returns the
// raw mutation object so the caller can decide between mutateAsync
// (for await-on-submit) and the imperative mutate.
export function useSummarizeRepo() {
	return usePostAiSummarizeRepo()
}
