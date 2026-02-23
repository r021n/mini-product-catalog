import { Tabs, Tab, Card, CardBody } from "@heroui/react";
import { CategoriesManager } from "./admin/CategoriesManager";
import { ProductsManager } from "./admin/ProductsManager";

export function AdminPage() {
  return (
    <Card>
      <CardBody className="space-y-4">
        <div className="text-xl font-semibold">Admin Panel</div>
        <Tabs aria-label="admin-tabs">
          <Tab key="categories" title="Categories">
            <CategoriesManager />
          </Tab>
          <Tab key="products" title="Products">
            <ProductsManager />
          </Tab>
        </Tabs>
      </CardBody>
    </Card>
  );
}
