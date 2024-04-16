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
      apply: (_, configEnv) => {
        return !configEnv.isSsrBuild;
      },
    },
  ],
  test: {
    include: ["src/**/*.{test,spec}.{js,ts}"],
  },
});
