import { redirect } from "@sveltejs/kit";
import type { Actions } from "./$types";
import { doRequest } from "$lib/action-api";

export const actions: Actions = {
  modify: async ({ fetch, request, params }) => {
    await doRequest({
      method: "PUT",
      path: `/v1/groups/${params.id}`,
      request,
      fetch,
    });

    throw redirect(303, `/modify/${params.id}/complete`);
  },
};
