"use client"

import type { AiworkflowsInterfacesHttpRepoSummaryResponse } from "@shared/api/models"
import { Button } from "@shared/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@shared/ui/card"
import { Input } from "@shared/ui/input"
import { Progress } from "@shared/ui/progress"
import { CheckCircle2, FileSearch, Loader2, Sparkles, XCircle } from "lucide-react"
import { useMemo, useState } from "react"
import { type AIProgressStep, useAIProgress } from "../model/use-ai-progress"
import { useSummarizeRepo } from "../model/use-summarize"
import { useRepoSummary } from "../model/use-summary"

const stepLabel: Record<AIProgressStep, string> = {
	idle: "Bereit",
	started: "Repository klonen…",
	summarize_file: "Dateien analysieren…",
	store: "Zusammenfassung speichern…",
}

export function SummarizeCard() {
	const [repoUrl, setRepoUrl] = useState("")
	const [summaryId, setSummaryId] = useState<number | null>(null)

	const mutation = useSummarizeRepo()
	const progress = useAIProgress(summaryId)
	const query = useRepoSummary(summaryId)

	// customFetch wraps the response as { data, status, headers }; the real
	// payload lives at query.data.data.
	const queryEnvelope = query.data as
		| { data?: AiworkflowsInterfacesHttpRepoSummaryResponse }
		| undefined
	const result = queryEnvelope?.data
	const isTerminal =
		result?.status === "completed" || result?.status === "failed" || result?.status === "cancelled"

	const progressPercent = useMemo(() => {
		if (!summaryId) return 0
		if (isTerminal) return 100
		if (progress.step === "started") return 10
		if (progress.step === "summarize_file" && progress.fileCount && progress.fileIndex) {
			// Map per-file progress to 20-80% band; remaining 20% covers aggregate+store.
			const band = (progress.fileIndex / progress.fileCount) * 60
			return Math.min(80, 20 + band)
		}
		if (progress.step === "store") return 95
		return 5
	}, [summaryId, isTerminal, progress])

	const submitDisabled = mutation.isPending || (summaryId !== null && !isTerminal)

	const handleSubmit = async (e: React.FormEvent) => {
		e.preventDefault()
		setSummaryId(null)
		try {
			// customFetch wraps the body as { data, status, headers } so the
			// 202 payload's `summaryId` lives at `response.data.summaryId`.
			const response = (await mutation.mutateAsync({ data: { repoUrl } })) as {
				data?: { summaryId?: number }
				status?: number
			}
			const id = response?.data?.summaryId
			if (typeof id === "number") {
				setSummaryId(id)
			}
		} catch {
			// mutation error surfaces via mutation.error below
		}
	}

	const reset = () => {
		setSummaryId(null)
		setRepoUrl("")
		mutation.reset()
	}

	return (
		<Card className="w-full">
			<CardHeader>
				<CardTitle className="flex items-center gap-2">
					<Sparkles className="h-5 w-5" />
					Repository zusammenfassen
				</CardTitle>
				<CardDescription>
					Eine öffentliche Git-Repository-URL eingeben. Wir klonen das Repo, analysieren relevante
					Dateien per LLM und liefern eine Zusammenfassung.
				</CardDescription>
			</CardHeader>
			<CardContent className="space-y-6">
				<form className="space-y-3" onSubmit={handleSubmit}>
					<div className="space-y-2">
						<label htmlFor="repo-url" className="text-sm font-medium">
							Repository-URL
						</label>
						<Input
							id="repo-url"
							type="url"
							placeholder="https://github.com/owner/repo"
							value={repoUrl}
							onChange={(e) => setRepoUrl(e.target.value)}
							required
							disabled={submitDisabled}
						/>
						<p className="text-xs text-muted-foreground">
							Nur http(s)-URLs. Privatrepos werden derzeit nicht unterstützt.
						</p>
					</div>
					<Button type="submit" className="w-full" disabled={submitDisabled || !repoUrl}>
						{mutation.isPending ? (
							<>
								<Loader2 className="mr-2 h-4 w-4 animate-spin" />
								Starten…
							</>
						) : (
							<>
								<FileSearch className="mr-2 h-4 w-4" />
								Zusammenfassen
							</>
						)}
					</Button>
				</form>

				{mutation.isError && (
					<div className="flex items-center gap-2 rounded-lg border border-destructive/40 bg-destructive/5 p-3 text-sm text-destructive">
						<XCircle className="h-4 w-4" />
						<span>{(mutation.error as Error | null)?.message ?? "Anfrage fehlgeschlagen"}</span>
					</div>
				)}

				{summaryId !== null && (
					<div className="space-y-3 rounded-lg border bg-muted/30 p-4">
						<div className="flex items-center gap-3">
							{!isTerminal && <Loader2 className="h-5 w-5 animate-spin text-primary" />}
							{result?.status === "completed" && (
								<CheckCircle2 className="h-5 w-5 text-green-500" />
							)}
							{result?.status === "failed" && <XCircle className="h-5 w-5 text-destructive" />}
							{result?.status === "cancelled" && (
								<XCircle className="h-5 w-5 text-muted-foreground" />
							)}
							<div className="flex-1">
								<div className="font-medium">
									{result?.status === "completed" && "Fertig"}
									{result?.status === "failed" && "Fehlgeschlagen"}
									{result?.status === "cancelled" && "Abgebrochen"}
									{!isTerminal && stepLabel[progress.step]}
								</div>
								{progress.step === "summarize_file" && progress.fileCount ? (
									<div className="text-xs text-muted-foreground">
										Datei {progress.fileIndex} von {progress.fileCount}
										{progress.filename ? ` — ${progress.filename}` : ""}
									</div>
								) : null}
								{result?.failReason && (
									<div className="text-xs text-destructive">{result.failReason}</div>
								)}
							</div>
						</div>
						<Progress value={progressPercent} className="h-2" />
					</div>
				)}

				{result?.status === "completed" && result.summary && (
					<div className="space-y-4">
						<div>
							<h3 className="mb-2 text-sm font-semibold uppercase tracking-wide text-muted-foreground">
								Repository-Zusammenfassung
							</h3>
							<p className="whitespace-pre-wrap text-sm">{result.summary}</p>
						</div>

						{result.files && result.files.length > 0 && (
							<details className="rounded-lg border">
								<summary className="cursor-pointer p-3 text-sm font-medium">
									Dateien ({result.files.length})
								</summary>
								<div className="divide-y border-t">
									{result.files.map((f) => (
										<div key={f.filename} className="space-y-1 p-3">
											<div className="font-mono text-xs">{f.filename}</div>
											<div className="text-sm text-muted-foreground">{f.summary}</div>
										</div>
									))}
								</div>
							</details>
						)}

						<Button variant="outline" onClick={reset}>
							Neue Zusammenfassung
						</Button>
					</div>
				)}
			</CardContent>
		</Card>
	)
}
