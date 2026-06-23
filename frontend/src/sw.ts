/// <reference lib="webworker" />

declare const self: ServiceWorkerGlobalScope;

const CACHE_VERSION = "v1";
const APP_SHELL_CACHE = `goweb-shell-${CACHE_VERSION}`;
const APP_SHELL: readonly string[] = [
  "/",
  "/index.html",
  "/favicon.svg",
  "/manifest.webmanifest",
];
const ASSET_CACHE = `goweb-assets-${CACHE_VERSION}`;

const isBackendPath = (pathname: string): boolean =>
  pathname.startsWith("/api/") ||
  pathname.startsWith("/auth/") ||
  pathname.startsWith("/.well-known/");

self.addEventListener("install", (event) => {
  event.waitUntil(
    (async () => {
      const cache = await caches.open(APP_SHELL_CACHE);
      await cache.addAll(APP_SHELL);
      await self.skipWaiting();
    })(),
  );
});

self.addEventListener("activate", (event) => {
  event.waitUntil(
    (async () => {
      const keys = await caches.keys();
      await Promise.all(
        keys
          .filter((key) => key !== APP_SHELL_CACHE && key !== ASSET_CACHE)
          .map((key) => caches.delete(key)),
      );
      await self.clients.claim();
    })(),
  );
});

self.addEventListener("fetch", (event) => {
  const { request } = event;

  // Only cache GET requests.
  if (request.method !== "GET") {
    return;
  }

  // Only cache same-origin requests.
  const url = new URL(request.url);
  if (url.origin !== self.location.origin) {
    return;
  }

  // Don't cache backend requests.
  if (isBackendPath(url.pathname)) {
    return;
  }

  if (request.mode === "navigate") {
    event.respondWith(
      (async () => {
        try {
          const res = await fetch(request);
          const cache = await caches.open(APP_SHELL_CACHE);
          cache.put("/index.html", res.clone());
          return res;
        } catch {
          const res = await caches.match("/index.html");
          if (res) return res;
          return new Response("Offline", {
            status: 503,
            statusText: "Service Unavailable",
            headers: { "Content-Type": "text/plain" },
          });
        }
      })(),
    );
    return;
  }

  if (url.pathname.startsWith("/assets/")) {
    event.respondWith(
      (async () => {
        const cache = await caches.open(ASSET_CACHE);
        let res = await cache.match(request);
        if (res) {
          return res;
        }

        res = await fetch(request);
        if (res.ok) {
          cache.put(request, res.clone());
        }

        return res;
      })(),
    );
    return;
  }

  event.respondWith(
    (async () => {
      const cache = await caches.open(APP_SHELL_CACHE);
      let res = await cache.match(request);
      if (res) {
        return res;
      }

      res = await fetch(request);
      if (res.ok) {
        cache.put(request, res.clone());
      }

      return res;
    })(),
  );
});

self.addEventListener("message", (event) => {
  if (event.data === "SKIP_WAITING") {
    self.skipWaiting();
  }
});
