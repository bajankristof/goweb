import { createBrowserRouter } from "react-router";

import DefaultLayout, { defaultLayoutLoader } from "./layouts/DefaultLayout";
import SignIn, { signInLoader } from "./pages/SignIn";
import Home from "./pages/Home";

const router = createBrowserRouter([
  {
    path: "signin",
    loader: signInLoader,
    Component: SignIn,
  },
  {
    path: "/",
    loader: defaultLayoutLoader,
    Component: DefaultLayout,
    children: [{ index: true, Component: Home }],
  },
]);

export default router;
