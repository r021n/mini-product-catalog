import { useEffect, useState } from "react";
import { useParams, Link } from "react-router-dom";
import { Card, CardBody, Spinner, Button } from "@heroui/react";
import { apiFetch, type SuccessEnvelope } from "../lib/api";
import { type Product } from "../lib/types";
import { formatIDR } from "../lib/format";

export function ProductDetailPage() {
  const { id } = useParams();
  const [item, setItem] = useState<Product | null>(null);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!id) return;

    let ignore = false;

    const fetchProduct = async () => {
      // Reset state sebelum mulai fetch
      setError(null);
      setItem(null);

      try {
        const res = await apiFetch<SuccessEnvelope<Product>>(`/products/${id}`);
        if (!ignore) {
          setItem(res.data);
        }
      } catch (e: any) {
        if (!ignore) {
          setError(e.message ?? "Failed to load product");
        }
      }
    };

    fetchProduct();

    return () => {
      ignore = true;
    };
  }, [id]);

  if (error)
    return (
      <Card>
        <CardBody className="space-y-3">
          <div className="text-red-600 text-sm">{error}</div>
          <Button as={Link} to="/" variant="flat">
            Back to catalog
          </Button>
        </CardBody>
      </Card>
    );

  if (!item)
    return (
      <div className="flex justify-center py-10">
        <Spinner />
      </div>
    );

  return (
    <Card>
      <CardBody className="space-y-2">
        <div className="text-xs text-slate-500">{item.category_name}</div>
        <div className="text-xl font-semibold">{item.name}</div>
        <div className="text-sm text-slate-700">{item.description}</div>
        <div className="font-mono">{formatIDR(item.price)}</div>

        <div className="pt-2">
          <Button as={Link} to="/" variant="flat">
            Back
          </Button>
        </div>
      </CardBody>
    </Card>
  );
}
