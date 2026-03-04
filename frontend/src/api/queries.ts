import { queryOptions } from "@tanstack/react-query";

import { get } from "./client";
import type { AuthWellKnown, User } from "./index";

export const authWellKnownQuery = queryOptions({
  queryKey: ["auth", "well-known"],
  queryFn: () => get<AuthWellKnown>("/auth/well-known"),
  staleTime: Infinity,
});

export const currentUserQuery = queryOptions({
  queryKey: ["auth", "user"],
  queryFn: () => get<User>("/api/v1/users/u").catch(() => undefined),
  staleTime: 5 * 60 * 1000,
});
