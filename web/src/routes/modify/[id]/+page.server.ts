import { redirect } from "@sveltejs/kit";
import type { Actions } from "./$types";

export const actions: Actions = {
  modify: async ({ fetch, request, params }) => {
    const data = await request.formData();
    const attendees = data.get("attendees")?.toString();
    if (!attendees || Number.isNaN(+attendees)) {
      throw Error("Invalid attendees.");
    }

    const res = await fetch(`/api/v1/groups/${params.id}`, {
      method: "PUT",
      body: JSON.stringify({ attendees: +attendees }),
    });

    if (!res.ok) {
      throw Error(await res.text());
    }

    throw redirect(303, `/modify/${params.id}/complete`);
  },
};
