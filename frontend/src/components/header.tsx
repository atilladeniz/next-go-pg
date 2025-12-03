"use client"

import Link from "next/link"
import { usePathname, useRouter } from "next/navigation"
import { ModeToggle } from "@/components/mode-toggle"
import { Button } from "@/components/ui/button"
import { broadcastSignOut, useAuthSync } from "@/hooks/use-auth-sync"
import { signOut } from "@/lib/auth-client"
import { cn } from "@/lib/utils"

type User = {
	id: string
	name: string
	email: string
}

const navItems = [
	{ href: "/dashboard", label: "Dashboard" },
	{ href: "/api-test", label: "API Test" },
]

export function Header({ user }: { user: User }) {
	const router = useRouter()
	const pathname = usePathname()

	// Cross-Tab Auth Synchronisation
	useAuthSync()

	const handleSignOut = async () => {
		await signOut()
		await broadcastSignOut() // Alle anderen Tabs benachrichtigen
		router.push("/")
	}

	return (
		<header className="border-b">
			<div className="mx-auto flex max-w-7xl items-center justify-between px-4 py-4">
				<div className="flex items-center gap-6">
					<Link href="/dashboard" className="text-xl font-semibold">
						GoCa
					</Link>
					<nav className="flex items-center gap-1">
						{navItems.map((item) => (
							<Link
								key={item.href}
								href={item.href}
								className={cn(
									"rounded-md px-3 py-2 text-sm font-medium transition-colors",
									pathname === item.href
										? "bg-accent text-accent-foreground"
										: "text-muted-foreground hover:bg-accent hover:text-accent-foreground",
								)}
							>
								{item.label}
							</Link>
						))}
					</nav>
				</div>
				<div className="flex items-center gap-3">
					<span className="text-sm text-muted-foreground">{user.name || user.email}</span>
					<ModeToggle />
					<Button variant="outline" size="sm" onClick={handleSignOut}>
						Abmelden
					</Button>
				</div>
			</div>
		</header>
	)
}
