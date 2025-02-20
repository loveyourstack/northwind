import { ref } from 'vue'
import { defineStore } from 'pinia'
import { SelectionItem } from '@/types/system'
import ax from '@/api'

export const useCoreStore = defineStore('core', () => {

  const booleanOptions = [
    {value: true, title: 'Yes'},
    {value: false, title: 'No'}
  ]

  const categoriesList = ref<SelectionItem[]>([])
  const productsList = ref<SelectionItem[]>([])
  const suppliersList = ref<SelectionItem[]>([])

  function loadCategoriesList() {
    var myURL = '/a/core/categories?xfields=id,name&xsort=name&xper_page=5000'
    ax.get(myURL)
    .then(response => {
      categoriesList.value = response.data.data
    })
    .catch() // handled by interceptor
  }

  function loadProductsList() {
    var myURL = '/a/core/products?xfields=id,name&xsort=name&xper_page=5000'
    ax.get(myURL)
    .then(response => {
      productsList.value = response.data.data
    })
    .catch() // handled by interceptor
  }

  function loadSuppliersList() {
    var myURL = '/a/core/suppliers?xfields=id,name&xsort=name&xper_page=5000'
    ax.get(myURL)
    .then(response => {
      suppliersList.value = response.data.data
    })
    .catch() // handled by interceptor
  }

  return { booleanOptions, categoriesList, suppliersList, loadCategoriesList, loadSuppliersList, loadProductsList, productsList }
})

