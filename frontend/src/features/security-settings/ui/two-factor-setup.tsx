"use client"

import { Button } from "@shared/ui/button"
import {
	Dialog,
	DialogContent,
	DialogDescription,
	DialogFooter,
	DialogHeader,
	DialogTitle,
} from "@shared/ui/dialog"
import { Input } from "@shared/ui/input"
import { Label } from "@shared/ui/label"
import { CheckCircle, Copy, Loader2, Shield } from "lucide-react"
import { useEffect, useRef, useState } from "react"
import { useTwoFactor } from "../model/use-two-factor"

export function TwoFactorSetup() {
	const {
		isEnabled,
		isLoading,
		error,
		totpUri,
		backupCodes,
		showSetup,
		showBackupCodes,
		enable,
		verifyAndEnable,
		disable,
		generateBackupCodes,
		closeSetup,
		closeBackupCodes,
	} = useTwoFactor()

	const [code, setCode] = useState("")
	const [trustDevice, setTrustDevice] = useState(true)
	const [copied, setCopied] = useState(false)
	const [qrLoaded, setQrLoaded] = useState(false)
	const inputRef = useRef<HTMLInputElement>(null)

	useEffect(() => {
		if (showSetup && inputRef.current) {
			inputRef.current.focus()
		}
	}, [showSetup])

	const handleVerify = async () => {
		if (code.length !== 6) return
		const success = await verifyAndEnable(code, trustDevice)
		if (success) {
			setCode("")
		}
	}

	const handleCopyBackupCodes = () => {
		if (!backupCodes) return
		navigator.clipboard.writeText(backupCodes.join("\n"))
		setCopied(true)
		setTimeout(() => setCopied(false), 2000)
	}

	const qrCodeUrl = totpUri
		? `https://api.qrserver.com/v1/create-qr-code/?size=200x200&data=${encodeURIComponent(totpUri)}`
		: null

	return (
		<div className="space-y-4">
			{/* Status */}
			<div className="flex items-center justify-between">
				<div className="flex items-center gap-3">
					<Shield className={`h-5 w-5 ${isEnabled ? "text-green-500" : "text-muted-foreground"}`} />
					<div>
						<p className="font-medium">Authenticator-App (TOTP)</p>
						<p className="text-sm text-muted-foreground">
							{isEnabled
								? "Aktiviert - Du wirst bei der Anmeldung nach einem Code gefragt"
								: "Nicht aktiviert - Schütze dein Konto mit einer Authenticator-App"}
						</p>
					</div>
				</div>
				<Button
					variant={isEnabled ? "outline" : "default"}
					onClick={isEnabled ? disable : enable}
					disabled={isLoading}
				>
					{isLoading && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
					{isEnabled ? "Deaktivieren" : "Aktivieren"}
				</Button>
			</div>

			{/* Backup Codes regenerieren */}
			{isEnabled && (
				<div className="flex items-center justify-between border-t pt-4">
					<div>
						<p className="text-sm font-medium">Backup-Codes</p>
						<p className="text-sm text-muted-foreground">
							Neue Backup-Codes generieren (alte werden ungültig)
						</p>
					</div>
					<Button variant="outline" size="sm" onClick={generateBackupCodes} disabled={isLoading}>
						Neue Codes generieren
					</Button>
				</div>
			)}

			{error && <p className="text-sm text-destructive">{error}</p>}

			{/* Setup Dialog */}
			<Dialog open={showSetup} onOpenChange={(open) => !open && closeSetup()}>
				<DialogContent className="sm:max-w-md">
					<DialogHeader>
						<DialogTitle>Authenticator-App einrichten</DialogTitle>
						<DialogDescription>
							Scanne den QR-Code mit deiner Authenticator-App (z.B. Google Authenticator, Authy)
						</DialogDescription>
					</DialogHeader>

					<div className="space-y-4">
						{/* QR Code */}
						{qrCodeUrl && (
							<div className="flex justify-center">
								<div className="rounded-lg border bg-white p-4">
									{!qrLoaded && (
										<div className="flex h-[200px] w-[200px] items-center justify-center">
											<Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
										</div>
									)}
									{/* biome-ignore lint/performance/noImgElement: External QR API cannot use next/image */}
									<img
										src={qrCodeUrl}
										alt="TOTP QR Code"
										width={200}
										height={200}
										onLoad={() => setQrLoaded(true)}
										className={qrLoaded ? "" : "hidden"}
									/>
								</div>
							</div>
						)}

						{/* Code Input */}
						<div className="space-y-2">
							<Label htmlFor="totp-code">6-stelliger Code</Label>
							<Input
								ref={inputRef}
								id="totp-code"
								value={code}
								onChange={(e) => setCode(e.target.value.replace(/\D/g, "").slice(0, 6))}
								placeholder="000000"
								className="text-center text-2xl tracking-widest"
								maxLength={6}
								onKeyDown={(e) => e.key === "Enter" && handleVerify()}
							/>
						</div>

						{/* Trust Device */}
						<div className="flex items-center gap-2">
							<input
								type="checkbox"
								id="trust-device"
								checked={trustDevice}
								onChange={(e) => setTrustDevice(e.target.checked)}
								className="rounded"
							/>
							<Label htmlFor="trust-device" className="text-sm font-normal">
								Diesem Gerät 30 Tage vertrauen
							</Label>
						</div>

						{error && <p className="text-sm text-destructive">{error}</p>}
					</div>

					<DialogFooter>
						<Button variant="outline" onClick={closeSetup}>
							Abbrechen
						</Button>
						<Button onClick={handleVerify} disabled={code.length !== 6 || isLoading}>
							{isLoading && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
							Bestätigen
						</Button>
					</DialogFooter>
				</DialogContent>
			</Dialog>

			{/* Backup Codes Dialog */}
			<Dialog open={showBackupCodes} onOpenChange={(open) => !open && closeBackupCodes()}>
				<DialogContent className="sm:max-w-md">
					<DialogHeader>
						<DialogTitle className="flex items-center gap-2">
							<CheckCircle className="h-5 w-5 text-green-500" />
							2FA aktiviert!
						</DialogTitle>
						<DialogDescription>
							Speichere diese Backup-Codes an einem sicheren Ort. Du kannst sie verwenden, falls du
							keinen Zugriff auf deine Authenticator-App hast.
						</DialogDescription>
					</DialogHeader>

					{backupCodes && (
						<div className="space-y-4">
							<div className="rounded-lg border bg-muted p-4">
								<div className="grid grid-cols-2 gap-2 font-mono text-sm">
									{backupCodes.map((code) => (
										<div key={code} className="rounded bg-background px-2 py-1">
											{code}
										</div>
									))}
								</div>
							</div>

							<Button variant="outline" className="w-full" onClick={handleCopyBackupCodes}>
								<Copy className="mr-2 h-4 w-4" />
								{copied ? "Kopiert!" : "Alle Codes kopieren"}
							</Button>

							<p className="text-xs text-muted-foreground">
								Jeder Code kann nur einmal verwendet werden. Nach der Verwendung wird er ungültig.
							</p>
						</div>
					)}

					<DialogFooter>
						<Button onClick={closeBackupCodes}>Verstanden</Button>
					</DialogFooter>
				</DialogContent>
			</Dialog>
		</div>
	)
}
