import { expect } from "@playwright/test";
import { test } from "./setup";

test("has footer with event title and link", async ({ page }) => {
  // search page
  await page.goto("/");

  const footer = page.getByRole("contentinfo");
  const link = footer.getByRole("link", { name: "An Event" });

  await expect(link).toBeVisible();
  await expect(link).toHaveAttribute("href", "http://localhost:5555");
});

test("shows 404 page", async ({ page }) => {
  // unknown page
  await page.goto("/unknown");

  await expect(
    page.getByText("Sorry, the requested page was not found."),
  ).toBeVisible();

  const link = page.getByRole("link", { name: "Back to home" });
  await expect(link).toBeVisible();
  await expect(link).toHaveAttribute("href", "/");
});
