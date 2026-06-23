import { queryOptions, useQuery } from "@tanstack/react-query";

import { getCurrentUser } from "../api";

export const currentUserQuery = queryOptions({
  queryKey: ["auth", "user"],
  queryFn: ({ signal }) => getCurrentUser({ signal }),
  staleTime: 5 * 60 * 1000,
});

export default function useCurrentUser() {
  return useQuery(currentUserQuery);
}
