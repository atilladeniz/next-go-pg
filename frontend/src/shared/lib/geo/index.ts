// Client-safe utilities (no fs dependencies)

export function parseUserAgent(userAgent: string | null): {
	browser: string
	os: string
	device: string
} {
	if (!userAgent) {
		return { browser: "Unbekannt", os: "Unbekannt", device: "Unbekannt" }
	}

	let browser = "Unbekannt"
	let os = "Unbekannt"
	let device = "Desktop"

	if (userAgent.includes("Firefox")) {
		browser = "Firefox"
	} else if (userAgent.includes("Edg/")) {
		browser = "Edge"
	} else if (userAgent.includes("Chrome")) {
		browser = "Chrome"
	} else if (userAgent.includes("Safari")) {
		browser = "Safari"
	} else if (userAgent.includes("Opera") || userAgent.includes("OPR")) {
		browser = "Opera"
	}

	if (userAgent.includes("Windows")) {
		os = "Windows"
	} else if (userAgent.includes("Mac OS")) {
		os = "macOS"
	} else if (userAgent.includes("Linux")) {
		os = "Linux"
	} else if (userAgent.includes("Android")) {
		os = "Android"
		device = "Mobile"
	} else if (userAgent.includes("iPhone") || userAgent.includes("iPad")) {
		os = "iOS"
		device = userAgent.includes("iPad") ? "Tablet" : "Mobile"
	}

	return { browser, os, device }
}

export function formatLocationFromIP(ip: string | null): string {
	if (!ip || ip === "127.0.0.1" || ip === "::1" || ip.startsWith("192.168.")) {
		return "Lokales Netzwerk"
	}
	return ip
}
