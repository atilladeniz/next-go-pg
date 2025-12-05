import { getSession } from "@shared/lib/auth-server"
import { Header } from "@widgets/header"
import { redirect } from "next/navigation"
import { ApiTestClient } from "./client"

export default async function ApiTestPage() {
	const session = await getSession()
	if (!session) redirect("/login")

	const { user } = session

	return (
		<div className="min-h-screen bg-background">
			<Header user={user} />
			<ApiTestClient userEmail={user.email} />
		</div>
	)
}
