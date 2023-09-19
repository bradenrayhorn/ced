import { expect } from "@playwright/test";
import { createGroup, doSearch, getSearchBox, test } from "./setup";

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

test("can handle search failure", async ({
  prefix: { prefix },
  page,
  mockRequest,
}) => {
  await mockRequest({
    path: "/api/v1/groups/search",
    method: "get",
    status: 400,
    body: {},
  });

  await page.goto("/");

  await doSearch({ page, prefix, search: "Fred" });

  await expect(
    page.getByRole("complementary").getByText("400: Unknown error"),
  ).toBeVisible();
});
