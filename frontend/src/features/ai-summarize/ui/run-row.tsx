"use client"

import type { AiworkflowsInterfacesHttpRepoSummaryResponse } from "@shared/api/models"
import { AccordionContent, AccordionItem, AccordionTrigger } from "@shared/ui/accordion"
import { CheckCircle2, Circle, Loader2, XCircle } from "lucide-react"
import { useMemo } from "react"
import {
	STEP_ORDER,
	type StepName,
	type StepStatus,
	type StepView,
	useAIProgress,
} from "../model/use-ai-progress"
import { useRepoSummary } from "../model/use-summary"

const stepLabel: Record<StepName, string> = {
	clone: "Repository klonen",
	traverse: "Dateien analysieren",
	summarize_files: "Dateien zusammenfassen",
	aggregate: "Gesamt-Zusammenfassung",
	store: "Ergebnis speichern",
}

function formatRepoLabel(url: string): string {
	try {
		const u = new URL(url)
		const parts = u.pathname.replace(/^\/+|\/+$/g, "").split("/")
		if (parts.length >= 2) return `${parts[0]}/${parts[1].replace(/\.git$/, "")}`
		return u.hostname + u.pathname
	} catch {
		return url
	}
}

function formatDuration(ms?: number): string {
	if (typeof ms !== "number" || ms <= 0) return ""
	if (ms < 1000) return `${ms}ms`
	const s = ms / 1000
	if (s < 60) return `${s.toFixed(2)}s`
	const m = Math.floor(s / 60)
	const rem = (s - m * 60).toFixed(0)
	return `${m}m ${rem}s`
}

function formatRelative(iso?: string): string {
	if (!iso) return ""
	const d = new Date(iso)
	if (Number.isNaN(d.getTime())) return ""
	const diffMs = Date.now() - d.getTime()
	const m = Math.floor(diffMs / 60_000)
	if (m < 1) return "gerade"
	if (m < 60) return `vor ${m} min`
	const h = Math.floor(m / 60)
	if (h < 24) return `vor ${h} h`
	return `vor ${Math.floor(h / 24)} d`
}

type Tone = "idle" | "running" | "success" | "error"

const statusTone = (status?: string): Tone => {
	if (status === "completed") return "success"
	if (status === "failed" || status === "cancelled") return "error"
	if (status === "running" || status === "pending") return "running"
	return "idle"
}

const toneDot: Record<Tone, string> = {
	idle: "bg-muted-foreground",
	running: "bg-primary animate-pulse",
	success: "bg-emerald-500",
	error: "bg-destructive",
}

const toneText: Record<Tone, string> = {
	idle: "text-muted-foreground",
	running: "text-primary",
	success: "text-emerald-600 dark:text-emerald-400",
	error: "text-destructive",
}

const statusLabel: Record<string, string> = {
	pending: "Wartet",
	running: "Läuft",
	completed: "Fertig",
	failed: "Fehlgeschlagen",
	cancelled: "Abgebrochen",
}

function StepIcon({ status }: { status: StepStatus }) {
	if (status === "completed") {
		return <CheckCircle2 className="h-4 w-4 text-emerald-500 fill-emerald-500/15" strokeWidth={2} />
	}
	if (status === "failed") return <XCircle className="h-4 w-4 text-destructive" />
	if (status === "running") return <Loader2 className="h-4 w-4 animate-spin text-primary" />
	return <Circle className="h-4 w-4 text-muted-foreground/40" />
}

