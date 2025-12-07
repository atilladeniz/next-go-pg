"use client"

import { twoFactor } from "@shared/lib/auth-client"
import { Button } from "@shared/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@shared/ui/card"
import { Input } from "@shared/ui/input"
import { Label } from "@shared/ui/label"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@shared/ui/tabs"
import { KeyRound, Loader2, Mail, Shield } from "lucide-react"
import { useRouter } from "next/navigation"
import { useEffect, useRef, useState } from "react"

export default function TwoFactorPage() {
	const router = useRouter()
	const [code, setCode] = useState("")
	const [backupCode, setBackupCode] = useState("")
	const [trustDevice, setTrustDevice] = useState(true)
	const [isLoading, setIsLoading] = useState(false)
	const [error, setError] = useState<string | null>(null)
	const [otpSent, setOtpSent] = useState(false)
	const [activeTab, setActiveTab] = useState("totp")
	const inputRef = useRef<HTMLInputElement>(null)

	// biome-ignore lint/correctness/useExhaustiveDependencies: Intentionally refocus input when tab changes
	useEffect(() => {
		if (inputRef.current) {
			inputRef.current.focus()
		}
	}, [activeTab])

	const handleVerifyTotp = async () => {
		if (code.length !== 6) return
		setIsLoading(true)
		setError(null)

		try {
			const result = await twoFactor.verifyTotp({
				code,
				trustDevice,
			})

			if (result.error) {
				setError(result.error.message || "Ungültiger Code")
				setIsLoading(false)
				return
			}

			router.push("/dashboard")
		} catch {
			setError("Ein Fehler ist aufgetreten")
			setIsLoading(false)
		}
	}

	const handleVerifyBackupCode = async () => {
		if (!backupCode.trim()) return
		setIsLoading(true)
		setError(null)

		try {
			const result = await twoFactor.verifyBackupCode({
				code: backupCode.trim(),
				trustDevice,
			})

			if (result.error) {
				setError(result.error.message || "Ungültiger Backup-Code")
				setIsLoading(false)
				return
			}

			router.push("/dashboard")
		} catch {
			setError("Ein Fehler ist aufgetreten")
			setIsLoading(false)
		}
	}

	const handleSendOtp = async () => {
		setIsLoading(true)
		setError(null)

		try {
			const result = await twoFactor.sendOtp()

			if (result.error) {
				setError(result.error.message || "Fehler beim Senden")
				setIsLoading(false)
				return
			}

			setOtpSent(true)
			setIsLoading(false)
		} catch {
			setError("Ein Fehler ist aufgetreten")
			setIsLoading(false)
		}
	}

	const handleVerifyOtp = async () => {
		if (code.length !== 6) return
		setIsLoading(true)
		setError(null)

		try {
			const result = await twoFactor.verifyOtp({
				code,
				trustDevice,
			})

			if (result.error) {
				setError(result.error.message || "Ungültiger Code")
				setIsLoading(false)
				return
			}

			router.push("/dashboard")
		} catch {
			setError("Ein Fehler ist aufgetreten")
			setIsLoading(false)
		}
	}

	return (
		<div className="flex min-h-screen items-center justify-center bg-background p-4">
			<Card className="w-full max-w-md">
				<CardHeader className="text-center">
					<div className="mx-auto mb-4 flex h-12 w-12 items-center justify-center rounded-full bg-primary/10">
						<Shield className="h-6 w-6 text-primary" />
					</div>
					<CardTitle>Zwei-Faktor-Authentifizierung</CardTitle>
					<CardDescription>Bestätige deine Identität mit einem zusätzlichen Code</CardDescription>
				</CardHeader>
				<CardContent>
					<Tabs value={activeTab} onValueChange={setActiveTab}>
						<TabsList className="grid w-full grid-cols-3">
							<TabsTrigger value="totp">
								<KeyRound className="mr-2 h-4 w-4" />
								App
							</TabsTrigger>
							<TabsTrigger value="email">
								<Mail className="mr-2 h-4 w-4" />
								E-Mail
							</TabsTrigger>
							<TabsTrigger value="backup">Backup</TabsTrigger>
						</TabsList>

						{/* TOTP Tab */}
						<TabsContent value="totp" className="space-y-4">
							<div className="space-y-2">
								<Label htmlFor="totp-code">Code aus Authenticator-App</Label>
								<Input
									ref={inputRef}
									id="totp-code"
									value={code}
									onChange={(e) => setCode(e.target.value.replace(/\D/g, "").slice(0, 6))}
									placeholder="000000"
									className="text-center text-2xl tracking-widest"
									maxLength={6}
									onKeyDown={(e) => e.key === "Enter" && handleVerifyTotp()}
								/>
							</div>

							<div className="flex items-center gap-2">
								<input
									type="checkbox"
									id="trust-device-totp"
									checked={trustDevice}
									onChange={(e) => setTrustDevice(e.target.checked)}
									className="rounded"
								/>
								<Label htmlFor="trust-device-totp" className="text-sm font-normal">
									Diesem Gerät 30 Tage vertrauen
								</Label>
							</div>

							{error && <p className="text-sm text-destructive">{error}</p>}

							<Button
								className="w-full"
								onClick={handleVerifyTotp}
								disabled={code.length !== 6 || isLoading}
							>
								{isLoading && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
								Bestätigen
							</Button>
						</TabsContent>

						{/* Email OTP Tab */}
						<TabsContent value="email" className="space-y-4">
							{!otpSent ? (
								<>
									<p className="text-sm text-muted-foreground">
										Wir senden dir einen Einmal-Code per E-Mail.
									</p>
									<Button className="w-full" onClick={handleSendOtp} disabled={isLoading}>
										{isLoading && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
										Code per E-Mail senden
									</Button>
								</>
							) : (
								<>
									<div className="space-y-2">
										<Label htmlFor="email-code">Code aus E-Mail</Label>
										<Input
											id="email-code"
											value={code}
											onChange={(e) => setCode(e.target.value.replace(/\D/g, "").slice(0, 6))}
											placeholder="000000"
											className="text-center text-2xl tracking-widest"
											maxLength={6}
											onKeyDown={(e) => e.key === "Enter" && handleVerifyOtp()}
										/>
									</div>

									<div className="flex items-center gap-2">
										<input
											type="checkbox"
											id="trust-device-email"
											checked={trustDevice}
											onChange={(e) => setTrustDevice(e.target.checked)}
											className="rounded"
										/>
										<Label htmlFor="trust-device-email" className="text-sm font-normal">
											Diesem Gerät 30 Tage vertrauen
										</Label>
									</div>

									<Button
										className="w-full"
										onClick={handleVerifyOtp}
										disabled={code.length !== 6 || isLoading}
									>
										{isLoading && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
										Bestätigen
									</Button>

									<Button
										variant="ghost"
										className="w-full"
										onClick={() => {
											setOtpSent(false)
											setCode("")
										}}
									>
										Erneut senden
									</Button>
								</>
							)}

							{error && <p className="text-sm text-destructive">{error}</p>}
						</TabsContent>

						{/* Backup Code Tab */}
						<TabsContent value="backup" className="space-y-4">
							<div className="space-y-2">
								<Label htmlFor="backup-code">Backup-Code</Label>
								<Input
									id="backup-code"
									value={backupCode}
									onChange={(e) => setBackupCode(e.target.value)}
									placeholder="xxxxxxxxxx"
									className="font-mono"
								/>
								<p className="text-xs text-muted-foreground">
									Verwende einen deiner einmaligen Backup-Codes
								</p>
							</div>

							<div className="flex items-center gap-2">
								<input
									type="checkbox"
									id="trust-device-backup"
									checked={trustDevice}
									onChange={(e) => setTrustDevice(e.target.checked)}
									className="rounded"
								/>
								<Label htmlFor="trust-device-backup" className="text-sm font-normal">
									Diesem Gerät 30 Tage vertrauen
								</Label>
							</div>

							{error && <p className="text-sm text-destructive">{error}</p>}

							<Button
								className="w-full"
								onClick={handleVerifyBackupCode}
								disabled={!backupCode.trim() || isLoading}
							>
								{isLoading && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
								Bestätigen
							</Button>
						</TabsContent>
					</Tabs>
				</CardContent>
			</Card>
		</div>
	)
}
