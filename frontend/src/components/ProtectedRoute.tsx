import { Navigate, Outlet, useLocation } from "react-router-dom";
import { useAuth } from "../lib/useAuth";

export function ProtectedRoute({ requireAdmin }: { requireAdmin?: boolean }) {
  const { user, loading } = useAuth();
  const loc = useLocation();

  if (loading) return null;

  if (!user)
    return <Navigate to="/login" replace state={{ from: loc.pathname }} />;

  if (requireAdmin && user.role !== "admin") return <Navigate to="/" replace />;

  return <Outlet />;
}
