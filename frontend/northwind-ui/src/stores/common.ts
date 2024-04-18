import { ref } from 'vue'
import { defineStore } from 'pinia'
import { Country } from '@/types/common'
import ax from '@/api'

export const useCommonStore = defineStore('common', () => {

  const countriesList = ref<Country[]>([])
  const activeCountriesList = ref<Country[]>([])

  function loadCountriesList() {
    var myURL = '/a/common/countries?xfields=id,is_active,iso2,name&xsort=name&xper_page=5000'
    ax.get(myURL)
    .then(response => {
      countriesList.value = response.data.data
      activeCountriesList.value = countriesList.value.filter((el: Country) => el.is_active)
    })
    .catch() // handled by interceptor
  }

  return { countriesList, activeCountriesList, loadCountriesList }
})

