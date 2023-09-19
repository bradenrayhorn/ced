import type { PlaywrightTestConfig } from "@playwright/test";
import { devices } from "@playwright/test";

const config: PlaywrightTestConfig = {
  webServer: [
    {
      command:
        "rm -f ced-e2e.db* && (cd ../server && ORIGIN=* go run ./cmd/cedd)",
      url: "http://localhost:8080/api/v1/live",
      env: {
        HTTP_PORT: "8080",
      },
      reuseExistingServer: !process.env.CI,
    },
  ],
  testDir: "tests",
  testMatch: /(.+\.)?(test|spec)\.[jt]s/,
  projects: [
    {
      name: "chromium",
      use: {
        ...devices["Desktop Chrome"],
      },
    },
    {
      name: "chromium",
      use: {
        ...devices["Desktop Chrome"],
        javaScriptEnabled: false,
      },
    },
  ],
};

export default config;
