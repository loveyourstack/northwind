
export interface CustomerInput {
  address?: string
  city?: string
  code: string
  company_name: string
  contact_name: string
  contact_title?: string
  country_fk: number | undefined
  phone?: string
  postal_code?: string
  state?: string
}
export interface Customer extends CustomerInput {
  id: number
  country: string
  name: string
  order_count: number
}
export function NewCustomer(): Customer {
  return  {
    address: '',
    city: '',
    code: '',
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
    order_count: 0
  }
}
export function GetCustomerInputFromItem(item: Customer): CustomerInput {
  return  {
    address: item.address,
    city: item.city,
    code: item.code,
    company_name: item.company_name,
    contact_name: item.contact_name,
    contact_title: item.contact_title,
    country_fk: item.country_fk,
    phone: item.phone,
    postal_code: item.postal_code,
    state: item.state
  }
}

// ------------------------------------------------------------------------------------------------------


export interface OrderInput {
  customer_fk: number | undefined
  dest_address: string
  dest_city: string
  dest_company_name: string
  dest_country_fk: number | undefined
  dest_postal_code?: string
  dest_state?: string
  freight_cost: number
  is_shipped: boolean
  order_date: Date | undefined
  order_number: number
  required_date: Date | undefined
  salesman_fk: number | undefined
  shipped_date?: Date | undefined
  shipper_fk?: number | undefined
}
export interface Order extends OrderInput {
  id: number
  customer_code: string
  customer_company_name: string
  dest_country_iso2: string
  salesman: string
  shipper_company_name: string
  order_detail_count: number
  order_value: number
}
export function NewOrder(): Order {
  return  {
    customer_fk: undefined,
    dest_address: '',
    dest_city: '',
    dest_company_name: '',
    dest_country_fk: undefined,
    dest_postal_code: '',
    dest_state: '',
    freight_cost: 0,
    is_shipped: false,
    order_date: undefined,
    order_number: 0,
    required_date: undefined,
    salesman_fk: undefined,
    shipped_date: undefined,
    shipper_fk: undefined,
    id: 0,
    customer_code: '',
    customer_company_name: '',
    dest_country_iso2: '',
    salesman: '',
    shipper_company_name: '',
    order_detail_count: 0,
    order_value: 0,
  }
}
export function GetOrderInputFromItem(item: Order): OrderInput {
  return  {
    customer_fk: item.customer_fk,
    dest_address: item.dest_address,
    dest_city: item.dest_city,
    dest_company_name: item.dest_company_name,
    dest_country_fk: item.dest_country_fk,
    dest_postal_code: item.dest_postal_code,
    dest_state: item.dest_state,
    freight_cost: item.freight_cost,
    is_shipped: item.is_shipped,
    order_date: item.order_date,
    order_number: item.order_number,
    required_date: item.required_date,
    salesman_fk: item.salesman_fk,
    shipped_date: item.shipped_date,
    shipper_fk: item.shipper_fk,
  }
}


// ------------------------------------------------------------------------------------------------------

export interface OrderDetail {
  discount: number
  order_fk: number
  order_number: string
  product_fk: number
  product_name: string
  quantity: number
  unit_price: number
}

// ------------------------------------------------------------------------------------------------------

export interface Shipper {
  id: number
  company_name: string
  phone: string
}