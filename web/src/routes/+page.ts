import type { PageLoad } from "./$types";

type Output = {
  foundGroup: Group | null;
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

  if (!!search) {
    const res = await fetch(`/api/v1/groups/search?search=${search}`);
    const data = await res.json();
    const result = data.data?.[0];
    if (result) {
      return { foundGroup: result };
    }
  }

  return { foundGroup: null };
};
