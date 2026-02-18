export type ApiErrorEnvelope = {
  error: { message: string; details?: any };
};

export class ApiError extends Error {
  status: number;
  details?: any;

  constructor(status: number, message: string, details?: any) {
    super(message);
    this.status = status;
    this.details = details;
  }
}

const API_URL = import.meta.env.VITE_API_URL as string;

type RequestOptions = {
  method?: string;
  token?: string | null;
  body?: any;
};

export async function apiFetch<T>(
  path: string,
  opts: RequestOptions = {},
): Promise<T> {
  const headers: Record<string, string> = {
    "Content-Type": "application/json",
  };

  if (opts.token) {
    headers["Authorization"] = `Bearer ${opts.token}`;
  }

  const res = await fetch(`${API_URL}${path}`, {
    method: opts.method ?? "GET",
    headers,
    body: opts.body ? JSON.stringify(opts.body) : undefined,
  });

  const text = await res.text();
  const data = text ? JSON.parse(text) : null;

  if (!res.ok) {
    const env = data as ApiErrorEnvelope | null;
    const msg =
      env?.error?.message ?? `Request failed with status ${res.status}`;
    throw new ApiError(res.status, msg, env?.error?.details);
  }

  return data as T;
}

export type SuccessEnvelope<T> = { data: T; meta?: any };
