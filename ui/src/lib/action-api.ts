import { api } from "./api";
import { getError } from "./fetch-error";

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

  const res = await internalFetch(api(path), {
    method,
    body: JSON.stringify(obj),
  });

  if (!res.ok) {
    await getError(res);
  }

  return res;
};
