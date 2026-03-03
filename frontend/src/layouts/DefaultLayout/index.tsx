import { Outlet, redirect, useLoaderData } from "react-router";

import type { User } from "../../api";
import useAuthInfo from "../../hooks/useAuthInfo";
import Header from "./components/Header";

type LoaderData = {
  user?: User;
};

export async function defaultLayoutLoader(): Promise<LoaderData> {
  const { wellKnown, user } = await useAuthInfo();
  const { providers } = wellKnown;

  if (user) {
    return { user };
  } else if (providers.length > 1) {
    throw redirect("/signin");
  }

  window.location.href = "/auth/signin";
  return {};
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
