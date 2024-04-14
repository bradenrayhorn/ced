import { read } from "$app/server";
import type { PageServerLoad } from "./$types";
import serverLicenses from "./server-licenses.txt";
import uiLicenses from "./ui-licenses.txt";

export const load: PageServerLoad = async () => {
  return {
    notices: {
      ui: await read(uiLicenses).text(),
      server: await read(serverLicenses).text(),
    },
  };
};
