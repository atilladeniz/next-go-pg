import { describe, expect, it } from "vitest"
import { cn } from "./utils"

describe("cn (className utility)", () => {
	it("should merge class names", () => {
		expect(cn("foo", "bar")).toBe("foo bar")
	})

	it("should handle conditional classes", () => {
		expect(cn("foo", false && "bar", "baz")).toBe("foo baz")
	})

	it("should handle undefined values", () => {
		expect(cn("foo", undefined, "bar")).toBe("foo bar")
	})

	it("should handle null values", () => {
		expect(cn("foo", null, "bar")).toBe("foo bar")
	})

	it("should merge tailwind classes correctly", () => {
		expect(cn("px-2 py-1", "px-4")).toBe("py-1 px-4")
	})

	it("should handle conflicting tailwind classes", () => {
		expect(cn("text-red-500", "text-blue-500")).toBe("text-blue-500")
	})

	it("should handle array of classes", () => {
		expect(cn(["foo", "bar"])).toBe("foo bar")
	})

	it("should handle object syntax", () => {
		expect(cn({ foo: true, bar: false, baz: true })).toBe("foo baz")
	})

	it("should handle empty input", () => {
		expect(cn()).toBe("")
	})

	it("should handle complex combinations", () => {
		expect(
			cn(
				"base-class",
				["array-class"],
				{ "object-class": true },
				undefined,
				null,
				false && "conditional-false",
				true && "conditional-true",
			),
		).toBe("base-class array-class object-class conditional-true")
	})
})
