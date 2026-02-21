import { useEffect, useMemo, useState } from "react";
import { Link } from "react-router-dom";
import {
  Card,
  CardBody,
  Input,
  Button,
  Spinner,
  Pagination,
  Select,
  SelectItem,
} from "@heroui/react";
import { apiFetch, type SuccessEnvelope } from "../lib/api";
import { type Category, type Product } from "../lib/types";
import { formatIDR } from "../lib/format";

function totalPages(total: number, limit: number) {
  return Math.max(1, Math.ceil(total / limit));
}

export function CatalogPage() {
  const [categories, setCategories] = useState<Category[]>([]);

  const [q, setQ] = useState("");
  const [categoryID, setCategoryID] = useState<string>("all");
  const [minPrice, setMinPrice] = useState<string>("");
  const [maxPrice, setMaxPrice] = useState<string>("");
  const [sort, setSort] = useState<string>("created_at");
  const [order, setOrder] = useState<string>("desc");

  const [page, setPage] = useState(1);
  const [limit] = useState(6);

  const [items, setItems] = useState<Product[]>([]);
  const [meta, setMeta] = useState<{
    total: number;
    page: number;
    limit: number;
  } | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const categoryOptions = [{ id: "all", name: "All" }, ...categories];

  useEffect(() => {
    apiFetch<SuccessEnvelope<Category[]>>("/categories")
      .then((res) => setCategories(res.data))
      .catch(() => {});
  }, []);

  const queryString = useMemo(() => {
    const p = new URLSearchParams();
    p.set("page", String(page));
    p.set("limit", String(limit));

    if (q.trim()) p.set("q", q.trim());
    if (categoryID !== "all") p.set("category_id", categoryID);
    if (minPrice.trim()) p.set("min_price", minPrice.trim());
    if (maxPrice.trim()) p.set("max_price", maxPrice.trim());

    p.set("sort", sort);
    p.set("order", order);

    return p.toString();
  }, [page, limit, q, categoryID, minPrice, maxPrice, sort, order]);

  useEffect(() => {
    let isMounted = true;

    const fetchProducts = async () => {
      await Promise.resolve();

      if (!isMounted) return;
      setLoading(true);
      setError(null);

      try {
        const res = await apiFetch<SuccessEnvelope<Product[]>>(
          `/products?${queryString}`,
        );
        if (isMounted) {
          setItems(res.data);
          setMeta(res.meta);
        }
      } catch (e: any) {
        if (isMounted) setError(e.message ?? "Failed to load products");
      } finally {
        if (isMounted) setLoading(false);
      }
    };

    fetchProducts();

    return () => {
      isMounted = false;
    };
  }, [queryString]);

  function resetToFirstPage() {
    setPage(1);
  }

  const pages = totalPages(meta?.total ?? 0, limit);

  return (
    <div className="space-y-4">
      <div className="text-xl font-semibold">Catalog</div>

      <Card>
        <CardBody className="grid md:grid-cols-3 gap-3">
          <Input
            label="Search"
            placeholder="e.g. mouse"
            value={q}
            onValueChange={(v) => {
              setQ(v);
              resetToFirstPage();
            }}
          />

          <Select
            label="Category"
            items={categoryOptions}
            selectedKeys={new Set([categoryID])}
            onSelectionChange={(keys) => {
              const v = Array.from(keys)[0] as string;
              setCategoryID(v);
              resetToFirstPage();
            }}
          >
            {(c) => <SelectItem key={c.id}>{c.name}</SelectItem>}
          </Select>

          <div className="grid grid-cols-2 gap-3">
            <Input
              label="Min price"
              placeholder="100000"
              value={minPrice}
              onValueChange={(v) => {
                setMinPrice(v);
                resetToFirstPage();
              }}
            />
            <Input
              label="Max price"
              placeholder="500000"
              value={maxPrice}
              onValueChange={(v) => {
                setMaxPrice(v);
                resetToFirstPage();
              }}
            />
          </div>

          <Select
            label="Sort"
            selectedKeys={new Set([sort])}
            onSelectionChange={(keys) => {
              const v = Array.from(keys)[0] as string;
              setSort(v);
              resetToFirstPage();
            }}
          >
            <SelectItem key="created_at">created_at</SelectItem>
            <SelectItem key="price">price</SelectItem>
          </Select>

          <Select
            label="Order"
            selectedKeys={new Set([order])}
            onSelectionChange={(keys) => {
              const v = Array.from(keys)[0] as string;
              setOrder(v);
              resetToFirstPage();
            }}
          >
            <SelectItem key="desc">desc</SelectItem>
            <SelectItem key="asc">asc</SelectItem>
          </Select>

          <Button
            variant="flat"
            onPress={() => {
              // basic validation min/max numeric
              if (minPrice && isNaN(Number(minPrice)))
                return setError("Min price must be a number");
              if (maxPrice && isNaN(Number(maxPrice)))
                return setError("Max price must be a number");
              setError(null);
            }}
          >
            Apply
          </Button>
        </CardBody>
      </Card>

      {meta && (
        <div className="text-sm text-slate-600">
          Total: <span className="font-semibold">{meta.total}</span>
        </div>
      )}
      {error && <div className="text-sm text-red-600">{error}</div>}

      {loading ? (
        <div className="flex justify-center py-10">
          <Spinner />
        </div>
      ) : items.length === 0 ? (
        <Card>
          <CardBody>No products found.</CardBody>
        </Card>
      ) : (
        <div className="grid md:grid-cols-2 gap-4">
          {items.map((p) => (
            <Card key={p.id} className="hover:shadow-md transition">
              <CardBody className="space-y-2">
                <div className="text-xs text-slate-500">{p.category_name}</div>

                <Link
                  to={`/products/${p.id}`}
                  className="font-semibold underline-offset-2 hover:underline"
                >
                  {p.name}
                </Link>

                <div className="text-sm text-slate-700 line-clamp-2">
                  {p.description}
                </div>
                <div className="font-mono text-sm">{formatIDR(p.price)}</div>
              </CardBody>
            </Card>
          ))}
        </div>
      )}

      <div className="flex justify-center pt-2">
        <Pagination page={page} total={pages} onChange={setPage} />
      </div>
    </div>
  );
}
