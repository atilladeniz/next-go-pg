"use client"

import * as AccordionPrimitive from "@radix-ui/react-accordion"
import { getGetAiSummariesQueryKey, useDeleteAiSummariesId } from "@shared/api/endpoints/ai/ai"
import type { AiworkflowsInterfacesHttpRepoSummaryResponse } from "@shared/api/models"
import { cn } from "@shared/lib/utils"
import { Accordion, AccordionContent, AccordionItem, AccordionTrigger } from "@shared/ui/accordion"
import {
	AlertDialog,
	AlertDialogAction,
	AlertDialogCancel,
	AlertDialogContent,
	AlertDialogDescription,
	AlertDialogFooter,
	AlertDialogHeader,
	AlertDialogTitle,
	AlertDialogTrigger,
} from "@shared/ui/alert-dialog"
import { Button } from "@shared/ui/button"
import { useQueryClient } from "@tanstack/react-query"
import { CheckCircle2, ChevronDown, Circle, Loader2, Trash2, XCircle } from "lucide-react"
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
			className={cn(
				"flex items-center gap-3 rounded-md px-3 py-2 transition-colors",
				step.status === "running" && "bg-primary/5",
			)}
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
	onDeleted?: (id: number) => void
}

// RunRow renders one run as a shadcn AccordionItem. The parent
// <Accordion type="multiple"> drives open/close state — we only need to
// fetch the detail when the card is currently open. Multi-mode lets the
// user inspect several runs side-by-side.
export function RunRow({
	id,
	repoUrl,
	status,
	fileCount,
	startedAt,
	isOpen,
	onDeleted,
}: RunRowProps) {
	// Always pass the real id so React Query's cache key stays stable.
	// Gating with `null` here would re-key the query on every collapse →
	// the previously fetched detail (incl. failure reason) would vanish
	// during the close animation, which the user sees as the failure
	// text dropping out before the panel finishes collapsing.
	const live = useAIProgress(id)
	const detail = useRepoSummary(id, isOpen)
	const queryClient = useQueryClient()
	const deleteMutation = useDeleteAiSummariesId()

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

	const handleConfirmDelete = async () => {
		try {
			await deleteMutation.mutateAsync({ id })
			queryClient.invalidateQueries({ queryKey: getGetAiSummariesQueryKey() })
			onDeleted?.(id)
		} catch {
			// Mutation surfaces the error via deleteMutation.isError;
			// the row stays visible so the user can retry.
		}
	}

	return (
		<AccordionItem
			value={String(id)}
			className={cn(
				"overflow-hidden rounded-xl border bg-card last:border-b",
				// 300ms with ease-out matches the accordion-down 200ms keyframe
				// closely enough that the bg fade looks coupled with the
				// height animation, without being so slow it feels laggy.
				"transition-colors duration-300 ease-out",
				"data-[state=open]:bg-muted/30",
			)}
		>
			{/* Header is a single flex row with `items-center` — the
			    Trigger expands to fill, the delete Button sits as a
			    proper flex sibling on the right. No absolute positioning,
			    no margin hacks: alignment is purely from flex centering,
			    so the button stays put whether the item is open or
			    closed. shadcn's <AccordionTrigger> wraps the trigger in
			    its own Header internally and can't be extended from
			    outside — we use the Radix primitive directly here.
			    The open-state background lives on the AccordionItem so
			    Trigger + delete Button visually share the same header
			    block; otherwise the Trigger looked like a separate inner
			    pill detached from the action button on its right. */}
			<AccordionPrimitive.Header className="flex items-center">
				<AccordionPrimitive.Trigger
					className={cn(
						"focus-visible:border-ring focus-visible:ring-ring/50",
						"flex flex-1 cursor-pointer items-center gap-3",
						"py-3 pl-4 pr-2 text-left text-sm font-medium",
						"outline-none focus-visible:ring-[3px]",
						"disabled:pointer-events-none disabled:opacity-50",
						"[&[data-state=open]>svg.chevron]:rotate-180",
					)}
				>
					<div className="min-w-0 flex-1 text-left">
						<div className="truncate font-medium">{formatRepoLabel(repoUrl)}</div>
						<div className="truncate text-xs font-normal text-muted-foreground">{subLabel}</div>
					</div>
					<div className={cn("flex items-center gap-2 text-sm font-medium", toneText[tone])}>
						<span className={cn("h-2 w-2 rounded-full", toneDot[tone])} />
						{statusLabel[effectiveStatus] ?? effectiveStatus}
					</div>
					{/* Mini step-progress counter, more honest than a half-baked
					    bar that depends on detail-fetch state. */}
					<div className="hidden shrink-0 font-mono text-xs tabular-nums text-muted-foreground sm:block">
						{doneCount} / {STEP_ORDER.length}
					</div>
					<ChevronDown className="chevron pointer-events-none size-4 shrink-0 text-muted-foreground transition-transform duration-200" />
				</AccordionPrimitive.Trigger>
				<AlertDialog>
					<AlertDialogTrigger asChild>
						<Button
							variant="ghost"
							size="icon-sm"
							disabled={deleteMutation.isPending}
							title="Run löschen"
							className="mr-2 text-muted-foreground hover:bg-destructive/10 hover:text-destructive"
						>
							<Trash2 className="h-4 w-4" />
							<span className="sr-only">Löschen</span>
						</Button>
					</AlertDialogTrigger>
					<AlertDialogContent>
						<AlertDialogHeader>
							<AlertDialogTitle>Run löschen?</AlertDialogTitle>
							<AlertDialogDescription>
								<span className="font-mono text-foreground">{formatRepoLabel(repoUrl)}</span> wird
								inklusive Datei-Zusammenfassungen unwiderruflich gelöscht.
							</AlertDialogDescription>
						</AlertDialogHeader>
						<AlertDialogFooter>
							<AlertDialogCancel>Abbrechen</AlertDialogCancel>
							<AlertDialogAction
								onClick={handleConfirmDelete}
								className="bg-destructive text-white hover:bg-destructive/90"
							>
								Löschen
							</AlertDialogAction>
						</AlertDialogFooter>
					</AlertDialogContent>
				</AlertDialog>
			</AccordionPrimitive.Header>
			<AccordionContent className="px-4 pb-4">
				<div className="space-y-4">
					{/* Inner panels keep `bg-card` so they pop against the muted
					    header background when the item is open. */}
					<div className="divide-y rounded-lg border bg-card">
						{STEP_ORDER.map((name) => (
							<StepRow key={name} step={steps[name]} />
						))}
					</div>

					{result?.status === "completed" && result.summary && (
						<div className="rounded-lg border bg-card p-4">
							<h3 className="mb-2 text-xs font-semibold uppercase tracking-wide text-muted-foreground">
								Zusammenfassung
							</h3>
							<p className="whitespace-pre-wrap text-sm leading-relaxed">{result.summary}</p>
						</div>
					)}

					{result?.files && result.files.length > 0 && (
						<Accordion
							type="single"
							collapsible
							className="overflow-hidden rounded-lg border bg-card"
						>
							<AccordionItem value="files" className="border-b-0">
								<AccordionTrigger className="items-center px-3 py-2 text-sm font-medium hover:no-underline data-[state=open]:border-b">
									Dateien ({result.files.length})
								</AccordionTrigger>
								<AccordionContent className="p-0">
									<div className="divide-y">
										{result.files.map((f) => (
											<div key={f.filename} className="space-y-1 p-3">
												<div className="font-mono text-xs">{f.filename}</div>
												<div className="text-sm leading-relaxed text-muted-foreground">
													{f.summary}
												</div>
											</div>
										))}
									</div>
								</AccordionContent>
							</AccordionItem>
						</Accordion>
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
