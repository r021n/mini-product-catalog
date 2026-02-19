import { Link } from "react-router-dom";
import {
  Navbar,
  NavbarBrand,
  NavbarContent,
  NavbarItem,
  Button,
  Chip,
} from "@heroui/react";
import { useAuth } from "../lib/useAuth";
import type React from "react";

export function AppShell({ children }: { children: React.ReactNode }) {
  const { user, logout } = useAuth();

  return (
    <div className="min-h-screen bg-slate-50">
      <Navbar>
        <NavbarBrand>
          <Link to="/" className="font-semibold">
            Mini Product Catalog
          </Link>
        </NavbarBrand>

        <NavbarContent justify="end">
          {!user ? (
            <>
              <NavbarItem>
                <Link to="/login">Login</Link>
              </NavbarItem>
              <NavbarItem>
                {" "}
                <Link to="/register">Register</Link>
              </NavbarItem>
            </>
          ) : (
            <>
              <NavbarItem className="flex items-center gap-2">
                <Chip
                  size="sm"
                  color={user.role === "admin" ? "danger" : "primary"}
                >
                  {user.role}
                </Chip>
                <span className="text-sm text-slate-700">{user.email}</span>
              </NavbarItem>
              {user.role === "admin" && (
                <NavbarItem>
                  <Link to="/admin">Admin</Link>
                </NavbarItem>
              )}
              <NavbarItem>
                <Button size="sm" variant="flat" onPress={logout}>
                  Logout
                </Button>
              </NavbarItem>
            </>
          )}
        </NavbarContent>
      </Navbar>

      <main className="max-w-4xl mx-auto p-6">{children}</main>
    </div>
  );
}
