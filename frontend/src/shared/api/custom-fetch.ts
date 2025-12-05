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

	// Return response with status for Orval's discriminated unions
	return {
		data,
		status: response.status,
		headers: response.headers,
	} as T
}
