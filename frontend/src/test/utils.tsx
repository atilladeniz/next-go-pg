import { QueryClient, QueryClientProvider } from "@tanstack/react-query"
import { type RenderOptions, render } from "@testing-library/react"
import type { ReactElement, ReactNode } from "react"

// Create a new QueryClient for each test
function createTestQueryClient() {
	return new QueryClient({
		defaultOptions: {
			queries: {
				retry: false,
				gcTime: 0,
				staleTime: 0,
			},
			mutations: {
				retry: false,
			},
		},
	})
}

interface WrapperProps {
	children: ReactNode
}

function createWrapper() {
	const queryClient = createTestQueryClient()

	return function Wrapper({ children }: WrapperProps) {
		return <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
	}
}

// Custom render that includes providers
function customRender(ui: ReactElement, options?: Omit<RenderOptions, "wrapper">) {
	return render(ui, { wrapper: createWrapper(), ...options })
}

// Re-export everything
export * from "@testing-library/react"
export { customRender as render }
export { createTestQueryClient }
