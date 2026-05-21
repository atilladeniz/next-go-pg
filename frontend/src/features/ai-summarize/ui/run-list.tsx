"use client"

import { cn } from "@shared/lib/utils"
import { Accordion } from "@shared/ui/accordion"
import { Sparkles } from "lucide-react"
import { useRouter, useSearchParams } from "next/navigation"
import { useCallback, useEffect, useMemo, useState } from "react"
import { useHistory } from "../model/use-history"
import { RunRow } from "./run-row"

// RunList is the single source of truth for "runs on this account".
// shadcn Accordion in `type="multiple"` mode lets the user pin several
// cards open at once. The open set is mirrored in `?ids=1,2,3` so a
// refresh keeps them open AND the server-side page can prefetch each
// summary's detail BEFORE hydration.
export function RunList() {
	const router = useRouter()
	const params = useSearchParams()
	const { items, isLoading, isError } = useHistory()

	// `mounted` suppresses the accordion-content open/close animations on
	// the very first paint after a refresh. Without it Radix runs
	// `animate-accordion-down` on mount because its `data-state` flips
	// from undefined → "open", which renders as a visible close-then-open
	// flicker for any card that the URL says should be expanded.
	// Pattern documented in radix-ui/primitives#1463 and shadcn-ui#1644.
	const [mounted, setMounted] = useState(false)
	useEffect(() => {
		setMounted(true)
	}, [])

	const openIds = useMemo(() => {
		const single = params.get("id")
		const multi = params.get("ids")
		const raw = multi ?? single
		if (!raw) return []
		return raw
			.split(",")
			.map((s) => s.trim())
			.filter(Boolean)
	}, [params])

	const writeOpenIds = useCallback(
		(next: string[]) => {
			const url = next.length ? `/ai/summarize?ids=${next.join(",")}` : "/ai/summarize"
			router.replace(url, { scroll: false })
		},
		[router],
	)

	const onValueChange = useCallback(
		(next: string[]) => {
			writeOpenIds(next)
		},
		[writeOpenIds],
	)

	const onDeleted = useCallback(
		(id: number) => {
			writeOpenIds(openIds.filter((x) => x !== String(id)))
		},
		[openIds, writeOpenIds],
	)

	if (isLoading) {
		return <p className="text-sm text-muted-foreground">Lade Runs…</p>
	}

	if (isError) {
		return <p className="text-sm text-destructive">Konnte Runs nicht laden. Backend erreichbar?</p>
	}

	if (items.length === 0) {
		return (
			<div className="rounded-xl border border-dashed bg-muted/20 p-8 text-center">
				<Sparkles className="mx-auto h-8 w-8 text-muted-foreground/60" />
				<p className="mt-3 text-sm font-medium">Noch keine Zusammenfassung</p>
				<p className="mt-1 text-xs text-muted-foreground">
					URL oben eintippen und starten — jeder Run erscheint hier.
				</p>
			</div>
		)
	}

	return (
		<Accordion
			type="multiple"
			value={openIds}
			onValueChange={onValueChange}
			className={cn(
				"space-y-3",
				// First paint only — kill the open/close animation so Radix
				// does not flash the just-restored card closed-then-open.
				!mounted && "[&_[data-slot=accordion-content]]:!animate-none",
			)}
		>
			{items.map((it) => (
				<RunRow
					key={it.id}
					id={it.id ?? 0}
					repoUrl={it.repoUrl ?? ""}
					status={it.status}
					fileCount={it.fileCount}
					startedAt={it.startedAt}
					updatedAt={it.updatedAt}
					isOpen={openIds.includes(String(it.id))}
					onDeleted={onDeleted}
				/>
			))}
		</Accordion>
	)
}
