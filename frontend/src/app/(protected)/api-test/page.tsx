"use client"

import { useState } from "react"
import { Header } from "@/components/header"
import { useServerSession } from "@/components/session-provider"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"

const API_BASE = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080"

type TestResult = {
	name: string
	status: "idle" | "loading" | "success" | "error"
	response?: unknown
	error?: string
	duration?: number
}

export default function ApiTestPage() {
	const { user } = useServerSession()
	const [results, setResults] = useState<Record<string, TestResult>>({
		health: { name: "Health Check", status: "idle" },
		healthReady: { name: "Health Ready", status: "idle" },
		healthLive: { name: "Health Live", status: "idle" },
		publicHello: { name: "Public Hello", status: "idle" },
		protectedHello: { name: "Protected Hello", status: "idle" },
		currentUser: { name: "Current User", status: "idle" },
	})

	const updateResult = (key: string, update: Partial<TestResult>) => {
		setResults((prev) => ({
			...prev,
			[key]: { ...prev[key], ...update },
		}))
	}

	const runTest = async (key: string, url: string, requiresAuth = false) => {
		updateResult(key, { status: "loading", response: undefined, error: undefined })
		const start = performance.now()

		try {
			const res = await fetch(url, {
				credentials: requiresAuth ? "include" : "omit",
			})
			const duration = Math.round(performance.now() - start)

			if (!res.ok) {
				const errorText = await res.text()
				updateResult(key, {
					status: "error",
					error: `${res.status} ${res.statusText}: ${errorText}`,
					duration,
				})
				return
			}

			const contentType = res.headers.get("content-type") || ""
			let data: unknown
			if (contentType.includes("application/json")) {
				data = await res.json()
			} else {
				const text = await res.text()
				data = { message: text }
			}
			updateResult(key, { status: "success", response: data, duration })
		} catch (err) {
			const duration = Math.round(performance.now() - start)
			updateResult(key, {
				status: "error",
				error: err instanceof Error ? err.message : "Unknown error",
				duration,
			})
		}
	}

	const runAllTests = async () => {
		await Promise.all([
			runTest("health", `${API_BASE}/health`),
			runTest("healthReady", `${API_BASE}/health/ready`),
			runTest("healthLive", `${API_BASE}/health/live`),
			runTest("publicHello", `${API_BASE}/api/v1/hello`),
			runTest("protectedHello", `${API_BASE}/api/v1/protected/hello`, true),
			runTest("currentUser", `${API_BASE}/api/v1/me`, true),
		])
	}

	return (
		<div className="min-h-screen bg-background">
			<Header />

			<main className="mx-auto max-w-7xl px-4 py-8">
				<div className="mb-6 flex items-center justify-between">
					<div>
						<h1 className="text-2xl font-bold">API Test</h1>
						<p className="text-muted-foreground">
							Teste die Verbindung zum Go Backend ({API_BASE})
						</p>
					</div>
					<Button onClick={runAllTests}>Alle Tests ausführen</Button>
				</div>

				<div className="grid gap-4 md:grid-cols-2">
					{/* Health Checks */}
					<Card>
						<CardHeader>
							<CardTitle className="text-lg">Health Endpoints</CardTitle>
							<CardDescription>Basis-Gesundheitschecks des Backends</CardDescription>
						</CardHeader>
						<CardContent className="space-y-3">
							<TestRow
								result={results.health}
								onRun={() => runTest("health", `${API_BASE}/health`)}
								endpoint="GET /health"
							/>
							<TestRow
								result={results.healthReady}
								onRun={() => runTest("healthReady", `${API_BASE}/health/ready`)}
								endpoint="GET /health/ready"
							/>
							<TestRow
								result={results.healthLive}
								onRun={() => runTest("healthLive", `${API_BASE}/health/live`)}
								endpoint="GET /health/live"
							/>
						</CardContent>
					</Card>

					{/* Public Endpoints */}
					<Card>
						<CardHeader>
							<CardTitle className="text-lg">Public Endpoints</CardTitle>
							<CardDescription>Öffentlich zugängliche API-Endpunkte</CardDescription>
						</CardHeader>
						<CardContent className="space-y-3">
							<TestRow
								result={results.publicHello}
								onRun={() => runTest("publicHello", `${API_BASE}/api/v1/hello`)}
								endpoint="GET /api/v1/hello"
							/>
						</CardContent>
					</Card>

					{/* Protected Endpoints */}
					<Card className="md:col-span-2">
						<CardHeader>
							<CardTitle className="text-lg">Protected Endpoints</CardTitle>
							<CardDescription>
								Authentifizierung erforderlich (eingeloggt als: {user.email})
							</CardDescription>
						</CardHeader>
						<CardContent className="space-y-3">
							<TestRow
								result={results.protectedHello}
								onRun={() => runTest("protectedHello", `${API_BASE}/api/v1/protected/hello`, true)}
								endpoint="GET /api/v1/protected/hello"
							/>
							<TestRow
								result={results.currentUser}
								onRun={() => runTest("currentUser", `${API_BASE}/api/v1/me`, true)}
								endpoint="GET /api/v1/me"
							/>
						</CardContent>
					</Card>
				</div>
			</main>
		</div>
	)
}

function TestRow({
	result,
	onRun,
	endpoint,
}: {
	result: TestResult
	onRun: () => void
	endpoint: string
}) {
	return (
		<div className="rounded-lg border p-3">
			<div className="flex items-center justify-between">
				<div className="flex items-center gap-3">
					<StatusBadge status={result.status} />
					<div>
						<div className="font-medium">{result.name}</div>
						<code className="text-xs text-muted-foreground">{endpoint}</code>
					</div>
				</div>
				<div className="flex items-center gap-2">
					{result.duration !== undefined && (
						<span className="text-xs text-muted-foreground">{result.duration}ms</span>
					)}
					<Button
						size="sm"
						variant="outline"
						onClick={onRun}
						disabled={result.status === "loading"}
					>
						{result.status === "loading" ? "..." : "Test"}
					</Button>
				</div>
			</div>
			{result.response !== undefined && (
				<pre className="mt-2 overflow-auto rounded bg-muted p-2 text-xs">
					{JSON.stringify(result.response, null, 2)}
				</pre>
			)}
			{result.error && (
				<div className="mt-2 rounded bg-destructive/10 p-2 text-xs text-destructive">
					{result.error}
				</div>
			)}
		</div>
	)
}

function StatusBadge({ status }: { status: TestResult["status"] }) {
	switch (status) {
		case "idle":
			return <Badge variant="secondary">Idle</Badge>
		case "loading":
			return <Badge variant="secondary">Loading...</Badge>
		case "success":
			return <Badge className="bg-green-600">Success</Badge>
		case "error":
			return <Badge variant="destructive">Error</Badge>
	}
}
