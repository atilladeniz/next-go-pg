import { expect, test } from "@playwright/test"

test.describe("Authentication", () => {
	test("should show login page", async ({ page }) => {
		await page.goto("/login")
		await expect(page.getByRole("heading", { name: /anmelden/i })).toBeVisible()
	})

	test("should have email input on login page", async ({ page }) => {
		await page.goto("/login")
		const emailInput = page.getByPlaceholder(/e-mail/i)
		await expect(emailInput).toBeVisible()
	})

	test("should show error for invalid email", async ({ page }) => {
		await page.goto("/login")

		const emailInput = page.getByPlaceholder(/e-mail/i)
		await emailInput.fill("invalid-email")

		const submitButton = page.getByRole("button", { name: /link senden/i })
		if (await submitButton.isVisible()) {
			await submitButton.click()
			// Should show validation error
			await expect(page.getByText(/gültige/i)).toBeVisible({ timeout: 5000 })
		}
	})

	test("should redirect unauthenticated users from dashboard", async ({
		page,
	}) => {
		await page.goto("/dashboard")
		// Should redirect to login
		await expect(page).toHaveURL(/login/)
	})

	test("should redirect unauthenticated users from settings", async ({
		page,
	}) => {
		await page.goto("/settings")
		// Should redirect to login
		await expect(page).toHaveURL(/login/)
	})
})
