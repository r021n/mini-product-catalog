import { useEffect, useMemo, useState } from "react";
import {
  Button,
  Input,
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
} from "@heroui/react";

import { apiFetch, type SuccessEnvelope, ApiError } from "../../lib/api";
import { type Category } from "../../lib/types";
import { useAuth } from "../../lib/useAuth";

type FormMode = "create" | "edit";

export function CategoriesManager() {
  const { token } = useAuth();
  const [items, setItems] = useState<Category[]>([]);
  const [loading, setLoading] = useState(false);

  const [message, setMessage] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);

  const { isOpen, onOpen, onOpenChange, onClose } = useDisclosure();
  const [mode, setMode] = useState<FormMode>("create");
  const [active, setActive] = useState<Category | null>(null);
  const [name, setName] = useState("");

  async function load() {
    setLoading(true);
    setError(null);

    try {
      const res = await apiFetch<SuccessEnvelope<Category[]>>("/categories");
      setItems(res.data);
    } catch (e: any) {
      setError(e.message ?? "Failed to load categories");
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    load();
  }, []);

  const sorted = useMemo(() => {
    return [...items].sort((a, b) => a.name.localeCompare(b.name));
  }, [items]);

  function openCreate() {
    setMode("create");
    setActive(null);
    setName("");
    setError(null);
    setMessage(null);
    onOpen();
  }

  function openEdit(c: Category) {
    setMode("edit");
    setActive(c);
    setName(c.name);
    setError(null);
    setMessage(null);
    onOpen();
  }

  async function submit() {
    setError(null);
    setMessage(null);

    if (!token) return setError("No token (please re-login)");
    if (!name.trim()) return setError("Name is required");

    try {
      if (mode === "create") {
        await apiFetch("/categories", {
          method: "POST",
          token,
          body: { name: name.trim() },
        });
        setMessage("Category created");
      } else {
        await apiFetch(`/categories/${active!.id}`, {
          method: "PUT",
          token,
          body: { name: name.trim() },
        });
        setMessage("Category updated");
      }

      await load();
      onClose();
    } catch (e) {
      if (e instanceof ApiError) setError(e.message);
      else setError("Failed to submit");
    }
  }

  async function remove(c: Category) {
    setError(null);
    setMessage(null);

    if (!token) return setError("No token (please re-login)");
    if (!confirm(`Delete category "${c.name}"?`)) return;

    try {
      await apiFetch(`/categories/${c.id}`, { method: "DELETE", token });
      setMessage("Category deleted");
      await load();
    } catch (e) {
      if (e instanceof ApiError) setError(e.message);
      else setError("Failed to delete");
    }
  }

  return (
    <div className="space-y-3">
      <div className="flex items-center justify-between gap-3">
        <div className="font-semibold">Manage Categories</div>
        <Button color="primary" onPress={openCreate}>
          New Category
        </Button>
      </div>

      {message && <div className="text-sm text-green-700">{message}</div>}
      {error && <div className="text-sm text-red-600">{error}</div>}

      <Table aria-label="categories-table" isStriped>
        <TableHeader>
          <TableColumn>Name</TableColumn>
          <TableColumn width={220}>Actions</TableColumn>
        </TableHeader>
        <TableBody
          items={sorted}
          isLoading={loading}
          emptyContent={"No categories."}
        >
          {(c) => (
            <TableRow key={c.id}>
              <TableCell>{c.name}</TableCell>
              <TableCell>
                <div className="flex gap-2">
                  <Button size="sm" variant="flat" onPress={() => openEdit(c)}>
                    Edit
                  </Button>
                  <Button
                    size="sm"
                    color="danger"
                    variant="flat"
                    onPress={() => remove(c)}
                  >
                    Delete
                  </Button>
                </div>
              </TableCell>
            </TableRow>
          )}
        </TableBody>
      </Table>

      <Modal isOpen={isOpen} onOpenChange={onOpenChange}>
        <ModalContent>
          {(onClose) => (
            <>
              <ModalHeader>
                {mode === "create" ? "New Category" : "Edit Category"}
              </ModalHeader>
              <ModalBody className="space-y-2">
                <Input label="Name" value={name} onValueChange={setName} />
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
