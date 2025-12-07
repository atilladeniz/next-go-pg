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

/**
 * Two-factor authentication API response types.
 * Note: With skipVerificationOnEnable: true, password is not required.
 * These types match the actual runtime responses from Better Auth.
 */
interface TwoFactorEnableResponse {
	error?: { message?: string }
	data?: { totpURI?: string; backupCodes?: string[] }
}

interface TwoFactorDisableResponse {
	error?: { message?: string }
	data?: unknown
}

interface TwoFactorBackupCodesResponse {
	error?: { message?: string }
	data?: { backupCodes?: string[] }
}

/**
 * Helper to call twoFactor methods without password parameter.
 * skipVerificationOnEnable: true is configured in auth.ts, so password is not needed.
 * We use explicit typing instead of type assertions on the method signatures.
 */
async function enableTwoFactor(): Promise<TwoFactorEnableResponse> {
	// @ts-expect-error - skipVerificationOnEnable makes password optional at runtime
	return twoFactor.enable({})
}

async function disableTwoFactor(): Promise<TwoFactorDisableResponse> {
	// @ts-expect-error - skipVerificationOnEnable makes password optional at runtime
	return twoFactor.disable({})
}

async function generateNewBackupCodes(): Promise<TwoFactorBackupCodesResponse> {
	// @ts-expect-error - skipVerificationOnEnable makes password optional at runtime
	return twoFactor.generateBackupCodes({})
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
			const result = await enableTwoFactor()

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
				const errorMessage = result.error.message || "Ungültiger Code"
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
			const result = await disableTwoFactor()

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
			const result = await generateNewBackupCodes()

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
