import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import { RouterProvider } from "react-router/dom";

import "./index.scss";
import { createRouter } from "./router";

const queryClient = new QueryClient();
const router = createRouter(queryClient);

const $root = document.getElementById("root");
if (!$root) {
  throw new Error("Root element not found");
}

createRoot($root).render(
  <StrictMode>
    <QueryClientProvider client={queryClient}>
      <RouterProvider router={router} />
    </QueryClientProvider>
  </StrictMode>,
);

if (import.meta.env.PROD && "serviceWorker" in navigator) {
  window.addEventListener("load", () => {
    navigator.serviceWorker
      .register("/sw.js", { type: "module", scope: "/" })
      .catch((error) => {
        console.error("Service worker registration failed", error);
      });
  });
}
