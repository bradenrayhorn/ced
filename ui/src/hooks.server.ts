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

const enableCloudflareForwarding = privateEnv.ENABLE_CLOUDFLARE_FORWARDING;

export const handleFetch: HandleFetch = ({ event, request, fetch }) => {
  if (enableCloudflareForwarding) {
    request.headers.set(
      "X-Real-IP",
      event.request.headers.get("CF-Connecting-IP") ?? "",
    );
  }

  return fetch(request);
};
