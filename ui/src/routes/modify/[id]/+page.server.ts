import { redirect } from "@sveltejs/kit";
import type { Actions } from "./$types";
import { doRequest } from "$lib/action-api";

export const actions: Actions = {
  modify: async ({ fetch, request, params }) => {
    await doRequest({
      method: "PUT",
      path: `/api/v1/groups/${params.id}`,
      request,
      fetch,
    });

    redirect(303, `/modify/${params.id}/complete`);
  },
};
