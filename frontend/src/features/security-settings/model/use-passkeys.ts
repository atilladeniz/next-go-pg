"use client"

import { passkey } from "@shared/lib/auth-client"
import { useCallback, useEffect, useState } from "react"

interface Passkey {
	id: string
	name: string | null
	createdAt: Date
	deviceType: string
}

interface PasskeysState {
	passkeys: Passkey[]
	isLoading: boolean
	error: string | null
	isAdding: boolean
}

export function usePasskeys() {
	const [state, setState] = useState<PasskeysState>({
		passkeys: [],
		isLoading: true,
		error: null,
		isAdding: false,
	})

	const loadPasskeys = useCallback(async () => {
		setState((prev) => ({ ...prev, isLoading: true, error: null }))

		try {
			const result = await passkey.listUserPasskeys()

			if (result.error) {
				setState((prev) => ({
					...prev,
					isLoading: false,
					error: result.error.message || "Fehler beim Laden",
				}))
				return
			}

			setState((prev) => ({
				...prev,
				isLoading: false,
				passkeys: (result.data || []) as Passkey[],
			}))
		} catch {
			setState((prev) => ({
				...prev,
				isLoading: false,
				error: "Unbekannter Fehler",
			}))
		}
	}, [])

	useEffect(() => {
		loadPasskeys()
	}, [loadPasskeys])

	const addPasskey = useCallback(
		async (name?: string) => {
			setState((prev) => ({ ...prev, isAdding: true, error: null }))

			try {
				const result = await passkey.addPasskey({
					name: name || undefined,
				})

				if (result.error) {
					setState((prev) => ({
						...prev,
						isAdding: false,
						error: result.error.message || "Fehler beim Hinzufügen",
					}))
					return false
				}

				// Reload passkeys after adding
				await loadPasskeys()
				setState((prev) => ({ ...prev, isAdding: false }))
				return true
			} catch (err) {
				// WebAuthn errors
				const errorMessage =
					err instanceof Error ? err.message : "Passkey konnte nicht erstellt werden"
				setState((prev) => ({
					...prev,
					isAdding: false,
					error: errorMessage,
				}))
				return false
			}
		},
		[loadPasskeys],
	)

	const deletePasskey = useCallback(async (id: string) => {
		setState((prev) => ({ ...prev, isLoading: true, error: null }))

		try {
			const result = await passkey.deletePasskey({ id })

			if (result.error) {
				setState((prev) => ({
					...prev,
					isLoading: false,
					error: result.error.message || "Fehler beim Löschen",
				}))
				return false
			}

			// Remove from local state
			setState((prev) => ({
				...prev,
				isLoading: false,
				passkeys: prev.passkeys.filter((p) => p.id !== id),
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

	const updatePasskeyName = useCallback(async (id: string, name: string) => {
		try {
			const result = await passkey.updatePasskey({ id, name })

			if (result.error) {
				return false
			}

			// Update local state
			setState((prev) => ({
				...prev,
				passkeys: prev.passkeys.map((p) => (p.id === id ? { ...p, name } : p)),
			}))
			return true
		} catch {
			return false
		}
	}, [])

	return {
		...state,
		addPasskey,
		deletePasskey,
		updatePasskeyName,
		refresh: loadPasskeys,
	}
}
