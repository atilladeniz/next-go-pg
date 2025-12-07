import { Card, CardDescription, CardHeader, CardTitle } from "@shared/ui/card"
import { Loader2 } from "lucide-react"
import { Suspense } from "react"
import { MagicLinkVerifyContent } from "./verify-content"

function LoadingFallback() {
	return (
		<div className="flex min-h-screen items-center justify-center bg-background">
			<Card className="w-full max-w-md">
				<CardHeader className="text-center">
					<div className="mx-auto mb-4 flex h-12 w-12 items-center justify-center">
						<Loader2 className="h-8 w-8 animate-spin text-primary" />
					</div>
					<CardTitle className="text-2xl">Anmeldung wird verifiziert...</CardTitle>
					<CardDescription>Bitte warte einen Moment.</CardDescription>
				</CardHeader>
			</Card>
		</div>
	)
}

export default function MagicLinkVerifyPage() {
	return (
		<Suspense fallback={<LoadingFallback />}>
			<MagicLinkVerifyContent />
		</Suspense>
	)
}
