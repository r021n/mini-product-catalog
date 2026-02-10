import { useEffect, useState } from "react";
import { Navbar, NavbarBrand, Card, CardBody, Chip } from "@heroui/react";

type Health = { status: string };

export default function App() {
  const [health, setHealth] = useState<Health | null>(null);
  const [error, setError] = useState<string | null>(null);

  const apiUrl = import.meta.env.VITE_API_URL as string;

  useEffect(() => {
    fetch(`${apiUrl}/health`)
      .then(async (r) => {
        if (!r.ok) throw new Error(`HTTP ${r.status}`);
        return (await r.json()) as Health;
      })
      .then(setHealth)
      .catch((e) => setError(String(e)));
  }, [apiUrl]);

  return (
    <div className="min-h-screen bg-slate-50">
      <Navbar>
        <NavbarBrand className="font-semibold">
          Mini Product Catalog
        </NavbarBrand>
      </Navbar>

      <main className="max-w-3xl mx-auto p-6">
        <Card>
          <CardBody className="space-y-3">
            <div className="text-lg font-semibold">Home</div>

            <div className="flex items-center gap-2">
              <div>Backend health:</div>
              {health && <Chip color="success">{health.status}</Chip>}
              {error && <Chip color="danger">{error}</Chip>}
              {!health && !error && <Chip>loading...</Chip>}
            </div>

            <div className="text-sm text-slate-600">
              API: <span className="font-mono">{apiUrl}</span>
            </div>
          </CardBody>
        </Card>
      </main>
    </div>
  );
}
