import { expect, type Locator, type Page } from "@playwright/test";
import { createGroup, test } from "./setup";

function getSearchBox(page: Page): Locator {
  return page.getByRole("textbox", { name: "Search" });
}

async function doSearch({
  page,
  prefix,
  search,
}: {
  page: Page;
  prefix: string;
  search: string;
}) {
  await getSearchBox(page).fill(prefix + search);
  await page.getByRole("button", { name: "Search" }).click();
}

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
  await page.getByLabel("Two guests").click();
  await page.getByRole("button", { name: "Confirm" }).click();

  // confirmed page
  await expect(page.getByText("Your RSVP has been received!")).toBeVisible();

  // should have appropriate links
  await expect(
    page.getByRole("link", { name: "View event details" }),
  ).toHaveAttribute("href", "http://localhost:5555");

  await page.getByRole("link", { name: "Edit RSVP" }).click();
  await expect(page).toHaveTitle("Modify RSVP");
});

test("can have search return no results", async ({
  prefix: { prefix },
  page,
}) => {
  await page.goto("/");

  await doSearch({ page, prefix, search: "Fred" });

  await expect(
    page
      .getByRole("complementary")
      .getByText("Please refine your search and try again."),
  ).toBeVisible();
});

test("can answer no to is this you", async ({ prefix: { prefix }, page }) => {
  await createGroup(prefix, "Fred", 5);

  // search page
  await page.goto("/");
  await expect(page).toHaveTitle("RSVP for An Event");
  await doSearch({ page, prefix, search: "Fred" });

  // is this you page
  await expect(page.getByText("Is this you?")).toBeVisible();
  await page.getByRole("link", { name: "No" }).click();

  // search page
  await expect(getSearchBox(page)).toBeVisible();
  await expect(getSearchBox(page)).toHaveValue("");
});

test("has footer with event title and link", async ({ page }) => {
  // search page
  await page.goto("/");

  const footer = page.getByRole("contentinfo");
  const link = footer.getByRole("link", { name: "An Event" });

  await expect(link).toBeVisible();
  await expect(link).toHaveAttribute("href", "http://localhost:5555");
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
  await page.getByLabel("Fred Jam").click();
  await page.getByRole("button", { name: "Continue" }).click();

  // modify page
  await page.getByLabel("Two guests").click();
  await page.getByRole("button", { name: "Confirm" }).click();

  // confirmed page
  await expect(page.getByText("Your RSVP has been received!")).toBeVisible();
});

test("can go back to search from rsvp with two results", async ({
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
  await page.getByLabel("None of the above").click();
  await page.getByRole("button", { name: "Continue" }).click();

  // search page
  await expect(getSearchBox(page)).toBeVisible();
  await expect(getSearchBox(page)).toHaveValue("");
});
