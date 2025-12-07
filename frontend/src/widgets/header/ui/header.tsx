"use client"

import type { SessionUser } from "@entities/user"
import { useAuthSync } from "@features/auth"
import { cn } from "@shared/lib"
import Link from "next/link"
import { usePathname } from "next/navigation"
import { ModeToggle } from "./mode-toggle"
import { UserMenu } from "./user-menu"

const navItems = [
	{ href: "/dashboard", label: "Dashboard" },
	{ href: "/api-test", label: "API Test" },
]

export function Header({ user }: { user: SessionUser }) {
	const pathname = usePathname()

	useAuthSync()

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
					<ModeToggle />
					<UserMenu user={user} />
				</div>
			</div>
		</header>
	)
}
