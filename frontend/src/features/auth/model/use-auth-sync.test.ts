import { describe, expect, it, vi } from "vitest"

// Mock broadcast-channel before importing the module
vi.mock("broadcast-channel", () => {
	return {
		BroadcastChannel: class MockBroadcastChannel {
			postMessage = vi.fn().mockResolvedValue(undefined)
			close = vi.fn()
			onmessage: ((msg: unknown) => void) | null = null
		},
	}
})

// Import after mock
const { broadcastSignIn, broadcastSignOut } = await import("./use-auth-sync")

describe("Auth Sync", () => {
	describe("broadcastSignOut", () => {
		it("should not throw when called in browser environment", async () => {
			await expect(broadcastSignOut()).resolves.toBeUndefined()
		})

		it("should return early when window is undefined", async () => {
			const originalWindow = global.window
			// @ts-expect-error - Testing SSR environment
			delete global.window

			const result = await broadcastSignOut()
			expect(result).toBeUndefined()

			global.window = originalWindow
		})
	})

	describe("broadcastSignIn", () => {
		it("should not throw when called in browser environment", async () => {
			await expect(broadcastSignIn()).resolves.toBeUndefined()
		})

		it("should return early when window is undefined", async () => {
			const originalWindow = global.window
			// @ts-expect-error - Testing SSR environment
			delete global.window

			const result = await broadcastSignIn()
			expect(result).toBeUndefined()

			global.window = originalWindow
		})
	})
})
