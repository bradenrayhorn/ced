import { expect } from "@playwright/test";
import { test } from "./setup";

test("has about page", async ({ page }) => {
  // search page
  await page.goto("/about");

  await expect(
    page.getByText("ced is a self-hosted RSVP service."),
  ).toBeVisible();
});
