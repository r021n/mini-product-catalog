import { useEffect, useState } from "react";
import { Card, CardBody } from "@heroui/react";
import { apiFetch, type SuccessEnvelope } from "../lib/api";

type Product = {
  id: string;
  category_id: string;
  category_name: string;
  name: string;
  description: string;
  price: number;
  created_at: string;
  updated_at: string;
};

export function CatalogPage() {
  const [items, setItems] = useState<Product[]>([]);
  const [meta, setMeta] = useState<any>(null);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    apiFetch<SuccessEnvelope<Product[]>>("/products?page=1&limit=10")
      .then((res) => {
        setItems(res.data);
        setMeta(res.meta);
      })
      .catch((e) => setError(e.message ?? "Failed to load products"));
  }, []);

  return (
    <div className="space-y-4">
      <div className="text-xl font-semibold">Catalog</div>

      {meta && (
        <div className="text-sm text-slate-600">Total: {meta.total}</div>
      )}
      {error && <div className="text-sm text-red-600">{error}</div>}

      <div className="grid md:grid-cols-2 gap-4">
        {items.map((p) => (
          <Card key={p.id}>
            <CardBody className="space-y-2">
              <div className="font-semibold">{p.name}</div>
              <div className="text-sm text-slate-600">{p.category_name}</div>
              <div className="text-sm">{p.description}</div>
              <div className="font-mono text-sm">Rp {p.price}</div>
            </CardBody>
          </Card>
        ))}
      </div>
    </div>
  );
}
