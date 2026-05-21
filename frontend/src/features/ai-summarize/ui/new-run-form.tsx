"use client"

import { getGetAiSummariesQueryKey } from "@shared/api/endpoints/ai/ai"
import { Button } from "@shared/ui/button"
import { Card, CardContent } from "@shared/ui/card"
import { Input } from "@shared/ui/input"
import { useQueryClient } from "@tanstack/react-query"
import { FileSearch, Loader2, XCircle } from "lucide-react"
import { useRouter, useSearchParams } from "next/navigation"
import { useState } from "react"
import { useSummarizeRepo } from "../model/use-summarize"

// NewRunForm is the always-visible header — paste a repo URL, hit submit,
// the new run appears as a card in the list below and auto-expands.
export function NewRunForm() {
	const router = useRouter()
	const params = useSearchParams()
	const queryClient = useQueryClient()
	const [repoUrl, setRepoUrl] = useState("")
	const mutation = useSummarizeRepo()

	const handleSubmit = async (e: React.FormEvent) => {
		e.preventDefault()
		try {
			const response = (await mutation.mutateAsync({ data: { repoUrl } })) as {
				data?: { summaryId?: number }
			}
			const id = response?.data?.summaryId
			if (typeof id === "number") {
				setRepoUrl("")
				// Merge the new id into the open set so the user keeps any
				// previously-pinned runs expanded too.
				const existing = (params.get("ids") ?? params.get("id") ?? "")
					.split(",")
					.map((s) => s.trim())
					.filter(Boolean)
				const next = Array.from(new Set([String(id), ...existing]))
				router.replace(`/ai/summarize?ids=${next.join(",")}`, { scroll: false })
				queryClient.invalidateQueries({ queryKey: getGetAiSummariesQueryKey() })
			}
		} catch (err) {
			console.error("[ai-summarize] mutate failed", err)
		}
	}

	return (
		<Card>
			<CardContent className="p-4">
				<form onSubmit={handleSubmit} className="flex flex-col gap-3 sm:flex-row">
					<Input
						type="url"
						placeholder="https://github.com/owner/repo"
						value={repoUrl}
						onChange={(e) => setRepoUrl(e.target.value)}
						required
						disabled={mutation.isPending}
						className="flex-1"
					/>
					<Button type="submit" disabled={mutation.isPending || !repoUrl} className="sm:w-44">
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
					<div className="mt-3 flex items-center gap-2 rounded-md border border-destructive/40 bg-destructive/5 p-2 text-sm text-destructive">
						<XCircle className="h-4 w-4" />
						<span>{(mutation.error as Error | null)?.message ?? "Anfrage fehlgeschlagen"}</span>
					</div>
				)}

				<p className="mt-2 text-xs text-muted-foreground">
					Öffentliche http(s)-URLs. Mehrere Runs parallel möglich — jeder erscheint als eigene Karte
					unten und aktualisiert sich live.
				</p>
			</CardContent>
		</Card>
	)
}
