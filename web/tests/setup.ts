import { test as base } from "@playwright/test";
import { execSync } from "child_process";

export type PrefixFixture = {
  prefix: string;
};

type Fixtures = {
  prefix: PrefixFixture;
};

const randomString = () => Math.random().toString(36);

export const test = base.extend<Fixtures>({
  prefix: async ({}, use) => {
    await use({ prefix: randomString() });
  },
});

export const createGroup = async (
  prefix: string,
  name: string,
  maxAttendees: number,
): Promise<void> => {
  return new Promise((resolve) => {
    execSync(
      `go run ../cmd/ced group create --name="${prefix}${name}" --max-attendees=${maxAttendees}`,
    );
    resolve();
  });
};
