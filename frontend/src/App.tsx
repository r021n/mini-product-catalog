import { AppShell } from "./components/AppShell";
import { ProtectedRoute } from "./components/ProtectedRoute";

import { LoginPage } from "./routes/LoginPage";
import { RegisterPage } from "./routes/RegisterPage";
import { CatalogPage } from "./routes/CatalogPage";
import { AdminPage } from "./routes/AdminPage";
import { ProductDetailPage } from "./routes/ProductDetailPage";
import { Route, Routes } from "react-router-dom";

export default function App() {
  return (
    <AppShell>
      <Routes>
        <Route path="/" element={<CatalogPage />} />
        <Route path="/login" element={<LoginPage />} />
        <Route path="/register" element={<RegisterPage />} />
        <Route path="/products/:id" element={<ProductDetailPage />} />

        <Route element={<ProtectedRoute />}></Route>

        <Route element={<ProtectedRoute requireAdmin />}>
          <Route path="/admin" element={<AdminPage />} />
        </Route>

        <Route path="*" element={<div>Not Found</div>} />
      </Routes>
    </AppShell>
  );
}
