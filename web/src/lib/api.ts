import { env } from "$env/dynamic/public";

export function api(path: string): string {
  return `${env.PUBLIC_BASE_API_URL}/api${path}`;
}
