import { request } from "../client";

export type SignOutOptions = {
  signal?: AbortSignal;
};

export async function signOut({ signal }: SignOutOptions = {}): Promise<void> {
  await request("/auth/signout", {
    method: "POST",
    signal,
    retryOnUnauthorized: false,
  });
}
