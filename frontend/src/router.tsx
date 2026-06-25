import type { QueryClient } from "@tanstack/react-query";
import { createBrowserRouter } from "react-router";

import ErrorBoundary from "#components/ErrorBoundary";
import InsetLoader from "#components/InsetLoader";
import Home from "#pages/Home";
import NotFound from "#pages/NotFound";
import SignIn, { signInLoader } from "#pages/SignIn";
import DefaultLayout, { defaultLayoutLoader } from "./layouts/DefaultLayout";

export function createRouter(queryClient: QueryClient) {
  return createBrowserRouter([
    {
      path: "signin",
      loader: signInLoader(queryClient),
      Component: SignIn,
      HydrateFallback: InsetLoader,
      ErrorBoundary,
    },
    {
      path: "/",
      loader: defaultLayoutLoader(queryClient),
      Component: DefaultLayout,
      children: [{ index: true, Component: Home }],
      HydrateFallback: InsetLoader,
      ErrorBoundary,
    },
    {
      path: "*",
      Component: NotFound,
      ErrorBoundary,
    },
  ]);
}
