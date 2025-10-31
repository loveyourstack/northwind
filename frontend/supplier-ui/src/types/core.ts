
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
  category_color_is_light: boolean
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
    category_color_is_light: false,
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

export interface NumericFilter {
  operator: string
  value: number
  value_upper: number // treated as upper limit when operator is 'between'
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