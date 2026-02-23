import { useEffect, useMemo, useState } from "react";
import {
  Button,
  Input,
  Textarea,
  Modal,
  ModalContent,
  ModalHeader,
  ModalBody,
  ModalFooter,
  Table,
  TableHeader,
  TableColumn,
  TableBody,
  TableRow,
  TableCell,
  useDisclosure,
  Select,
  SelectItem,
  Pagination,
  Spinner,
} from "@heroui/react";
import { apiFetch, type SuccessEnvelope, ApiError } from "../../lib/api";
import type { Category, Product } from "../../lib/types";
import { useAuth } from "../../lib/useAuth";
import { formatIDR } from "../../lib/format";

type FormMode = "create" | "edit";

function totalPages(total: number, limit: number) {
  return Math.max(1, Math.ceil(total / limit));
}

export function ProductsManager() {
  const { token } = useAuth();

  const [categories, setCategories] = useState<Category[]>([]);
  const [items, setItems] = useState<Product[]>([]);
  const [meta, setMeta] = useState<{
    total: number;
    page: number;
    limit: number;
  } | null>(null);

  const [page, setPage] = useState(1);
  const limit = 8;

  const [loading, setLoading] = useState(false);
  const [message, setMessage] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);

  const { isOpen, onOpen, onOpenChange, onClose } = useDisclosure();
  const [mode, setMode] = useState<FormMode>("create");
  const [active, setActive] = useState<Product | null>(null);

  const [categoryID, setCategoryID] = useState<string>("all");
  const [name, setName] = useState("");
  const [description, setDescription] = useState<string>("");
  const [price, setPrice] = useState<string>("");

  async function loadCategories() {
    const res = await apiFetch<SuccessEnvelope<Category[]>>("/categories");
    setCategories(res.data);
  }

  async function loadProducts() {
    setLoading(true);
    setError(null);

    try {
      const res = await apiFetch<SuccessEnvelope<Product[]>>(
        `/products?page=${page}&limit=${limit}`,
      );
      setItems(res.data);
      setMeta(res.meta);
    } catch (e: any) {
      setError(e.message ?? "Failed to load products");
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    loadCategories().catch(() => {});
  }, []);

  useEffect(() => {
    loadProducts();
    // eslint-disable-next-line
  }, [page]);

  const pages = totalPages(meta?.total ?? 0, limit);

  const categoryOptions = useMemo(() => {
    return [...categories].sort((a, b) => a.name.localeCompare(b.name));
  }, [categories]);

  function openCreate() {
    setMode("create");
    setActive(null);
    setCategoryID(categoryOptions[0]?.id ?? "all");
    setName("");
    setDescription("");
    setPrice("");
    setError(null);
    setMessage(null);
    onOpen();
  }

  function openEdit(p: Product) {
    setMode("edit");
    setActive(p);
    setCategoryID(p.category_id);
    setName(p.name);
    setDescription(p.description);
    setPrice(String(p.price));
    setError(null);
    setMessage(null);
    onOpen();
  }

  async function submit() {
    setError(null);
    setMessage(null);

    if (!token) return setError("No token (please re-login)");
    if (!categoryID || categoryID === "all")
      return setError("Category is required");
    if (!name.trim()) return setError("Name is required");

    const priceNum = Number(price);
    if (!price || Number.isNaN(priceNum) || priceNum <= 0)
      return setError("Price must be > 0");

    const body = {
      category_id: categoryID,
      name: name.trim(),
      description: description ?? "",
      price: priceNum,
    };

    try {
      if (mode === "create") {
        await apiFetch("/products", { method: "POST", token, body });
        setMessage("Product created");
      } else {
        await apiFetch(`/products/${active!.id}`, {
          method: "PUT",
          token,
          body,
        });
        setMessage("Product updated");
      }

      await loadProducts();
      onClose();
    } catch (e) {
      if (e instanceof ApiError) setError(e.message);
      else setError("Failed to submit");
    }
  }

  async function remove(p: Product) {
    setError(null);
    setMessage(null);

    if (!token) return setError("No token (please re-login)");
    if (!confirm(`Delete product "${p.name}"?`)) return;

    try {
      await apiFetch(`/products/${p.id}`, { method: "DELETE", token });
      setMessage("Product deleted");
      await loadProducts();
    } catch (e) {
      if (e instanceof ApiError) setError(e.message);
      else setError("Failed to delete");
    }
  }

  return (
    <div className="space-y-3">
      <div className="flex items-center justify-between gap-3">
        <div className="font-semibold">Manage Products</div>
        <Button color="primary" onPress={openCreate}>
          New Product
        </Button>
      </div>

      {message && <div className="text-sm text-green-700">{message}</div>}
      {error && <div className="text-sm text-red-600">{error}</div>}

      {loading ? (
        <div className="flex justify-center py-6">
          <Spinner />
        </div>
      ) : (
        <>
          <Table aria-label="products-table" isStriped>
            <TableHeader>
              <TableColumn>Name</TableColumn>
              <TableColumn>Category</TableColumn>
              <TableColumn>Price</TableColumn>
              <TableColumn width={240}>Actions</TableColumn>
            </TableHeader>
            <TableBody items={items} emptyContent={"No products."}>
              {(p) => (
                <TableRow key={p.id}>
                  <TableCell>{p.name}</TableCell>
                  <TableCell>{p.category_name}</TableCell>
                  <TableCell className="font-mono">
                    {formatIDR(p.price)}
                  </TableCell>
                  <TableCell>
                    <div className="flex gap-2">
                      <Button
                        size="sm"
                        variant="flat"
                        onPress={() => openEdit(p)}
                      >
                        Edit
                      </Button>
                      <Button
                        size="sm"
                        color="danger"
                        variant="flat"
                        onPress={() => remove(p)}
                      >
                        Delete
                      </Button>
                    </div>
                  </TableCell>
                </TableRow>
              )}
            </TableBody>
          </Table>

          <div className="flex justify-center pt-2">
            <Pagination page={page} total={pages} onChange={setPage} />
          </div>
        </>
      )}

      <Modal isOpen={isOpen} onOpenChange={onOpenChange} size="lg">
        <ModalContent>
          {(onClose) => (
            <>
              <ModalHeader>
                {mode === "create" ? "New Product" : "Edit Product"}
              </ModalHeader>
              <ModalBody className="space-y-3">
                <Select
                  label="Category"
                  selectedKeys={new Set([categoryID])}
                  onSelectionChange={(keys) => {
                    const v = Array.from(keys)[0] as string;
                    setCategoryID(v);
                  }}
                >
                  {categoryOptions.map((c) => (
                    <SelectItem key={c.id}>{c.name}</SelectItem>
                  ))}
                </Select>

                <Input label="Name" value={name} onValueChange={setName} />
                <Textarea
                  label="Description"
                  value={description}
                  onValueChange={setDescription}
                />
                <Input
                  label="Price"
                  value={price}
                  onValueChange={setPrice}
                  placeholder="e.g. 299000"
                />
              </ModalBody>
              <ModalFooter>
                <Button variant="flat" onPress={onClose}>
                  Cancel
                </Button>
                <Button color="primary" onPress={submit}>
                  Save
                </Button>
              </ModalFooter>
            </>
          )}
        </ModalContent>
      </Modal>
    </div>
  );
}
