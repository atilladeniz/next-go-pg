"use client"

import { useCallback, useEffect, useRef, useState } from "react"

const API_BASE = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080"

export type ExportFormat = "csv" | "json"
export type ExportDataType = "stats" | "activity" | "all"
export type ExportStatus = "idle" | "pending" | "processing" | "completed" | "failed"

export interface ExportProgress {
	jobId: string
	status: ExportStatus
	progress: number
	message: string
	fileName?: string
	downloadId?: string
	error?: string
}

interface UseExportOptions {
	onComplete?: (downloadId: string, fileName: string) => void
	onError?: (error: string) => void
}

export function useExport(options?: UseExportOptions) {
	const [status, setStatus] = useState<ExportStatus>("idle")
	const [progress, setProgress] = useState(0)
	const [message, setMessage] = useState("")
	const [downloadId, setDownloadId] = useState<string | null>(null)
	const [fileName, setFileName] = useState<string | null>(null)
	const [error, setError] = useState<string | null>(null)
	const [currentJobId, setCurrentJobId] = useState<string | null>(null)

	const eventSourceRef = useRef<EventSource | null>(null)

	// Listen for export progress via SSE
	useEffect(() => {
		if (typeof window === "undefined") return

		const eventSource = new EventSource(`${API_BASE}/api/v1/events`)
		eventSourceRef.current = eventSource

		eventSource.addEventListener("export-progress", (event) => {
			try {
				const data: ExportProgress = JSON.parse(event.data)

				// Only process events for our current job
				if (currentJobId && data.jobId !== currentJobId) return

				setStatus(data.status as ExportStatus)
				setProgress(data.progress)
				setMessage(data.message)

				if (data.downloadId) {
					setDownloadId(data.downloadId)
				}
				if (data.fileName) {
					setFileName(data.fileName)
				}
				if (data.error) {
					setError(data.error)
					options?.onError?.(data.error)
				}

				if (data.status === "completed" && data.downloadId && data.fileName) {
					options?.onComplete?.(data.downloadId, data.fileName)
				}
			} catch {
				// Ignore parse errors
			}
		})

		return () => {
			eventSource.close()
		}
	}, [currentJobId, options])

	const startExport = useCallback(
		async (format: ExportFormat, dataType: ExportDataType) => {
			setStatus("pending")
			setProgress(0)
			setMessage("Export wird gestartet...")
			setDownloadId(null)
			setFileName(null)
			setError(null)

			try {
				const response = await fetch(`${API_BASE}/api/v1/export/start`, {
					method: "POST",
					headers: {
						"Content-Type": "application/json",
					},
					credentials: "include",
					body: JSON.stringify({ format, dataType }),
				})

				if (!response.ok) {
					const errorData = await response.json()
					throw new Error(errorData.error || "Export fehlgeschlagen")
				}

				const data = await response.json()
				setCurrentJobId(data.jobId)
			} catch (err) {
				setStatus("failed")
				setError(err instanceof Error ? err.message : "Unbekannter Fehler")
				options?.onError?.(err instanceof Error ? err.message : "Unbekannter Fehler")
			}
		},
		[options],
	)

	const downloadExport = useCallback(() => {
		if (!downloadId) return

		// Trigger download
		const downloadUrl = `${API_BASE}/api/v1/export/download/${downloadId}`
		const link = document.createElement("a")
		link.href = downloadUrl
		link.download = fileName || "export"
		document.body.appendChild(link)
		link.click()
		document.body.removeChild(link)
	}, [downloadId, fileName])

	const reset = useCallback(() => {
		setStatus("idle")
		setProgress(0)
		setMessage("")
		setDownloadId(null)
		setFileName(null)
		setError(null)
		setCurrentJobId(null)
	}, [])

	return {
		status,
		progress,
		message,
		downloadId,
		fileName,
		error,
		startExport,
		downloadExport,
		reset,
		isExporting: status === "pending" || status === "processing",
		isCompleted: status === "completed",
		isFailed: status === "failed",
	}
}
