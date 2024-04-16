import type { RequestHandler } from "@sveltejs/kit";
import { env as privateEnv } from "$env/dynamic/private";

export const proxyToServer: RequestHandler = ({ request, fetch, ...event }) => {
  const url = new URL(request.url);

  request = new Request(
    `${privateEnv.UNPROXIED_SERVER_URL ?? "http://localhost:8080"}${url.pathname}${url.search}`,
    request,
  );

  // pass through trusted client ip header if it is set
  request.headers.set("ced-connecting-ip", event.getClientAddress());

  return fetch(request);
};
