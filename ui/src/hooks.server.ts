import type { Handle, HandleFetch } from "@sveltejs/kit";
import { env } from "$env/dynamic/public";
import { env as privateEnv } from "$env/dynamic/private";

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

const forwarded = [privateEnv.CLIENT_IP_HEADER ?? ""];

export const handleFetch: HandleFetch = ({ event, request, fetch }) => {
  for (const header of forwarded.filter((h) => h.trim().length > 0)) {
    const value = event.request.headers.get(header);
    if (value !== null && !request.headers.has(header)) {
      request.headers.set(header, value);
    }
  }

  return fetch(request);
};
