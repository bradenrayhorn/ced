import type { PlaywrightTestConfig } from "@playwright/test";

const config: PlaywrightTestConfig = {
  webServer: [
    {
      command: "npm run build && npm run preview",
      port: 4173,
      reuseExistingServer: !process.env.CI,
      env: {
        PUBLIC_EVENT_TITLE: "An Event",
        PUBLIC_EVENT_URL: "http://localhost:5555",
        PUBLIC_BASE_API_URL: "",
      },
    },
    {
      command: "rm -f ced-e2e.db* && (cd ../server && go run ./cmd/cedd)",
      url: "http://localhost:8080/api/v1/live",
      env: {
        HTTP_PORT: "8080",
      },
      reuseExistingServer: !process.env.CI,
    },
  ],
  use: {
    baseURL: "http://localhost:4173",
  },
  testDir: "tests",
  testMatch: /(.+\.)?(test|spec)\.[jt]s/,
};

export default config;
