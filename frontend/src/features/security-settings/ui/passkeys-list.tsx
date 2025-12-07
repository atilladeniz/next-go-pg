"use client"

import { Button } from "@shared/ui/button"
import {
	Dialog,
	DialogContent,
	DialogDescription,
	DialogFooter,
	DialogHeader,
	DialogTitle,
	DialogTrigger,
} from "@shared/ui/dialog"
import { Input } from "@shared/ui/input"
import { Label } from "@shared/ui/label"
import { Key, Loader2, Plus, Smartphone, Trash2 } from "lucide-react"
import { useState } from "react"
import { usePasskeys } from "../model/use-passkeys"

export function PasskeysList() {
	const { passkeys, isLoading, error, isAdding, addPasskey, deletePasskey } = usePasskeys()
	const [newPasskeyName, setNewPasskeyName] = useState("")
	const [showAddDialog, setShowAddDialog] = useState(false)
	const [deleteConfirmId, setDeleteConfirmId] = useState<string | null>(null)

	const handleAdd = async () => {
		const success = await addPasskey(newPasskeyName || undefined)
		if (success) {
			setNewPasskeyName("")
			setShowAddDialog(false)
		}
	}

	const handleDelete = async (id: string) => {
		const success = await deletePasskey(id)
		if (success) {
			setDeleteConfirmId(null)
		}
	}

	const formatDate = (date: Date) => {
		return new Date(date).toLocaleDateString("de-DE", {
			day: "2-digit",
			month: "2-digit",
			year: "numeric",
		})
	}

	return (
		<div className="space-y-4">
			{/* Header */}
			<div className="flex items-center justify-between">
				<div className="flex items-center gap-3">
					<Key className="h-5 w-5 text-muted-foreground" />
					<div>
						<p className="font-medium">Passkeys</p>
						<p className="text-sm text-muted-foreground">
							Sichere Anmeldung mit Fingerabdruck, Gesichtserkennung oder Sicherheitsschlüssel
						</p>
					</div>
				</div>
				<Dialog open={showAddDialog} onOpenChange={setShowAddDialog}>
					<DialogTrigger asChild>
						<Button size="sm">
							<Plus className="mr-2 h-4 w-4" />
							Passkey hinzufügen
						</Button>
					</DialogTrigger>
					<DialogContent className="sm:max-w-md">
						<DialogHeader>
							<DialogTitle>Neuen Passkey hinzufügen</DialogTitle>
							<DialogDescription>
								Dein Gerät wird dich auffordern, einen Passkey zu erstellen. Du kannst
								Fingerabdruck, Gesichtserkennung oder einen Sicherheitsschlüssel verwenden.
							</DialogDescription>
						</DialogHeader>

						<div className="space-y-4">
							<div className="space-y-2">
								<Label htmlFor="passkey-name">Name (optional)</Label>
								<Input
									id="passkey-name"
									value={newPasskeyName}
									onChange={(e) => setNewPasskeyName(e.target.value)}
									placeholder="z.B. MacBook Pro, iPhone"
								/>
							</div>

							{error && <p className="text-sm text-destructive">{error}</p>}
						</div>

						<DialogFooter>
							<Button variant="outline" onClick={() => setShowAddDialog(false)}>
								Abbrechen
							</Button>
							<Button onClick={handleAdd} disabled={isAdding}>
								{isAdding && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
								Passkey erstellen
							</Button>
						</DialogFooter>
					</DialogContent>
				</Dialog>
			</div>

			{/* Error */}
			{error && !showAddDialog && <p className="text-sm text-destructive">{error}</p>}

			{/* Loading */}
			{isLoading && (
				<div className="flex items-center justify-center py-8">
					<Loader2 className="h-6 w-6 animate-spin text-muted-foreground" />
				</div>
			)}

			{/* Empty State */}
			{!isLoading && passkeys.length === 0 && (
				<div className="rounded-lg border border-dashed py-8 text-center">
					<Key className="mx-auto h-10 w-10 text-muted-foreground" />
					<p className="mt-2 text-sm text-muted-foreground">Noch keine Passkeys eingerichtet</p>
					<p className="text-xs text-muted-foreground">
						Passkeys sind die sicherste Methode, dich anzumelden
					</p>
				</div>
			)}

			{/* Passkey List */}
			{!isLoading && passkeys.length > 0 && (
				<div className="divide-y rounded-lg border">
					{passkeys.map((pk) => (
						<div key={pk.id} className="flex items-center justify-between p-4">
							<div className="flex items-center gap-3">
								<Smartphone className="h-5 w-5 text-muted-foreground" />
								<div>
									<p className="font-medium">{pk.name || "Unbenannter Passkey"}</p>
									<p className="text-sm text-muted-foreground">
										{pk.deviceType} &middot; Erstellt am {formatDate(pk.createdAt)}
									</p>
								</div>
							</div>
							<Dialog
								open={deleteConfirmId === pk.id}
								onOpenChange={(open) => !open && setDeleteConfirmId(null)}
							>
								<DialogTrigger asChild>
									<Button
										variant="ghost"
										size="icon"
										className="text-destructive hover:text-destructive"
										onClick={() => setDeleteConfirmId(pk.id)}
									>
										<Trash2 className="h-4 w-4" />
									</Button>
								</DialogTrigger>
								<DialogContent className="sm:max-w-md">
									<DialogHeader>
										<DialogTitle>Passkey entfernen?</DialogTitle>
										<DialogDescription>
											Möchtest du den Passkey &quot;{pk.name || "Unbenannt"}&quot; wirklich
											entfernen? Du kannst dich dann nicht mehr damit anmelden.
										</DialogDescription>
									</DialogHeader>
									<DialogFooter>
										<Button variant="outline" onClick={() => setDeleteConfirmId(null)}>
											Abbrechen
										</Button>
										<Button
											variant="destructive"
											onClick={() => handleDelete(pk.id)}
											disabled={isLoading}
										>
											{isLoading && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
											Entfernen
										</Button>
									</DialogFooter>
								</DialogContent>
							</Dialog>
						</div>
					))}
				</div>
			)}
		</div>
	)
}
