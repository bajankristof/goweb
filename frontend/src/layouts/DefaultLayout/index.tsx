import { Outlet } from "react-router";

import Header from "#components/Header";
import { requireSignIn } from "#hooks/guards";

export const defaultLayoutLoader = requireSignIn;

export default function DefaultLayout() {
  return (
    <>
      <Header />
      <main className="container-fluid">
        <Outlet />
      </main>
    </>
  );
}
