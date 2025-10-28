import { ref } from 'vue'
import { defineStore } from 'pinia'
import { type Shipper } from '@/types/sales'
import { type SelectionItem } from '@/types/system'
import ax from '@/api'

export const useSalesStore = defineStore('sales', () => {

  const customersList = ref<SelectionItem[]>([])
  const regions = ['Northern', 'Eastern', 'Southern', 'Western']
  const shippersList = ref<Shipper[]>([])

  function loadCustomersList() {
    var myURL = '/a/sales/customers?xfields=id,name&xsort=name&xper_page=5000'
    ax.get(myURL)
    .then(response => {
      customersList.value = response.data.data
    })
    .catch() // handled by interceptor
  }

  function loadShippersList() {
    var myURL = '/a/sales/shippers?xfields=id,company_name&xsort=company_name&xper_page=5000'
    ax.get(myURL)
    .then(response => {
      shippersList.value = response.data.data
    })
    .catch() // handled by interceptor
  }

  return { customersList, loadCustomersList, loadShippersList, regions, shippersList }
})

