"use client"

import { getGetAiSummariesIdQueryKey, getGetAiSummariesQueryKey } from "@shared/api/endpoints/ai/ai"
import { useQueryClient } from "@tanstack/react-query"
import { useEffect, useState } from "react"

const API_BASE = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080"

// Backend `events.Publisher` emits two SSE-payload kinds over the same
// `ai-progress` channel:
//   - kind=step       — fine-grained per-step / per-file transitions
//   - kind=lifecycle  — run-level (running/completed/failed/cancelled)

export type StepName = "clone" | "traverse" | "summarize_files" | "aggregate" | "store"
export type StepState = "started" | "completed" | "failed" | "progress"
export type RunStatus = "pending" | "running" | "completed" | "failed" | "cancelled"

interface ProgressPayload {
	kind?: "step" | "lifecycle"
	summaryId: number
	userId?: string
	step?: StepName
	state?: StepState
	status?: RunStatus
	durationMs?: number
	filename?: string
	fileIndex?: number
	fileCount?: number
	reason?: string
}

export type StepStatus = "pending" | "running" | "completed" | "failed"

export interface StepView {
	name: StepName
	status: StepStatus
	durationMs?: number
	fileIndex?: number
	fileCount?: number
	filename?: string
	reason?: string
}

export interface ProgressView {
	steps: Record<StepName, StepView>
	runStatus: RunStatus
}

export const STEP_ORDER: StepName[] = ["clone", "traverse", "summarize_files", "aggregate", "store"]

const initialView = (): ProgressView => ({
	steps: {
		clone: { name: "clone", status: "pending" },
		traverse: { name: "traverse", status: "pending" },
		summarize_files: { name: "summarize_files", status: "pending" },
		aggregate: { name: "aggregate", status: "pending" },
		store: { name: "store", status: "pending" },
	},
	runStatus: "pending",
})

// ─── shared SSE singleton ────────────────────────────────────────────
//
// One EventSource per browser tab, fanned out to every interested
// component. Running 3 workflows in parallel doesn't open 3 connections.

type Listener = (data: ProgressPayload) => void
let connection: EventSource | null = null
const listeners = new Set<Listener>()

function ensureConnection() {
	if (typeof window === "undefined") return null
	if (connection) return connection
	const es = new EventSource(`${API_BASE}/api/v1/events`)
	es.addEventListener("ai-progress", (event) => {
		let data: ProgressPayload
		try {
			data = JSON.parse((event as MessageEvent).data) as ProgressPayload
		} catch {
			return
		}
		for (const l of listeners) l(data)
	})
	connection = es
	return es
}

function subscribe(l: Listener): () => void {
	ensureConnection()
	listeners.add(l)
	return () => {
		listeners.delete(l)
		// Lazy: keep the connection open as long as the tab lives —
		// re-opening is more expensive than idling.
	}
}

// useAIProgress subscribes to the shared SSE stream and aggregates step
// events for ONE summary run. Pass null to disable.
export function useAIProgress(summaryId: number | null): ProgressView {
	const queryClient = useQueryClient()
	const [view, setView] = useState<ProgressView>(initialView())

	useEffect(() => {
		if (summaryId === null) {
			setView(initialView())
			return
		}
		setView(initialView())
		const off = subscribe((data) => {
			if (data.summaryId !== summaryId) return
			setView((prev) => reduce(prev, data))

			if (
				data.kind === "lifecycle" &&
				(data.status === "completed" || data.status === "failed" || data.status === "cancelled")
			) {
				queryClient.invalidateQueries({ queryKey: getGetAiSummariesIdQueryKey(summaryId) })
				queryClient.invalidateQueries({ queryKey: getGetAiSummariesQueryKey() })
			}
		})
		return off
	}, [summaryId, queryClient])

	return view
}

function reduce(prev: ProgressView, ev: ProgressPayload): ProgressView {
	if (ev.kind === "lifecycle" && ev.status) {
		return { ...prev, runStatus: ev.status }
	}
	if (ev.kind !== "step" || !ev.step) return prev

	const current = prev.steps[ev.step]
	if (!current) return prev

	let next: StepView = current
	switch (ev.state) {
		case "started":
			next = { ...current, status: "running", fileCount: ev.fileCount ?? current.fileCount }
			break
		case "progress":
			next = {
				...current,
				status: "running",
				fileIndex: ev.fileIndex,
				fileCount: ev.fileCount ?? current.fileCount,
				filename: ev.filename,
			}
			break
		case "completed":
			next = {
				...current,
				status: "completed",
				durationMs: ev.durationMs,
				fileCount: ev.fileCount ?? current.fileCount,
			}
			break
		case "failed":
			next = {
				...current,
				status: "failed",
				durationMs: ev.durationMs,
				reason: ev.reason,
			}
			break
		default:
			return prev
	}
	return { ...prev, steps: { ...prev.steps, [ev.step]: next } }
}
