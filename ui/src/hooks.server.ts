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

export const handleFetch: HandleFetch = ({ event, request, fetch }) => {
  const url = new URL(request.url);

  if (privateEnv.UNPROXIED_BASE_API_URL) {
    request = new Request(
      `${privateEnv.UNPROXIED_BASE_API_URL ?? ""}${url.pathname}${url.search}`,
      request,
    );
  }

  // pass through trusted client ip header if it is set
  const trustedClientIPHeader = privateEnv.TRUSTED_CLIENT_IP_HEADER;
  if (trustedClientIPHeader) {
    request.headers.set(
      trustedClientIPHeader,
      event.request.headers.get(trustedClientIPHeader) ?? "",
    );
  }

  return fetch(request);
};
