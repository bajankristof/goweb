import { requestJSON } from "../client";
import type { User } from "../types";

export type GetCurrentUserOptions = {
  signal?: AbortSignal;
};

export async function getCurrentUser({
  signal,
}: GetCurrentUserOptions = {}): Promise<User> {
  const res = await requestJSON<{ user: User }>("/api/v1/users/me", {
    signal,
  });
  return res.user;
}
