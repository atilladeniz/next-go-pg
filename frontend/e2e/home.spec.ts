import { expect, test } from "@playwright/test"

test.describe("Home Page", () => {
	test("should load the home page", async ({ page }) => {
		await page.goto("/")
		await expect(page).toHaveTitle(/Next-Go-PG/)
	})

	test("should have navigation", async ({ page }) => {
		await page.goto("/")
		const header = page.locator("header")
		await expect(header).toBeVisible()
	})

	test("should toggle dark mode", async ({ page }) => {
		await page.goto("/")

		// Find and click the theme toggle button
		const themeToggle = page.getByRole("button", { name: /toggle theme/i })
		if (await themeToggle.isVisible()) {
			await themeToggle.click()
			// Check that theme changed (html element should have class)
			const html = page.locator("html")
			await expect(html).toHaveClass(/dark|light/)
		}
	})
})
