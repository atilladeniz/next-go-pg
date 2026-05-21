export const customFetch = async <T>(url: string, options?: RequestInit): Promise<T> => {
	const response = await fetch(url, {
		...options,
		credentials: "include",
		headers: {
			"Content-Type": "application/json",
			...options?.headers,
		},
	})

	const contentType = response.headers.get("content-type")
	let data: unknown

	if (contentType?.includes("application/json")) {
		data = await response.json()
	} else {
		const text = await response.text()
		data = { message: text }
	}

	// Return response with status for Orval's discriminated unions. Headers
	// flattened to a plain Record so React Query's dehydrate → server-to-
	// client transfer can JSON-serialize the cached value (a live Headers
	// object would silently drop on the wire and break SSR hydration).
	const flatHeaders: Record<string, string> = {}
	response.headers.forEach((value, key) => {
		flatHeaders[key] = value
	})
	return {
		data,
		status: response.status,
		headers: flatHeaders,
	} as T
}
