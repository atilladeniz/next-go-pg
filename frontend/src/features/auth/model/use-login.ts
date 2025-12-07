"use client"

import { signIn } from "@shared/lib/auth-client"
import { useCallback, useEffect, useRef, useState } from "react"

export function useLogin() {
	const [email, setEmail] = useState("")
	const [error, setError] = useState("")
	const [loading, setLoading] = useState(false)
	const [sent, setSent] = useState(false)
	const [retryAfter, setRetryAfter] = useState(0)
	// Track if onError already set an error to avoid overwriting
	const errorSetByOnError = useRef(false)

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
			errorSetByOnError.current = false

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
							errorSetByOnError.current = true
						}
					},
				},
			)

			if (result.error) {
				// Only set error if onError didn't already set a more specific one (e.g., rate limit)
				if (!errorSetByOnError.current) {
					setError(result.error.message || "Anmeldung fehlgeschlagen")
				}
				setLoading(false)
				return
			}

			setSent(true)
			setLoading(false)
		},
		[email],
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
