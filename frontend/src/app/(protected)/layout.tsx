import { redirect } from "next/navigation"
import { SessionProvider } from "@/components/session-provider"
import { getSession } from "@/lib/auth-server"

export default async function ProtectedLayout({ children }: { children: React.ReactNode }) {
	const session = await getSession()

	if (!session) {
		redirect("/login")
	}

	return <SessionProvider session={session}>{children}</SessionProvider>
}
