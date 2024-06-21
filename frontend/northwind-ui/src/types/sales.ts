import { useDateFormat, useNow } from "@vueuse/core"

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
  order_date: string | undefined // string, not Date: use additional Date object with _d suffix for date picker element
  order_number: number
  required_date: string | undefined
  salesman_fk: number | undefined
  shipped_date?: string | undefined
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

  // date objects: to be assigned after load
  order_date_d: Date | undefined
  required_date_d: Date | undefined
  shipped_date_d: Date | undefined
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

    order_date_d: useNow().value, // new order date defaults to today
    required_date_d: undefined,
    shipped_date_d: undefined
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
    order_date: useDateFormat(item.order_date_d, 'YYYY-MM-DD').value, // format corresponding Date object for db input
    order_number: item.order_number,
    required_date: useDateFormat(item.required_date_d, 'YYYY-MM-DD').value,
    salesman_fk: item.salesman_fk,
    shipped_date: useDateFormat(item.shipped_date_d, 'YYYY-MM-DD').value,
    shipper_fk: item.shipper_fk,
  }
}


// ------------------------------------------------------------------------------------------------------

export interface OrderDetailInput {
  discount: number
  order_fk: number | undefined
  product_fk: number | undefined
  quantity: number
  unit_price: number
}
export interface OrderDetail extends OrderDetailInput  {
  id: number
  order_number: string
  product_name: string
}
export function NewOrderDetail(order_id: number): OrderDetail {
  return  {
    discount: 0,
    order_fk: order_id,
    product_fk: undefined,
    quantity: 1,
    unit_price: 0,
    id: 0,
    order_number: '',
    product_name: '',
  }
}
export function GetOrderDetailInputFromItem(item: OrderDetail): OrderDetailInput {
  return  {
    discount: item.discount,
    order_fk: item.order_fk,
    product_fk: item.product_fk,
    quantity: item.quantity,
    unit_price: item.unit_price,
  }
}

// ------------------------------------------------------------------------------------------------------

export interface Shipper {
  id: number
  company_name: string
  phone: string
}