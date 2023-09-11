import { api } from "$lib/api";
import type { PageLoad } from "./$types";

type Output = {
  foundGroups: Group[];
};

type Group = {
  id: string;
  name: string;
  attendees: number;
  max_attendees: number;
  has_responded: boolean;
};

export const load: PageLoad<Output> = async ({ fetch, url }) => {
  const search = url.searchParams.get("search");

  if (search) {
    const res = await fetch(api(`/v1/groups/search?search=${search}`));
    const data = await res.json();
    return { foundGroups: data.data ?? [] };
  }

  return { foundGroups: [] };
};
