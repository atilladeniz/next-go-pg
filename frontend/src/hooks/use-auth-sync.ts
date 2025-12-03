"use client"

import { BroadcastChannel } from "broadcast-channel"
import { useRouter } from "next/navigation"
import { useEffect, useRef } from "react"

type AuthMessage = {
	type: "SIGN_OUT" | "SIGN_IN"
	timestamp: number
}

const CHANNEL_NAME = "auth-sync"

// Singleton channel instance für broadcast functions
let broadcastChannelInstance: BroadcastChannel<AuthMessage> | null = null

function getBroadcastChannel(): BroadcastChannel<AuthMessage> {
	if (!broadcastChannelInstance) {
		broadcastChannelInstance = new BroadcastChannel<AuthMessage>(CHANNEL_NAME, {
			// localStorage fallback für bessere Browser-Kompatibilität (z.B. Safari Private Mode)
			type: "localstorage",
		})
	}
	return broadcastChannelInstance
}

/**
 * Hook für Cross-Tab Auth Synchronisation.
 *
 * Nutzt die broadcast-channel Library für bessere Browser-Kompatibilität
 * (inkl. Safari Private Mode, IE11, etc.)
 *
 * @example
 * // In einer Client Component:
 * useAuthSync()
 *
 * // Beim Logout broadcasten:
 * import { broadcastSignOut } from "@/hooks/use-auth-sync"
 * await signOut()
 * broadcastSignOut()
 */
export function useAuthSync() {
	const router = useRouter()
	const channelRef = useRef<BroadcastChannel<AuthMessage> | null>(null)

	useEffect(() => {
		// Channel nur im Browser erstellen
		if (typeof window === "undefined") return

		const channel = new BroadcastChannel<AuthMessage>(CHANNEL_NAME, {
			type: "localstorage",
		})
		channelRef.current = channel

		channel.onmessage = (msg: AuthMessage) => {
			switch (msg.type) {
				case "SIGN_OUT":
					// Andere Tabs haben sich ausgeloggt → Redirect zu Login
					router.push("/login")
					router.refresh()
					break
				case "SIGN_IN":
					// Andere Tabs haben sich eingeloggt → Seite refreshen
					router.refresh()
					break
			}
		}

		return () => {
			channel.close()
			channelRef.current = null
		}
	}, [router])
}

/**
 * Broadcast Logout zu allen anderen Tabs.
 * Aufrufen NACH signOut().
 */
export async function broadcastSignOut() {
	if (typeof window === "undefined") return

	const channel = getBroadcastChannel()
	await channel.postMessage({
		type: "SIGN_OUT",
		timestamp: Date.now(),
	})
}

/**
 * Broadcast Login zu allen anderen Tabs.
 * Aufrufen NACH signIn().
 */
export async function broadcastSignIn() {
	if (typeof window === "undefined") return

	const channel = getBroadcastChannel()
	await channel.postMessage({
		type: "SIGN_IN",
		timestamp: Date.now(),
	})
}
