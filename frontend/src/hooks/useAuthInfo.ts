import * as API from "../api";
import type { AuthWellKnown, User } from "../api";

export type AuthInfo = {
  wellKnown: AuthWellKnown;
  user?: User;
};

export default async function useAuthInfo(): Promise<AuthInfo> {
  const [wellKnown, user] = await Promise.all([
    API.get<AuthWellKnown>("/auth/well-known"),
    API.request("/auth/refresh")
      .then(() => API.get<User>("/api/v1/users/u"))
      .catch(() => undefined),
  ]);

  return {
    wellKnown,
    user,
  };
}
