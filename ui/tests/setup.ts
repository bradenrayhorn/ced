import { test as base, type Locator, type Page } from "@playwright/test";
import { execSync } from "child_process";
import getPort from "get-port";
import { http } from "msw";
import type { SetupServer } from "msw/node";
import { setupServer } from "msw/node";

type MockRequest = ({
  path,
  method,
  status,
  body,
}: {
  path: string;
  method: "get" | "post" | "put";
  status: number;
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  body: Record<string, any>;
}) => Promise<void>;

type Fixtures = {
  prefix: { prefix: string };
  http: typeof http;
  mockRequest: MockRequest;
};

type WorkerFixtures = {
  port: string;
  requestInterceptor: SetupServer;
};

const randomString = () => Math.random().toString(36);

export const test = base.extend<Fixtures, WorkerFixtures>({
  // eslint-disable-next-line no-empty-pattern
  prefix: async ({}, use) => {
    await use({ prefix: randomString() });
  },
  port: [
    // eslint-disable-next-line no-empty-pattern
    async ({}, use) => {
      const port = `${await getPort()}`;
      process.env.PORT = port;
      process.env.ORIGIN = `http://127.0.0.1:${port}`;
      process.env.PUBLIC_EVENT_TITLE = "An Event";
      process.env.PUBLIC_EVENT_URL = "http://localhost:5555";

      // run sveltekit server
      await import("../build/index.js");
      await use(port);
    },
    { scope: "worker", auto: true },
  ],
  requestInterceptor: [
    // eslint-disable-next-line no-empty-pattern
    async ({}, use) => {
      await use(
        (() => {
          const requestInterceptor = setupServer();

          requestInterceptor.listen({
            onUnhandledRequest: "bypass",
          });

          return requestInterceptor;
        })(),
      );
    },
    {
      scope: "worker",
    },
  ],
  http,
  baseURL: async ({ port }, use) => {
    await use(`http://127.0.0.1:${port}/`);
  },
  mockRequest: async ({ requestInterceptor, http }, use) => {
    await use(async ({ path, method, status, body }): Promise<void> => {
      // apply msw mock (for node server)
      requestInterceptor.use(
        http[method](`*${path}`, () => {
          return new Response(JSON.stringify(body), {
            status: status,
            headers: {
              "Content-Type": "application/json",
            },
          });
        }),
      );
    });
  },
});

export const createGroup = async (
  prefix: string,
  name: string,
  maxAttendees: number,
  searchHints: string = "",
): Promise<void> => {
  return new Promise((resolve) => {
    execSync(
      `(cd ../server && go run ./cmd/ced group create --name="${prefix}${name}" --max-attendees=${maxAttendees} --search-hints="${searchHints}")`,
    );
    resolve();
  });
};

export function getSearchBox(page: Page): Locator {
  return page.getByRole("textbox", { name: "Search" });
}

export async function doSearch({
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
