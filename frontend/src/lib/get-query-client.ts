import { defaultShouldDehydrateQuery, isServer, QueryClient } from "@tanstack/react-query"

function makeQueryClient() {
	return new QueryClient({
		defaultOptions: {
			queries: {
				// With SSR, we usually want to set some default staleTime
				// above 0 to avoid refetching immediately on the client
				staleTime: 60 * 1000,
			},
			dehydrate: {
				// Include pending queries in dehydration for streaming
				shouldDehydrateQuery: (query) =>
					defaultShouldDehydrateQuery(query) || query.state.status === "pending",
			},
		},
	})
}

let browserQueryClient: QueryClient | undefined

export function getQueryClient() {
	if (isServer) {
		// Server: always make a new query client
		return makeQueryClient()
	}
	// Browser: make a new query client if we don't already have one
	// This is very important, so we don't re-make a new client if React
	// suspends during the initial render
	if (!browserQueryClient) browserQueryClient = makeQueryClient()
	return browserQueryClient
}
