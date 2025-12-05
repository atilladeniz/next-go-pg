"use client"

import { getGetStatsQueryKey } from "@shared/api/endpoints/users/users"
import { useQueryClient } from "@tanstack/react-query"
import { useCallback, useEffect, useRef, useState } from "react"

const API_BASE = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080"

export function useSSE() {
	const queryClient = useQueryClient()
	const eventSourceRef = useRef<EventSource | null>(null)
	const reconnectTimeoutRef = useRef<NodeJS.Timeout | null>(null)
	const [isConnected, setIsConnected] = useState(false)

	const connect = useCallback(() => {
		// Only run in browser
		if (typeof window === "undefined") return

		// Close existing connection
		if (eventSourceRef.current) {
			eventSourceRef.current.close()
		}

		const eventSource = new EventSource(`${API_BASE}/api/v1/events`)
		eventSourceRef.current = eventSource

		eventSource.onopen = () => {
			console.log("[SSE] Connected")
			setIsConnected(true)
		}

		eventSource.onerror = () => {
			console.log("[SSE] Connection error, reconnecting in 5s...")
			setIsConnected(false)
			eventSource.close()

			// Reconnect after 5 seconds
			reconnectTimeoutRef.current = setTimeout(() => {
				connect()
			}, 5000)
		}

		// Handle connection event
		eventSource.addEventListener("connected", () => {
			console.log("[SSE] Server confirmed connection")
		})

		// Handle stats-updated event
		eventSource.addEventListener("stats-updated", () => {
			console.log("[SSE] Stats updated, invalidating cache...")
			queryClient.invalidateQueries({ queryKey: getGetStatsQueryKey() })
		})
	}, [queryClient])

	useEffect(() => {
		connect()

		return () => {
			if (eventSourceRef.current) {
				eventSourceRef.current.close()
			}
			if (reconnectTimeoutRef.current) {
				clearTimeout(reconnectTimeoutRef.current)
			}
		}
	}, [connect])

	return { isConnected }
}
