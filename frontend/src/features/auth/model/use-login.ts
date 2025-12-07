"use client"

import { signIn } from "@shared/lib/auth-client"
import { useCallback, useEffect, useState } from "react"

export function useLogin() {
	const [email, setEmail] = useState("")
	const [error, setError] = useState("")
	const [loading, setLoading] = useState(false)
	const [sent, setSent] = useState(false)
	const [retryAfter, setRetryAfter] = useState(0)

	useEffect(() => {
		if (retryAfter <= 0) return
		const timer = setInterval(() => {
			setRetryAfter((prev) => Math.max(0, prev - 1))
		}, 1000)
		return () => clearInterval(timer)
	}, [retryAfter])

	const handleSubmit = useCallback(
		async (e: React.FormEvent) => {
			e.preventDefault()
			setError("")
			setLoading(true)

			const result = await signIn.magicLink(
				{
					email,
					callbackURL: "/dashboard",
					newUserCallbackURL: "/dashboard",
					errorCallbackURL: "/login?error=verification_failed",
				},
				{
					onError: (ctx) => {
						if (ctx.error.status === 429) {
							const retry = ctx.response?.headers?.get("X-Retry-After")
							if (retry) {
								setRetryAfter(Number.parseInt(retry, 10))
								setError(`Zu viele Anfragen. Bitte warte ${retry} Sekunden.`)
							} else {
								setError("Zu viele Anfragen. Bitte versuche es später erneut.")
							}
						}
					},
				},
			)

			if (result.error) {
				if (!error) {
					setError(result.error.message || "Anmeldung fehlgeschlagen")
				}
				setLoading(false)
				return
			}

			setSent(true)
			setLoading(false)
		},
		[email, error],
	)

	const reset = useCallback(() => {
		setSent(false)
		setEmail("")
		setError("")
		setRetryAfter(0)
	}, [])

	return {
		email,
		setEmail,
		error,
		loading,
		sent,
		retryAfter,
		handleSubmit,
		reset,
	}
}
