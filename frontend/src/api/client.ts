export class APIError extends Error {
  public readonly response: Response;

  constructor(response: Response) {
    super(`${response.status} (${response.statusText})`);
    this.name = "APIError";
    this.response = response;
  }

  get status(): number {
    return this.response.status;
  }

  get statusText(): string {
    return this.response.statusText;
  }

  get url(): string {
    return this.response.url;
  }
}

export class NetworkError extends Error {
  constructor(
    message = "Network request failed",
    options?: { cause?: unknown },
  ) {
    super(message, options);
    this.name = "NetworkError";
  }
}

export type RequestOptions = RequestInit & {
  retryOnUnauthorized?: boolean;
};

let refresh: Promise<void> | null = null;

async function neofetch(
  input: RequestInfo | URL,
  init?: RequestInit,
): Promise<Response> {
  let res: Response;

  try {
    res = await fetch(input, init);
  } catch (cause) {
    throw new NetworkError("Network request failed", { cause });
  }

  if (!res.ok) {
    throw new APIError(res);
  }

  return res;
}

export async function request(
  path: string,
  { retryOnUnauthorized = true, ...init }: RequestOptions = {},
): Promise<Response> {
  try {
    return await neofetch(path, init);
  } catch (error) {
    if (
      !(error instanceof APIError) ||
      error.status !== 401 ||
      !retryOnUnauthorized
    ) {
      throw error;
    }
  }

  if (!refresh) {
    refresh = neofetch("/auth/refresh", { method: "POST" })
      .then(() => {})
      .finally(() => (refresh = null));
  }

  await refresh;

  return neofetch(path, init);
}

export async function requestJSON<T>(
  path: string,
  init: RequestOptions = {},
): Promise<T> {
  const headers = new Headers(init?.headers);
  if (!headers.has("Accept")) {
    headers.set("Accept", "application/json");
  }

  if (init?.body !== undefined && !headers.has("Content-Type")) {
    headers.set("Content-Type", "application/json");
  }

  const res = await request(path, { ...init, headers });
  if (res.status === 204 || res.headers.get("Content-Length") === "0") {
    return undefined as T;
  }

  return (await res.json()) as T;
}
