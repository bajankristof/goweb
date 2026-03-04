import type { QueryClient } from "@tanstack/react-query";
import { createBrowserRouter } from "react-router";

import DefaultLayout, { defaultLayoutLoader } from "./layouts/DefaultLayout";
import SignIn, { signInLoader } from "./pages/SignIn";
import Home from "./pages/Home";

export function createRouter(queryClient: QueryClient) {
  return createBrowserRouter([
    {
      path: "signin",
      loader: signInLoader(queryClient),
      Component: SignIn,
    },
    {
      path: "/",
      loader: defaultLayoutLoader(queryClient),
      Component: DefaultLayout,
      children: [{ index: true, Component: Home }],
    },
  ]);
}
