import { env } from "$env/dynamic/private";

const trustedClientIPHeader = env.TRUSTED_CLIENT_IP_HEADER;

export function unproxiedApi(path: string): string {
  return `${env.UNPROXIED_BASE_API_URL}/api${path}`;
}

export const doRequest = async ({
  method,
  path,
  fetch: internalFetch,
  request,
}: {
  method: string;
  path: string;
  request: Request;
  fetch: typeof fetch;
}): Promise<Response> => {
  const data = await request.formData();
  const obj: { [key: string]: string } = {};
  data.forEach((value, key) => {
    obj[key] = value.toString();
  });

  const headers: HeadersInit = {};

  if (trustedClientIPHeader) {
    headers[trustedClientIPHeader] =
      request.headers.get(trustedClientIPHeader) ?? "";
  }

  const res = await internalFetch(unproxiedApi(path), {
    method,
    body: JSON.stringify(obj),
    headers,
  });

  if (!res.ok) {
    throw Error(await res.text());
  }

  return res;
};
