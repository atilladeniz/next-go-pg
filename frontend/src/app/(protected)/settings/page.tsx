import { UserInfo } from "@entities/user"
import { SecurityCard } from "@features/security-settings"
import { SessionsList } from "@features/user-settings"
import { getSession } from "@shared/lib/auth-server"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@shared/ui/card"
import { Header } from "@widgets/header"
import { redirect } from "next/navigation"

export default async function SettingsPage() {
	const session = await getSession()
	if (!session) redirect("/login")

	const { user } = session

	return (
		<div className="min-h-screen bg-background">
			<Header user={user} />

			<main className="mx-auto max-w-3xl px-4 py-8 space-y-6">
				<h1 className="text-2xl font-bold">Einstellungen</h1>

				<Card>
					<CardHeader>
						<CardTitle>Profil</CardTitle>
						<CardDescription>Deine Kontoinformationen</CardDescription>
					</CardHeader>
					<CardContent>
						<UserInfo user={user} showEmail />
					</CardContent>
				</Card>

				<SecurityCard />

				<Card>
					<CardHeader>
						<CardTitle>Aktive Sessions</CardTitle>
						<CardDescription>Geräte und Browser, auf denen du angemeldet bist</CardDescription>
					</CardHeader>
					<CardContent>
						<SessionsList />
					</CardContent>
				</Card>
			</main>
		</div>
	)
}
