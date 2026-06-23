import type { QueryClient } from "@tanstack/react-query";
import { redirect } from "react-router";

import { APIError } from "../api";
import { currentUserQuery } from "./useCurrentUser";

export function requireSignIn(queryClient: QueryClient) {
  return async () => {
    try {
      await queryClient.ensureQueryData(currentUserQuery);
    } catch (err) {
      if (err instanceof APIError && err.status === 401) {
        throw redirect("/signin");
      }
      throw err;
    }
    return null;
  };
}

export function rejectSignedIn(queryClient: QueryClient) {
  return async () => {
    try {
      await queryClient.ensureQueryData(currentUserQuery);
    } catch (err) {
      if (err instanceof APIError && err.status === 401) {
        return null;
      }
      throw err;
    }
    throw redirect("/");
  };
}
