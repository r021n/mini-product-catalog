export type Category = {
  id: string;
  name: string;
  created_at: string;
};

export type Product = {
  id: string;
  category_id: string;
  category_name: string;
  name: string;
  description: string;
  price: number;
  created_at: string;
  updated_at: string;
};