function StepRow({ step }: { step: StepView }) {
	const showSub = step.name === "summarize_files" && step.fileCount && step.fileCount > 0
	return (
		<div
			className={`flex items-center gap-3 rounded-md px-3 py-2 transition-colors ${
				step.status === "running" ? "bg-primary/5" : ""
			}`}
		>
			<StepIcon status={step.status} />
			<div className="flex-1 min-w-0">
				<div className="text-sm">{stepLabel[step.name]}</div>
				{showSub && (
					<div className="text-xs text-muted-foreground truncate">
						{step.status === "completed"
							? `${step.fileCount} Datei${step.fileCount === 1 ? "" : "en"} verarbeitet`
							: `${step.fileIndex ?? 0} / ${step.fileCount}${
									step.filename ? ` — ${step.filename}` : ""
								}`}
					</div>
				)}
				{step.status === "failed" && step.reason && (
					<div className="text-xs text-destructive truncate" title={step.reason}>
						{step.reason}
					</div>
				)}
			</div>
			<div className="shrink-0 font-mono text-xs text-muted-foreground tabular-nums">
				{step.status === "completed" || step.status === "failed"
					? formatDuration(step.durationMs)
					: step.status === "running"
						? "…"
						: ""}
			</div>
		</div>
	)
}

export interface RunRowProps {
	id: number
	repoUrl: string
	status?: string
	fileCount?: number
	startedAt?: string
	updatedAt?: string
	isOpen: boolean
}

