import { error, type NumericRange } from "@sveltejs/kit";

const defaultError = "Unknown error";

export const getError = async (res: Response) => {
  const errorJson = await res
    .json()
    .catch(async () => await res.text().catch(() => defaultError));
  const msg = errorJson?.error ?? defaultError;

  if (res.status >= 400 && res.status <= 599) {
    error(res.status as NumericRange<400, 599>, msg);
  }

  error(500, msg);
};
