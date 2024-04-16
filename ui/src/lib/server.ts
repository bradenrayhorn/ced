import type { RequestHandler } from "@sveltejs/kit";
import { env as privateEnv } from "$env/dynamic/private";

export const proxyToServer: RequestHandler = ({ request, fetch }) => {
  const url = new URL(request.url);
  const headers = request.headers;

  request = new Request(
    `${privateEnv.UNPROXIED_SERVER_URL ?? "http://localhost:8080"}${url.pathname}${url.search}`,
    request,
  );

  // pass through trusted client ip header if it is set
  const trustedClientIPHeader = privateEnv.TRUSTED_CLIENT_IP_HEADER;
  if (trustedClientIPHeader) {
    request.headers.set(
      trustedClientIPHeader,
      headers.get(trustedClientIPHeader) ?? "",
    );
  }

  return fetch(request);
};
