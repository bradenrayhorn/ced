import { read } from "$app/server";
import type { PageServerLoad } from "./$types";
import serverLicenses from "./server-licenses.txt";
import uiLicenses from "./ui-licenses.txt";
import uiServerLicenses from "./ui-server-licenses.txt";

export const load: PageServerLoad = async () => {
  return {
    notices: {
      ui: await read(uiLicenses).text(),
      uiServer: await read(uiServerLicenses).text(),
      server: await read(serverLicenses).text(),
    },
  };
};
