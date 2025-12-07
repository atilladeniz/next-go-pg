"use client"

import { twoFactor, useSession } from "@shared/lib/auth-client"
import { useCallback, useState } from "react"

interface TwoFactorState {
	isEnabled: boolean
	isLoading: boolean
	error: string | null
	totpUri: string | null
	backupCodes: string[] | null
	showSetup: boolean
	showBackupCodes: boolean
}

// Type for 2FA API calls without password (skipVerificationOnEnable: true)
type TwoFactorResult = {
	error?: { message?: string }
	data?: { totpURI?: string; backupCodes?: string[] }
}

export function useTwoFactor() {
	const { data: session } = useSession()
	const [state, setState] = useState<TwoFactorState>({
		isEnabled: session?.user?.twoFactorEnabled ?? false,
		isLoading: false,
		error: null,
		totpUri: null,
		backupCodes: null,
		showSetup: false,
		showBackupCodes: false,
	})

	const enable = useCallback(async () => {
		setState((prev) => ({ ...prev, isLoading: true, error: null }))

		try {
			// skipVerificationOnEnable is true, so password is not required
			const result = (await (twoFactor.enable as (arg: object) => Promise<unknown>)(
				{},
			)) as TwoFactorResult

			if (result.error) {
				const errorMessage = result.error.message || "Fehler beim Aktivieren"
				setState((prev) => ({
					...prev,
					isLoading: false,
					error: errorMessage,
				}))
				return
			}

			setState((prev) => ({
				...prev,
				isLoading: false,
				totpUri: result.data?.totpURI || null,
				backupCodes: result.data?.backupCodes || null,
				showSetup: true,
			}))
		} catch {
			setState((prev) => ({
				...prev,
				isLoading: false,
				error: "Unbekannter Fehler",
			}))
		}
	}, [])

	const verifyAndEnable = useCallback(async (code: string, trustDevice: boolean = false) => {
		setState((prev) => ({ ...prev, isLoading: true, error: null }))

		try {
			const result = await twoFactor.verifyTotp({
				code,
				trustDevice,
			})

			if (result.error) {
				setState((prev) => ({
					...prev,
					isLoading: false,
					error: result.error.message || "Ungültiger Code",
				}))
				return false
			}

			setState((prev) => ({
				...prev,
				isLoading: false,
				isEnabled: true,
				showSetup: false,
				showBackupCodes: true,
			}))
			return true
		} catch {
			setState((prev) => ({
				...prev,
				isLoading: false,
				error: "Unbekannter Fehler",
			}))
			return false
		}
	}, [])

	const disable = useCallback(async () => {
		setState((prev) => ({ ...prev, isLoading: true, error: null }))

		try {
			// skipVerificationOnEnable is true, so password is not required for disable either
			const result = (await (twoFactor.disable as (arg: object) => Promise<unknown>)(
				{},
			)) as TwoFactorResult

			if (result.error) {
				const errorMessage = result.error.message || "Fehler beim Deaktivieren"
				setState((prev) => ({
					...prev,
					isLoading: false,
					error: errorMessage,
				}))
				return false
			}

			setState((prev) => ({
				...prev,
				isLoading: false,
				isEnabled: false,
				totpUri: null,
				backupCodes: null,
			}))
			return true
		} catch {
			setState((prev) => ({
				...prev,
				isLoading: false,
				error: "Unbekannter Fehler",
			}))
			return false
		}
	}, [])

	const generateBackupCodes = useCallback(async () => {
		setState((prev) => ({ ...prev, isLoading: true, error: null }))

		try {
			// skipVerificationOnEnable is true, so password is not required
			const result = (await (twoFactor.generateBackupCodes as (arg: object) => Promise<unknown>)(
				{},
			)) as TwoFactorResult

			if (result.error) {
				const errorMessage = result.error.message || "Fehler beim Generieren"
				setState((prev) => ({
					...prev,
					isLoading: false,
					error: errorMessage,
				}))
				return null
			}

			const codes = result.data?.backupCodes || []
			setState((prev) => ({
				...prev,
				isLoading: false,
				backupCodes: codes,
				showBackupCodes: true,
			}))
			return codes
		} catch {
			setState((prev) => ({
				...prev,
				isLoading: false,
				error: "Unbekannter Fehler",
			}))
			return null
		}
	}, [])

	const closeSetup = useCallback(() => {
		setState((prev) => ({ ...prev, showSetup: false }))
	}, [])

	const closeBackupCodes = useCallback(() => {
		setState((prev) => ({ ...prev, showBackupCodes: false, backupCodes: null }))
	}, [])

	return {
		...state,
		enable,
		verifyAndEnable,
		disable,
		generateBackupCodes,
		closeSetup,
		closeBackupCodes,
	}
}
