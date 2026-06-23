import process from "node:process";
import { fileURLToPath } from "node:url";

import react from "@vitejs/plugin-react";
import { defineConfig } from "vite";

const backendUrl = process.env.BACKEND_URL || "http://backend:8080";

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  html: {
    cspNonce: "NONCEMUSTBEREPLACED",
  },
  build: {
    rolldownOptions: {
      input: {
        // The HTML entry must remain so Vite still treats this as an
        // app build and emits index.html with hashed asset references.
        main: fileURLToPath(new URL("./index.html", import.meta.url)),
        // The service worker must live at the site root with a stable
        // filename so navigator.serviceWorker.register("/sw.js") works
        // and its scope covers the entire app.
        sw: fileURLToPath(new URL("./src/sw.ts", import.meta.url)),
      },
      output: {
        // Place sw.js at the dist root with no hash; everything else
        // keeps Vite's default hashed location under /assets/.
        entryFileNames: (chunk) =>
          chunk.name === "sw" ? "[name].js" : "assets/[name]-[hash].js",
      },
    },
  },
  server: {
    proxy: {
      "^/(\\.well-known|auth|api)": {
        target: backendUrl,
      },
    },
  },
});
