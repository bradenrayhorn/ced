import { api } from "$lib/api";
import { getError } from "$lib/fetch-error";
import type { PageLoad } from "./$types";

type OutputProps = {
  data: Group;
};

type Group = {
  id: string;
  name: string;
  attendees: number;
  max_attendees: number;
  has_responded: boolean;
};

export const load: PageLoad<OutputProps> = async ({ fetch, params }) => {
  const res = await fetch(api(`/v1/groups/${params.id}`));
  if (!res.ok) {
    await getError(res);
  }
  return (await res.json()) as OutputProps;
};
