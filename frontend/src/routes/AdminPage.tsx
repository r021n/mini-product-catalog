import { Card, CardBody } from "@heroui/react";
import { useAuth } from "../lib/useAuth";

export function AdminPage() {
  const { user } = useAuth();

  return (
    <Card>
      <CardBody className="space-y-2">
        <div className="text-xl font-semibold">Admin</div>
        <div className="text-sm text-slate-600">
          Welcome, {user?.email} (role: {user?.role})
        </div>
        <div className="text-sm">
          Next chapter kita isi halaman ini untuk CRUD categories & products.
        </div>
      </CardBody>
    </Card>
  );
}
