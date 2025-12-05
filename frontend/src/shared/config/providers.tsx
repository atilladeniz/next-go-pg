"use client"

import { QueryClientProvider } from "@tanstack/react-query"
import type { ReactNode } from "react"
import { getQueryClient } from "../lib/query-client"
import { ThemeProvider } from "./theme-provider"

export function Providers({ children }: { children: ReactNode }) {
	// NOTE: Avoid useState when initializing the query client if you don't
	// have a suspense boundary between this and the code that may
	// suspend because React will throw away the client on the initial
	// render if it suspends and there is no boundary
	const queryClient = getQueryClient()

	return (
		<ThemeProvider attribute="class" defaultTheme="system" enableSystem disableTransitionOnChange>
			<QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
		</ThemeProvider>
	)
}