// RunRow renders one run as a shadcn AccordionItem. The parent
// <Accordion type="multiple"> drives open/close state — we only need to
// fetch the detail when the card is currently open. Multi-mode lets the
// user inspect several runs side-by-side.
export function RunRow({ id, repoUrl, status, fileCount, startedAt, isOpen }: RunRowProps) {
	const live = useAIProgress(isOpen ? id : null)
	const detail = useRepoSummary(isOpen ? id : null)

	const detailEnvelope = detail.data as
		| { data?: AiworkflowsInterfacesHttpRepoSummaryResponse }
		| undefined
	const result = detailEnvelope?.data

	const effectiveStatus =
		live.runStatus !== "pending"
			? live.runStatus
			: ((result?.status as string | undefined) ?? status ?? "pending")

	const steps = useMemo(() => {
		const merged = { ...live.steps }
		const persisted = result?.stepDurations ?? {}

		for (const name of STEP_ORDER) {
			const dur = persisted[name]
			if (typeof dur === "number" && merged[name].status === "pending") {
				merged[name] = { ...merged[name], status: "completed", durationMs: dur }
			} else if (
				typeof dur === "number" &&
				merged[name].status === "completed" &&
				!merged[name].durationMs
			) {
				merged[name] = { ...merged[name], durationMs: dur }
			}
		}

		const filesDone = (result?.files?.length ?? fileCount ?? 0) > 0
		const hasSummary = (result?.summary?.length ?? 0) > 0
		const isCompleted = effectiveStatus === "completed"
		const isFailed = effectiveStatus === "failed" || effectiveStatus === "cancelled"

		const ensure = (name: StepName, next: StepStatus) => {
			if (merged[name].status === "pending") merged[name] = { ...merged[name], status: next }
		}

		if (isCompleted) {
			for (const name of STEP_ORDER) ensure(name, "completed")
		} else {
			if (filesDone || hasSummary) {
				ensure("clone", "completed")
				ensure("traverse", "completed")
			}
			if (filesDone) ensure("summarize_files", "completed")
			if (hasSummary) ensure("aggregate", "completed")
		}

		if (isFailed) {
			for (const k of STEP_ORDER) {
				if (merged[k].status !== "completed") {
					merged[k] = { ...merged[k], status: "failed", reason: result?.failReason }
					break
				}
			}
		}

		if (effectiveStatus === "running") {
			for (const name of STEP_ORDER) {
				if (merged[name].status === "pending") {
					merged[name] = {
						...merged[name],
						status: "running",
						fileCount:
							name === "summarize_files"
								? (merged[name].fileCount ?? fileCount ?? undefined)
								: merged[name].fileCount,
					}
					break
				}
			}
		}

		return merged
	}, [live.steps, result, effectiveStatus, fileCount])

	const tone = statusTone(effectiveStatus)
	const doneCount = STEP_ORDER.filter((s) => steps[s].status === "completed").length

	// Sub-label under the title — packs the most useful summary into one line.
	const subLabel = (() => {
		if (effectiveStatus === "running") {
			const sf = steps.summarize_files
			if (sf.status === "running" && sf.fileCount) {
				return `${sf.fileIndex ?? 0} / ${sf.fileCount} Dateien · ${formatRelative(startedAt)}`
			}
			const stepIdx = STEP_ORDER.findIndex((s) => steps[s].status === "running")
			const stepName = stepIdx >= 0 ? stepLabel[STEP_ORDER[stepIdx]] : "Läuft"
			return `${stepName} · ${formatRelative(startedAt)}`
		}
		if (effectiveStatus === "completed") {
			const files = result?.files?.length ?? fileCount ?? 0
			return `${files} Datei${files === 1 ? "" : "en"} · ${formatRelative(startedAt)}`
		}
		return `${statusLabel[effectiveStatus] ?? effectiveStatus} · ${formatRelative(startedAt)}`
	})()

	return (
		<AccordionItem value={String(id)} className="rounded-xl border bg-card last:border-b">
			{/* `items-center` overrides shadcn's `items-start` so the chevron
			    lines up with the title row instead of floating at the top.
			    `[&>svg]:translate-y-0` cancels the +2px nudge baked into the
			    primitive's ChevronDownIcon. */}
			<AccordionTrigger className="items-center px-4 py-3 hover:no-underline data-[state=open]:bg-muted/30 [&>svg]:translate-y-0">
				<div className="flex flex-1 items-center gap-3 min-w-0">
					<div className="flex-1 min-w-0 text-left">
						<div className="font-medium truncate">{formatRepoLabel(repoUrl)}</div>
						<div className="text-xs text-muted-foreground truncate font-normal">{subLabel}</div>
					</div>
					<div className={`flex items-center gap-2 text-sm font-medium ${toneText[tone]}`}>
						<span className={`h-2 w-2 rounded-full ${toneDot[tone]}`} />
						{statusLabel[effectiveStatus] ?? effectiveStatus}
					</div>
					{/* Mini step-progress counter, much more honest than a half-baked
					    bar that depends on detail-fetch state. */}
					<div className="hidden shrink-0 text-xs font-mono tabular-nums text-muted-foreground sm:block">
						{doneCount} / {STEP_ORDER.length}
					</div>
				</div>
			</AccordionTrigger>
			<AccordionContent className="px-4 pb-4">
				<div className="space-y-4">
					<div className="divide-y rounded-lg border">
						{STEP_ORDER.map((name) => (
							<StepRow key={name} step={steps[name]} />
						))}
					</div>

					{result?.status === "completed" && result.summary && (
						<div>
							<h3 className="mb-2 text-xs font-semibold uppercase tracking-wide text-muted-foreground">
								Zusammenfassung
							</h3>
							<p className="whitespace-pre-wrap text-sm leading-relaxed">{result.summary}</p>
						</div>
					)}

					{result?.files && result.files.length > 0 && (
						<details className="rounded-lg border">
							<summary className="cursor-pointer p-3 text-sm font-medium">
								Dateien ({result.files.length})
							</summary>
							<div className="divide-y border-t">
								{result.files.map((f) => (
									<div key={f.filename} className="space-y-1 p-3">
										<div className="font-mono text-xs">{f.filename}</div>
										<div className="text-sm text-muted-foreground leading-relaxed">{f.summary}</div>
									</div>
								))}
							</div>
						</details>
					)}

					{result?.status === "failed" && result.failReason && (
						<div className="rounded-md border border-destructive/40 bg-destructive/5 p-3 text-sm text-destructive">
							{result.failReason}
						</div>
					)}

					{!result && (effectiveStatus === "pending" || effectiveStatus === "running") && (
						<p className="text-xs text-muted-foreground">
							Verbindung zum Run hergestellt — Schritte werden live aktualisiert.
						</p>
					)}
				</div>
			</AccordionContent>
		</AccordionItem>
	)
}
