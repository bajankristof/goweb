import { queryOptions, useQuery } from "@tanstack/react-query";

import { getInfo } from "../api";

export const wellKnownInfoQuery = queryOptions({
  queryKey: ["wellKnown", "info"],
  queryFn: ({ signal }) => getInfo({ signal }),
  staleTime: Infinity,
});

export default function useWellKnownInfo() {
  return useQuery(wellKnownInfoQuery);
}
