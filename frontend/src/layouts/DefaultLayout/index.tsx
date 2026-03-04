import { QueryClient } from "@tanstack/react-query";
import { Outlet, redirect, useLoaderData } from "react-router";

import type { User } from "../../api";
import { authWellKnownQuery, currentUserQuery } from "../../api/queries";
import Header from "./components/Header";

type LoaderData = {
  user?: User;
};

export function defaultLayoutLoader(queryClient: QueryClient) {
  return async (): Promise<LoaderData> => {
    const [wellKnown, user] = await Promise.all([
      queryClient.fetchQuery(authWellKnownQuery),
      queryClient.fetchQuery(currentUserQuery),
    ]);

    if (user) {
      return { user };
    }

    const { providers } = wellKnown;
    if (providers.length > 1) {
      throw redirect("/signin");
    }

    window.location.href = providers.length === 1 ? "/auth/signin" : "/auth/authless";
    return {};
  };
}

export default function DefaultLayout() {
  const { user } = useLoaderData<LoaderData>();
  if (!user) {
    return <main aria-busy="true"></main>;
  }

  return (
    <>
      <Header />
      <main className="container">
        <Outlet />
      </main>
    </>
  );
}
