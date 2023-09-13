import { error } from "@sveltejs/kit";

const defaultError = "Unknown error";

export const getError = async (res: Response) => {
  const errorJson = await res
    .json()
    .catch(async () => await res.text().catch(() => defaultError));
  const msg = errorJson?.error ?? defaultError;
  return error(res.status, msg);
};
