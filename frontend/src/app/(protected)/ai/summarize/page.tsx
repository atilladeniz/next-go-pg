import { SummarizeCard } from "@features/ai-summarize"
import { getSession } from "@shared/lib/auth-server"
import { Header } from "@widgets/header"
import { redirect } from "next/navigation"

export default async function AISummarizePage() {
	const session = await getSession()
	if (!session) redirect("/login")

	return (
		<div className="min-h-screen bg-background">
			<Header user={session.user} />
			<main className="mx-auto max-w-3xl px-4 py-8">
				<SummarizeCard />
			</main>
		</div>
	)
}
