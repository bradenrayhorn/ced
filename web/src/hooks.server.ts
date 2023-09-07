import type { Handle } from "@sveltejs/kit";
import { env } from "$env/dynamic/public";

export const handle: Handle = async ({ event, resolve }) => {
  let theme = env.PUBLIC_THEME ?? "";
  const validThemes = ["hamlindigo", "cardstock"];
  if (!validThemes.includes(theme)) {
    theme = "hamlindigo";
  }

  return await resolve(event, {
    transformPageChunk: ({ html }) =>
      html.replace('data-theme=""', `data-theme="${theme}"`),
  });
};
