import { redirect } from "@sveltejs/kit";
import type { Actions } from "./$types";

export const actions: Actions = {
  toModify: async ({ request }) => {
    const data = await request.formData();
    const groupID = data.get("group");

    if (groupID && groupID !== "0") {
      throw redirect(303, `/modify/${groupID}`);
    } else {
      throw redirect(303, `/`);
    }
  },
};
