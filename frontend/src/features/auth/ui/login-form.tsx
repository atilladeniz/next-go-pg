"use client"

import { useSearchParams } from "next/navigation"
import { useEffect, useState } from "react"
import { useLogin } from "../model/use-login"
import { EmailSentCard } from "./email-sent-card"
import { LoginCard } from "./login-card"

const errorMessages: Record<string, string> = {
	verification_failed: "Der Link ist ungültig oder abgelaufen. Bitte fordere einen neuen an.",
	invalid_token: "Ungültiger Anmelde-Link.",
	expired: "Der Link ist abgelaufen. Bitte fordere einen neuen an.",
}

export function LoginForm() {
	const { email, setEmail, error, loading, sent, retryAfter, handleSubmit, reset } = useLogin()
	const searchParams = useSearchParams()
	const [urlError, setUrlError] = useState<string | null>(null)

	useEffect(() => {
		const errorParam = searchParams.get("error")
		if (errorParam) {
			setUrlError(errorMessages[errorParam] || "Ein Fehler ist aufgetreten.")
		}
	}, [searchParams])

	const displayError = urlError || error

	return (
		<div className="flex min-h-screen items-center justify-center bg-background">
			{sent ? (
				<EmailSentCard email={email} onReset={reset} />
			) : (
				<LoginCard
					email={email}
					error={displayError}
					loading={loading}
					retryAfter={retryAfter}
					onEmailChange={setEmail}
					onSubmit={handleSubmit}
				/>
			)}
		</div>
	)
}
