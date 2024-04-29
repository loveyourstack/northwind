
export interface CategoryInput {
  color_hex: string
  description: string
  name: string
}
export interface Category extends CategoryInput {
  id: number
  active_product_count: number
}
export function NewCategory(): Category {
  return  {
    color_hex: '',
    description: '',
    name: '',
    id: 0,
    active_product_count: 0
  }
}
export function GetCategoryInputFromItem(item: Category): CategoryInput {
  return  {
    color_hex: item.color_hex,
    description: item.description,
    name: item.name,
  }
}

// ------------------------------------------------------------------------------------------------------

export interface ProductInput {
  category_fk: number | undefined
  is_discontinued?: boolean
  name: string
  quantity_per_unit: string
  reorder_level?: number
  supplier_fk: number | undefined
  tags?: string[]
  unit_price: number
  units_in_stock?: number
  units_on_order?: number
}
export interface Product extends ProductInput {
  id: number
  category: string
  category_color_hex: string
  supplier_company_name: string
}
export function NewProduct(): Product {
  return  {
    category_fk: undefined,
    is_discontinued: false,
    name: '',
    quantity_per_unit: '',
    supplier_fk: undefined,
    unit_price: 0,
    id: 0,
    category: '',
    category_color_hex: '',
    supplier_company_name: ''
  }
}
export function GetProductInputFromItem(item: Product): ProductInput {

  if (!item.tags) {
    item.tags = []
  }
  return  {
    category_fk: item.category_fk,
    is_discontinued: item.is_discontinued,
    name: item.name,
    quantity_per_unit: item.quantity_per_unit,
    reorder_level: item.reorder_level,
    supplier_fk: item.supplier_fk,
    tags: item.tags,
    unit_price: item.unit_price,
    units_in_stock: item.units_in_stock,
    units_on_order: item.units_on_order,
  }
}

// ------------------------------------------------------------------------------------------------------

export interface SupplierInput {
  address?: string
  city?: string
  company_name: string
  contact_name: string
  contact_title?: string
  country_fk: number | undefined
  phone?: string
  postal_code?: string
  state?: string
}
export interface Supplier extends SupplierInput {
  id: number
  country: string
  name: string
  active_product_count: number
}
export function NewSupplier(): Supplier {
  return  {
    address: '',
    city: '',
    company_name: '',
    contact_name: '',
    contact_title: '',
    country_fk: undefined,
    phone: '',
    postal_code: '',
    state: '',
    id: 0,
    country: '',
    name: '',
    active_product_count: 0
  }
}
export function GetSupplierInputFromItem(item: Supplier): SupplierInput {
  return  {
    address: item.address,
    city: item.city,
    company_name: item.company_name,
    contact_name: item.contact_name,
    contact_title: item.contact_title,
    country_fk: item.country_fk,
    phone: item.phone,
    postal_code: item.postal_code,
    state: item.state
  }
}