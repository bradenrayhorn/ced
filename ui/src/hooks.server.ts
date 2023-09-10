import type { Handle, HandleFetch } from "@sveltejs/kit";
import { env } from "$env/dynamic/public";
import { env as privateEnv } from "$env/dynamic/private";

const trustedClientIPHeader = privateEnv.TRUSTED_CLIENT_IP_HEADER;

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
  if (env.PUBLIC_BASE_API_URL === "") {
    // if base api url is empty, put unproxied base api url at start
    request = new Request(
      `${privateEnv.UNPROXIED_BASE_API_URL ?? ""}${request.url}`,
      request,
    );
  } else if (
    !!env.PUBLIC_BASE_API_URL &&
    request.url.startsWith(env.PUBLIC_BASE_API_URL)
  ) {
    // if base api url is set, replace with unproxied base api url
    request = new Request(
      request.url.replace(
        env.PUBLIC_BASE_API_URL,
        privateEnv.UNPROXIED_BASE_API_URL ?? "",
      ),
      request,
    );
  }

  // pass through trusted client ip header if it is set
  if (trustedClientIPHeader) {
    request.headers.set(
      trustedClientIPHeader,
      event.request.headers.get(trustedClientIPHeader) ?? "",
    );
  }

  return fetch(request);
};
