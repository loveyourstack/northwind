import { ref } from 'vue'
import { defineStore } from 'pinia'
import { type Country } from '@/types/core'
import { type SelectionItem } from '@/types/system'
import ax from '@/api'

export const useCoreStore = defineStore('core', () => {

  const activeCountriesList = ref<Country[]>([])
  const countriesList = ref<Country[]>([])
  const categoriesList = ref<SelectionItem[]>([])
  const frequenciesList = ref<string[]>(['Daily', 'Weekly'])
  const operatorsList = <string[]>(['<', '<=', '=', '>=', '>', '<=>', '!=' ])
  const productsList = ref<SelectionItem[]>([])
  const suppliersList = ref<SelectionItem[]>([])
  const weekdaysList = ref<string[]>(['None', 'Sunday', 'Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday'])

  function loadCategoriesList() {
    var myURL = '/a/core/categories?xfields=id,name&xsort=name&xper_page=5000'
    ax.get(myURL)
    .then(response => {
      categoriesList.value = response.data.data
    })
    .catch() // handled by interceptor
  }

  function loadCountriesList() {
    var myURL = '/a/core/countries?xfields=id,is_active,iso2,name&xsort=name&xper_page=5000'
    ax.get(myURL)
    .then(response => {
      countriesList.value = response.data.data
      activeCountriesList.value = countriesList.value.filter((el: Country) => el.is_active)
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

  return { activeCountriesList, categoriesList, countriesList, frequenciesList, operatorsList, productsList, suppliersList, weekdaysList,
    loadCategoriesList, loadCountriesList, loadSuppliersList, loadProductsList  }
})

