import React, { useEffect, useMemo, useState } from "react";
import { apiFetch, type SuccessEnvelope } from "./api";
import { AuthContext, type User, type AuthContextValue } from "./AuthContext";

const TOKEN_KEY = "access_token";

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [token, setToken] = useState<string | null>(() =>
    localStorage.getItem(TOKEN_KEY),
  );
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);

  async function refreshMe() {
    if (!token) {
      setUser(null);
      return;
    }
    const res = await apiFetch<SuccessEnvelope<User>>("/me", { token });
    setUser(res.data);
  }

  useEffect(() => {
    (async () => {
      try {
        await refreshMe();
      } catch {
        localStorage.removeItem(TOKEN_KEY);
        setToken(null);
        setUser(null);
      } finally {
        setLoading(false);
      }
    })();
    // eslint-disable-next-line
  }, []);

  async function login(email: string, password: string) {
    const res = await apiFetch<SuccessEnvelope<{ access_token: string }>>(
      "/auth/login",
      {
        method: "POST",
        body: { email, password },
      },
    );

    const t = res.data.access_token;
    localStorage.setItem(TOKEN_KEY, t);
    setToken(t);

    await refreshMe();
  }

  async function register(name: string, email: string, password: string) {
    await apiFetch("/auth/register", {
      method: "POST",
      body: { name, email, password },
    });

    await login(email, password);
  }

  function logout() {
    localStorage.removeItem(TOKEN_KEY);
    setToken(null);
    setUser(null);
  }

  const value = useMemo<AuthContextValue>(
    () => ({
      token,
      user,
      loading,
      login,
      register,
      logout,
      refreshMe,
    }),
    // eslint-disable-next-line
    [token, user, loading],
  );

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}
