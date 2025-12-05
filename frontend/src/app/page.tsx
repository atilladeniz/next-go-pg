import { Button } from "@shared/ui/button"
import Image from "next/image"
import Link from "next/link"

export default function Home() {
	return (
		<div className="flex min-h-screen items-center justify-center bg-background">
			<main className="flex min-h-screen w-full max-w-3xl flex-col items-center justify-between py-32 px-16 sm:items-start">
				<div className="flex items-center gap-4">
					<Image
						className="dark:invert"
						src="/next.svg"
						alt="Next.js logo"
						width={100}
						height={20}
						priority
					/>
					<span className="text-2xl text-muted-foreground">+</span>
					<Image src="/go.svg" alt="Go logo" width={80} height={30} priority />
				</div>
				<div className="flex flex-col items-center gap-6 text-center sm:items-start sm:text-left">
					<h1 className="max-w-xs text-3xl font-semibold leading-10 tracking-tight">
						Next.js + Go Backend mit Better Auth
					</h1>
					<p className="max-w-md text-lg leading-8 text-muted-foreground">
						Full-Stack Starter mit TypeScript, TanStack Query und PostgreSQL.
					</p>
				</div>
				<div className="flex flex-col gap-4 sm:flex-row">
					<Button asChild size="lg">
						<Link href="/login">Anmelden</Link>
					</Button>
					<Button asChild variant="outline" size="lg">
						<Link href="/register">Registrieren</Link>
					</Button>
				</div>
			</main>
		</div>
	)
}
