"use client"

import { useGetStats, usePostStats } from "@/api/endpoints/users/users"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { useSSE } from "@/hooks/use-sse"

export function StatsGrid() {
	// Connect to SSE for real-time updates
	useSSE()

	// Data is already hydrated from server via HydrationBoundary
	const { data: statsResponse } = useGetStats()

	// Generated mutation hook from Orval
	const { mutate: updateStats } = usePostStats()

	const stats = statsResponse?.status === 200 ? statsResponse.data : null

	const handleUpdate = (field: string, delta: number) => {
		updateStats({ data: { field, delta } })
	}

	return (
		<div className="mt-6 grid gap-6 md:grid-cols-2 lg:grid-cols-3">
			<StatsCard
				title="Projekte"
				value={stats?.projectCount ?? 0}
				label="Aktive Projekte"
				field="projects"
				onUpdate={handleUpdate}
			/>
			<StatsCard
				title="Aktivität"
				value={stats?.activityToday ?? 0}
				label="Heute"
				valueClassName="text-green-600"
				field="activity"
				onUpdate={handleUpdate}
			/>
			<StatsCard
				title="Benachrichtigungen"
				value={stats?.notifications ?? 0}
				label="Ungelesen"
				valueClassName="text-orange-600"
				field="notifications"
				onUpdate={handleUpdate}
			/>
		</div>
	)
}

function StatsCard({
	title,
	value,
	label,
	field,
	onUpdate,
	valueClassName = "text-primary",
}: {
	title: string
	value: number
	label: string
	field: string
	onUpdate: (field: string, delta: number) => void
	valueClassName?: string
}) {
	return (
		<Card>
			<CardHeader>
				<CardTitle className="text-base">{title}</CardTitle>
			</CardHeader>
			<CardContent>
				<div className="flex items-center gap-3">
					<Button
						variant="outline"
						size="icon"
						className="h-8 w-8"
						onClick={() => onUpdate(field, -1)}
					>
						-
					</Button>
					<p className={`text-3xl font-bold ${valueClassName} min-w-[3ch] text-center`}>{value}</p>
					<Button
						variant="outline"
						size="icon"
						className="h-8 w-8"
						onClick={() => onUpdate(field, 1)}
					>
						+
					</Button>
				</div>
				<p className="mt-2 text-sm text-muted-foreground">{label}</p>
			</CardContent>
		</Card>
	)
}
