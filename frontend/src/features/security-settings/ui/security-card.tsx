"use client"

import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@shared/ui/card"
import { Shield } from "lucide-react"
import { PasskeysList } from "./passkeys-list"
import { TwoFactorSetup } from "./two-factor-setup"

export function SecurityCard() {
	return (
		<Card>
			<CardHeader>
				<CardTitle className="flex items-center gap-2">
					<Shield className="h-5 w-5" />
					Sicherheit
				</CardTitle>
				<CardDescription>Schütze dein Konto mit zusätzlichen Sicherheitsmethoden</CardDescription>
			</CardHeader>
			<CardContent className="space-y-6">
				<TwoFactorSetup />
				<div className="border-t pt-6">
					<PasskeysList />
				</div>
			</CardContent>
		</Card>
	)
}
