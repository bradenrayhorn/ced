import { sveltekit } from "@sveltejs/kit/vite";
import { defineConfig } from "vitest/config";
import Icons from "unplugin-icons/vite";
import license from "rollup-plugin-license";
import path from "path";

export default defineConfig({
  plugins: [
    sveltekit(),
    Icons({ compiler: "svelte" }),
    {
      ...license({
        sourcemap: true,
        thirdParty: {
          multipleVersions: true,
          allow: {
            test: "MIT",
            failOnViolation: true,
            failOnUnlicensed: true,
          },
          output: {
            file: path.join(
              __dirname,
              "src",
              "routes",
              "about",
              "ui-licenses.txt",
            ),
          },
        },
      }),
      enforce: "post",
      apply: (config) => {
        return (
          // frontend build has assets limit
          !!process.env.GENERATE_LICENSES && !!config.build?.assetsInlineLimit
        );
      },
    },
    {
      ...license({
        sourcemap: true,
        thirdParty: {
          multipleVersions: true,
          allow: {
            test: "MIT",
            failOnViolation: true,
            failOnUnlicensed: true,
          },
          output: {
            file: path.join(
              __dirname,
              "src",
              "routes",
              "about",
              "ui-server-licenses.txt",
            ),
          },
        },
      }),
      enforce: "post",
      apply: (config) => {
        return (
          // server-side build has no assets limit
          !!process.env.GENERATE_LICENSES && !config.build?.assetsInlineLimit
        );
      },
    },
  ],
  test: {
    include: ["src/**/*.{test,spec}.{js,ts}"],
  },
});
