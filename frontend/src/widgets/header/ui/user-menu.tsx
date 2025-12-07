"use client"

import type { SessionUser } from "@entities/user"
import { broadcastSignOut } from "@features/auth"
import { signOut } from "@shared/lib/auth-client"
import { Avatar, AvatarFallback } from "@shared/ui/avatar"
import { Button } from "@shared/ui/button"
import {
	DropdownMenu,
	DropdownMenuContent,
	DropdownMenuItem,
	DropdownMenuLabel,
	DropdownMenuSeparator,
	DropdownMenuTrigger,
} from "@shared/ui/dropdown-menu"
import { LogOut, Settings, User } from "lucide-react"
import Link from "next/link"
import { useRouter } from "next/navigation"

interface UserMenuProps {
	user: SessionUser
}

function getInitials(name: string | null, email: string): string {
	if (name) {
		return name
			.split(" ")
			.map((n) => n[0])
			.join("")
			.toUpperCase()
			.slice(0, 2)
	}
	return email.slice(0, 2).toUpperCase()
}

export function UserMenu({ user }: UserMenuProps) {
	const router = useRouter()

	const handleSignOut = async () => {
		await signOut()
		await broadcastSignOut()
		router.push("/")
	}

	return (
		<DropdownMenu>
			<DropdownMenuTrigger asChild>
				<Button variant="ghost" className="relative h-9 w-9 rounded-full">
					<Avatar className="h-9 w-9">
						<AvatarFallback>{getInitials(user.name, user.email)}</AvatarFallback>
					</Avatar>
				</Button>
			</DropdownMenuTrigger>
			<DropdownMenuContent className="w-56" align="end" forceMount>
				<DropdownMenuLabel className="font-normal">
					<div className="flex flex-col space-y-1">
						<p className="text-sm font-medium leading-none">{user.name || "User"}</p>
						<p className="text-xs leading-none text-muted-foreground">{user.email}</p>
					</div>
				</DropdownMenuLabel>
				<DropdownMenuSeparator />
				<DropdownMenuItem asChild>
					<Link href="/settings" className="flex items-center">
						<User className="mr-2 h-4 w-4" />
						Profil
					</Link>
				</DropdownMenuItem>
				<DropdownMenuItem asChild>
					<Link href="/settings" className="flex items-center">
						<Settings className="mr-2 h-4 w-4" />
						Einstellungen
					</Link>
				</DropdownMenuItem>
				<DropdownMenuSeparator />
				<DropdownMenuItem
					onClick={handleSignOut}
					className="text-destructive focus:text-destructive"
				>
					<LogOut className="mr-2 h-4 w-4" />
					Abmelden
				</DropdownMenuItem>
			</DropdownMenuContent>
		</DropdownMenu>
	)
}
