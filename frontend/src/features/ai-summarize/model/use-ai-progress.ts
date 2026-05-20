"use client"

import { getGetAiSummariesIdQueryKey } from "@shared/api/endpoints/ai/ai"
import { useQueryClient } from "@tanstack/react-query"
import { useEffect, useState } from "react"

const API_BASE = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080"

export type AIProgressStep = "idle" | "started" | "summarize_file" | "store"

export type AIProgressStatus = "running" | "completed" | "failed" | "cancelled"

export interface AIProgressEvent {
	summaryId: number
	userId: string
	step: AIProgressStep
	status?: AIProgressStatus
	filename?: string
	fileIndex?: number
	fileCount?: number
	reason?: string
}

export interface AIProgressState {
	step: AIProgressStep
	status: AIProgressStatus | "pending"
	filename?: string
	fileIndex?: number
	fileCount?: number
	reason?: string
}

const initial: AIProgressState = { step: "idle", status: "pending" }

// useAIProgress subscribes to the platform SSE stream and surfaces
// `ai-progress` events for a single summary run. Pass null to disable.
export function useAIProgress(summaryId: number | null): AIProgressState {
	const queryClient = useQueryClient()
	const [state, setState] = useState<AIProgressState>(initial)

	useEffect(() => {
		if (summaryId === null) {
			setState(initial)
			return
		}
		if (typeof window === "undefined") return

		const eventSource = new EventSource(`${API_BASE}/api/v1/events`)
		eventSource.addEventListener("ai-progress", (event) => {
			try {
				const data = JSON.parse((event as MessageEvent).data) as AIProgressEvent
				if (data.summaryId !== summaryId) return
				setState({
					step: data.step,
					status: data.status ?? "running",
					filename: data.filename,
					fileIndex: data.fileIndex,
					fileCount: data.fileCount,
					reason: data.reason,
				})
				if (
					data.status === "completed" ||
					data.status === "failed" ||
					data.status === "cancelled"
				) {
					queryClient.invalidateQueries({ queryKey: getGetAiSummariesIdQueryKey(summaryId) })
				}
			} catch {
				// Ignore parse errors
			}
		})

		return () => {
			eventSource.close()
		}
	}, [summaryId, queryClient])

	return state
}
