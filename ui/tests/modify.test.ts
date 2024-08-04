import { expect } from "@playwright/test";
import { createGroup, doSearch, getSearchBox, test } from "./setup";

test("can complete an rsvp", async ({ prefix: { prefix }, page }) => {
  await createGroup(prefix, "Fred", 5);

  // search page
  await page.goto("/");
  await expect(page).toHaveTitle("RSVP for An Event");
  await doSearch({ page, prefix, search: "Fred" });

  // is this you page
  await expect(page.getByText("Is this you?")).toBeVisible();
  await page.getByRole("link", { name: "Yes" }).click();

  // modify page
  await page.getByLabel("Two attendees").click();
  await page.getByRole("button", { name: "Confirm" }).click();

  // confirmed page
  await expect(page.getByText("Your RSVP has been received!")).toBeVisible();

  // should have appropriate links
  await expect(
    page.getByRole("link", { name: "View event details" }),
  ).toHaveAttribute("href", "http://localhost:5555");

  await page.getByRole("link", { name: "Edit RSVP" }).click();
  await expect(page).toHaveTitle(`RSVP - ${prefix}Fred`);
});

test("can go back to search from modify page", async ({
  prefix: { prefix },
  page,
}) => {
  await createGroup(prefix, "Fred", 5);

  // search page
  await page.goto("/");
  await expect(page).toHaveTitle("RSVP for An Event");
  await doSearch({ page, prefix, search: "Fred" });

  // is this you page
  await expect(page.getByText("Is this you?")).toBeVisible();
  await page.getByRole("link", { name: "Yes" }).click();

  // modify page
  await page.getByRole("link", { name: "Back to search" }).click();

  // search page
  await expect(page).toHaveURL("/");
  await expect(getSearchBox(page)).toBeVisible();
});

test("can complete rsvp with two results", async ({
  prefix: { prefix },
  page,
}) => {
  await createGroup(prefix, "Fred Jam", 5, "Fred");
  await createGroup(prefix, "Fred Ham", 5, "Fred");

  // search page
  await page.goto("/");
  await expect(page).toHaveTitle("RSVP for An Event");
  await doSearch({ page, prefix, search: "Fred" });

  // select option page
  await expect(page.getByText("Is one of these options you?")).toBeVisible();
  await page.getByLabel(`${prefix}Fred Jam`).click();
  await page.getByRole("button", { name: "Continue" }).click();

  // modify page
  await page.getByLabel("Two attendees").click();
  await page.getByRole("button", { name: "Confirm" }).click();

  // confirmed page
  await expect(page.getByText("Your RSVP has been received!")).toBeVisible();
});

test("shows 404 page when group id not found", async ({ page }) => {
  // unknown page
  await page.goto("/modify/2VPZfkAdudOi1Qm5zaqwnqJHjGy");

  await expect(
    page.getByText("Sorry, the requested page was not found."),
  ).toBeVisible();

  const link = page.getByRole("link", { name: "Back to home" });
  await expect(link).toBeVisible();
  await expect(link).toHaveAttribute("href", "/");
});

test("can handle modify get failure", async ({ page, mockRequest }) => {
  await mockRequest({
    path: "/api/v1/groups/2VPZfkAdudOi1Qm5zaqwnqJHjGy",
    method: "get",
    status: 400,
    body: {},
  });

  await page.goto("/modify/2VPZfkAdudOi1Qm5zaqwnqJHjGy");

  await expect(
    page.getByRole("complementary").getByText("400: Unknown error"),
  ).toBeVisible();
});

test("can handle modify put failure", async ({
  prefix: { prefix },
  page,
  mockRequest,
}) => {
  await createGroup(prefix, "Fred", 5);

  // search page
  await page.goto("/");
  await expect(page).toHaveTitle("RSVP for An Event");
  await doSearch({ page, prefix, search: "Fred" });

  // is this you page
  await expect(page.getByText("Is this you?")).toBeVisible();
  await page.getByRole("link", { name: "Yes" }).click();

  await mockRequest({
    path: "/api/v1/groups/*",
    method: "put",
    status: 400,
    body: {},
  });

  // modify page
  await page.getByLabel("Two attendees").click();
  await page.getByRole("button", { name: "Confirm" }).click();

  await expect(
    page.getByRole("complementary").getByText("400: Unknown error"),
  ).toBeVisible();
});
