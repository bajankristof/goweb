import { fileURLToPath } from "node:url";

import react from "@vitejs/plugin-react";
import { defineConfig } from "vite";

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  html: {
    cspNonce: "NONCEPLACEHOLDER",
  },
  css: {
    preprocessorOptions: {
      scss: {
        loadPaths: [fileURLToPath(new URL("./src/styles", import.meta.url))],
      },
    },
  },
  resolve: {
    alias: {
      "#components": fileURLToPath(
        new URL("./src/components", import.meta.url),
      ),
      "#hooks": fileURLToPath(new URL("./src/hooks", import.meta.url)),
      "#pages": fileURLToPath(new URL("./src/pages", import.meta.url)),
    },
  },
  build: {
    rolldownOptions: {
      input: {
        main: fileURLToPath(new URL("./index.html", import.meta.url)),
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
});
