/**
 * Loki Transport for Pino
 * Sends logs directly to Loki via HTTP
 */

interface LokiEntry {
	timestamp: number
	line: string
	level: string
}

class LokiTransport {
	private url: string
	private serviceName: string
	private batch: LokiEntry[] = []
	private batchSize = 20
	private flushInterval: ReturnType<typeof setInterval> | null = null

	constructor(url: string, serviceName: string) {
		this.url = url
		this.serviceName = serviceName

		// Start background flusher
		if (typeof window === "undefined") {
			// Server-side only
			this.flushInterval = setInterval(() => this.flush(), 2000)
		}
	}

	push(level: string, msg: string, data: Record<string, unknown>) {
		const logLine = JSON.stringify({
			level,
			service: this.serviceName,
			message: msg,
			...data,
			timestamp: new Date().toISOString(),
		})

		this.batch.push({
			timestamp: Date.now() * 1000000, // nanoseconds
			line: logLine,
			level,
		})

		if (this.batch.length >= this.batchSize) {
			this.flush()
		}
	}

	async flush() {
		if (this.batch.length === 0) return

		const entries = [...this.batch]
		this.batch = []

		// Group by level
		const streams: Record<string, LokiEntry[]> = {}
		for (const entry of entries) {
			if (!streams[entry.level]) {
				streams[entry.level] = []
			}
			streams[entry.level].push(entry)
		}

		// Build payload
		const streamArr = Object.entries(streams).map(([level, levelEntries]) => {
			const values = levelEntries.map((e) => [String(e.timestamp), e.line])
			return {
				stream: {
					service: this.serviceName,
					level,
					job: "frontend",
				},
				values,
			}
		})

		const payload = { streams: streamArr }

		try {
			await fetch(this.url, {
				method: "POST",
				headers: { "Content-Type": "application/json" },
				body: JSON.stringify(payload),
			})
		} catch {
			// Silently fail - don't break the app if Loki is down
		}
	}

	close() {
		if (this.flushInterval) {
			clearInterval(this.flushInterval)
		}
		this.flush()
	}
}

// Singleton instance
let lokiTransport: LokiTransport | null = null

export function getLokiTransport(): LokiTransport | null {
	if (typeof window !== "undefined") {
		// Don't send logs from browser to Loki
		return null
	}

	if (!lokiTransport) {
		const lokiUrl = process.env.LOKI_URL
		if (lokiUrl) {
			lokiTransport = new LokiTransport(lokiUrl, "next-go-pg-frontend")
		}
	}

	return lokiTransport
}

export function sendToLoki(level: string, msg: string, data: Record<string, unknown>) {
	const transport = getLokiTransport()
	if (transport) {
		transport.push(level, msg, data)
	}
}
