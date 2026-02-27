const API_BASE =
  typeof window === "undefined"
    ? process.env.API_SERVER_BASE_URL || process.env.NEXT_PUBLIC_API_BASE_URL || "http://backend:8080/api/v1"
    : process.env.NEXT_PUBLIC_API_BASE_URL || "/api/v1";

export type AuthUser = {
  id: number;
  email: string;
  name: string;
  avatar_url?: string;
  provider: string;
  created_at: string;
};

type AuthEnvelope = {
  data: AuthUser;
};

async function authFetch<T>(path: string, init?: RequestInit): Promise<T> {
  const res = await fetch(`${API_BASE}${path}`, {
    credentials: "include",
    ...init,
  });
  if (!res.ok) {
    const body = await res.json().catch(() => ({}));
    throw new Error(body.error || `API ${res.status}`);
  }
  const body = await res.json();
  return body.data;
}

export async function register(email: string, password: string, name: string): Promise<AuthUser> {
  return authFetch<AuthUser>("/auth/register", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ email, password, name }),
  });
}

export async function login(email: string, password: string): Promise<AuthUser> {
  return authFetch<AuthUser>("/auth/login", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ email, password }),
  });
}

export async function fetchMe(): Promise<AuthUser | null> {
  try {
    return await authFetch<AuthUser>("/auth/me");
  } catch {
    return null;
  }
}

export async function logout(): Promise<void> {
  await fetch(`${API_BASE}/auth/logout`, {
    method: "POST",
    credentials: "include",
  });
}

export function googleLoginUrl(): string {
  const base =
    typeof window !== "undefined"
      ? process.env.NEXT_PUBLIC_API_BASE_URL || "/api/v1"
      : "";
  return `${base}/auth/google`;
}
