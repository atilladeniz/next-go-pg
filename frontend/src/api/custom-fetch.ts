export const customFetch = async <T>(url: string, options?: RequestInit): Promise<T> => {
	const response = await fetch(url, {
		...options,
		credentials: "include",
		headers: {
			"Content-Type": "application/json",
			...options?.headers,
		},
	})

	if (!response.ok) {
		const error = await response.json().catch(() => ({ error: "Unknown error" }))
		throw new Error(error.error || `HTTP ${response.status}`)
	}

	return response.json()
}
