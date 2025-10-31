import { ref } from 'vue'
import { defineStore } from 'pinia'
import { type SelectionItem } from '@/types/system'
import ax from '@/api'

export const useCoreStore = defineStore('core', () => {

  const categoriesList = ref<SelectionItem[]>([])
  const suppliersList = ref<SelectionItem[]>([])

  function loadCategoriesList() {
    var myURL = '/a/core/categories?xfields=id,name&xsort=name&xper_page=5000'
    ax.get(myURL)
    .then(response => {
      categoriesList.value = response.data.data
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

  return { categoriesList, suppliersList, loadCategoriesList, loadSuppliersList }
})

