"use client"

import { useLogin } from "../model/use-login"
import { EmailSentCard } from "./email-sent-card"
import { LoginCard } from "./login-card"

export function LoginForm() {
	const { email, setEmail, error, loading, sent, retryAfter, handleSubmit, reset } = useLogin()

	return (
		<div className="flex min-h-screen items-center justify-center bg-background">
			{sent ? (
				<EmailSentCard email={email} onReset={reset} />
			) : (
				<LoginCard
					email={email}
					error={error}
					loading={loading}
					retryAfter={retryAfter}
					onEmailChange={setEmail}
					onSubmit={handleSubmit}
				/>
			)}
		</div>
	)
}
