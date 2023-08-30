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
  const res = await fetch(`/api/v1/groups/${params.id}`);
  return (await res.json()) as OutputProps;
};
