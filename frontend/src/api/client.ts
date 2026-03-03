export class APIError extends Error {
  public response: Response;

  constructor(res: Response) {
    super(`${res.status} (${res.statusText})`);
    this.response = res;
  }
}

export async function get<T>(path: string): Promise<T> {
  const res = await request(path);
  return res.json() as Promise<T>;
}

export async function request(path: string, options?: RequestInit): Promise<Response> {
  let res = await fetch(path, options);
  if (res.status !== 401) {
    return res;
  }

  res = await fetch("/auth/refresh");
  if (!res.ok) {
    await res.text();
    throw new APIError(res);
  }

  res = await fetch(path, options);
  if (!res.ok) {
    await res.text();
    throw new APIError(res);
  }

  return res;
}
