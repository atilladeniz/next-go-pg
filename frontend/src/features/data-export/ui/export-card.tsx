"use client"

import { Button } from "@shared/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@shared/ui/card"
import { Progress } from "@shared/ui/progress"
import { CheckCircle2, Download, FileJson, FileSpreadsheet, Loader2, XCircle } from "lucide-react"
import { useState } from "react"
import { type ExportDataType, type ExportFormat, useExport } from "../model/use-export"

const formatOptions: { value: ExportFormat; label: string; icon: typeof FileJson }[] = [
	{ value: "csv", label: "CSV", icon: FileSpreadsheet },
	{ value: "json", label: "JSON", icon: FileJson },
]

const dataTypeOptions: { value: ExportDataType; label: string; description: string }[] = [
	{ value: "stats", label: "Statistiken", description: "Projektübersicht der letzten 7 Tage" },
	{ value: "activity", label: "Aktivität", description: "Letzte Aktionen und Ereignisse" },
	{ value: "all", label: "Alles", description: "Vollständiger Daten-Export" },
]

export function ExportCard() {
	const [selectedFormat, setSelectedFormat] = useState<ExportFormat>("csv")
	const [selectedDataType, setSelectedDataType] = useState<ExportDataType>("all")

	const {
		progress,
		message,
		fileName,
		startExport,
		downloadExport,
		reset,
		isExporting,
		isCompleted,
		isFailed,
	} = useExport()

	const handleExport = () => {
		startExport(selectedFormat, selectedDataType)
	}

	return (
		<Card className="w-full">
			<CardHeader>
				<CardTitle className="flex items-center gap-2">
					<Download className="h-5 w-5" />
					Daten exportieren
				</CardTitle>
				<CardDescription>Exportiere deine Daten als CSV oder JSON Datei</CardDescription>
			</CardHeader>
			<CardContent className="space-y-6">
				{/* Format Selection */}
				<div className="space-y-2">
					<span className="text-sm font-medium">Format</span>
					<div className="flex gap-2">
						{formatOptions.map((option) => {
							const Icon = option.icon
							return (
								<Button
									key={option.value}
									variant={selectedFormat === option.value ? "default" : "outline"}
									className="flex-1"
									onClick={() => setSelectedFormat(option.value)}
									disabled={isExporting}
								>
									<Icon className="mr-2 h-4 w-4" />
									{option.label}
								</Button>
							)
						})}
					</div>
				</div>

				{/* Data Type Selection */}
				<div className="space-y-2">
					<span className="text-sm font-medium">Daten</span>
					<div className="grid gap-2">
						{dataTypeOptions.map((option) => (
							<button
								key={option.value}
								type="button"
								className={`rounded-lg border p-3 text-left transition-colors ${
									selectedDataType === option.value
										? "border-primary bg-primary/5"
										: "border-border hover:border-primary/50"
								} ${isExporting ? "opacity-50 cursor-not-allowed" : "cursor-pointer"}`}
								onClick={() => !isExporting && setSelectedDataType(option.value)}
								disabled={isExporting}
							>
								<div className="font-medium">{option.label}</div>
								<div className="text-sm text-muted-foreground">{option.description}</div>
							</button>
						))}
					</div>
				</div>

				{/* Progress Section */}
				{(isExporting || isCompleted || isFailed) && (
					<div className="space-y-3 rounded-lg border bg-muted/30 p-4">
						{/* Status Icon */}
						<div className="flex items-center gap-3">
							{isExporting && <Loader2 className="h-5 w-5 animate-spin text-primary" />}
							{isCompleted && <CheckCircle2 className="h-5 w-5 text-green-500" />}
							{isFailed && <XCircle className="h-5 w-5 text-destructive" />}
							<span className="font-medium">
								{isExporting && "Export läuft..."}
								{isCompleted && "Export abgeschlossen!"}
								{isFailed && "Export fehlgeschlagen"}
							</span>
						</div>

						{/* Progress Bar */}
						{isExporting && (
							<div className="space-y-2">
								<Progress value={progress} className="h-2" />
								<div className="flex justify-between text-sm text-muted-foreground">
									<span>{message}</span>
									<span>{progress}%</span>
								</div>
							</div>
						)}

						{/* Success Message */}
						{isCompleted && fileName && (
							<div className="text-sm text-muted-foreground">
								Datei: <span className="font-mono">{fileName}</span>
							</div>
						)}
					</div>
				)}

				{/* Action Buttons */}
				<div className="flex gap-2">
					{!isCompleted && !isFailed && (
						<Button className="flex-1" onClick={handleExport} disabled={isExporting}>
							{isExporting ? (
								<>
									<Loader2 className="mr-2 h-4 w-4 animate-spin" />
									Exportiere...
								</>
							) : (
								<>
									<Download className="mr-2 h-4 w-4" />
									Export starten
								</>
							)}
						</Button>
					)}

					{isCompleted && (
						<>
							<Button className="flex-1" onClick={downloadExport}>
								<Download className="mr-2 h-4 w-4" />
								Herunterladen
							</Button>
							<Button variant="outline" onClick={reset}>
								Neuer Export
							</Button>
						</>
					)}

					{isFailed && (
						<Button variant="outline" className="flex-1" onClick={reset}>
							Erneut versuchen
						</Button>
					)}
				</div>
			</CardContent>
		</Card>
	)
}
