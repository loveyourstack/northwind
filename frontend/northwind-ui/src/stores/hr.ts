import { ref } from 'vue'
import { defineStore } from 'pinia'
import { type Employee } from '@/types/hr'
import ax from '@/api'

export const useHRStore = defineStore('hr', () => {

  const employeesList = ref<Employee[]>([])

  function loadEmployeesList() {
    var myURL = '/a/hr/employees?xfields=id,name&xsort=last_name,first_name&xper_page=5000'
    ax.get(myURL)
    .then(response => {
      employeesList.value = response.data.data
    })
    .catch() // handled by interceptor
  }

  return { employeesList, loadEmployeesList }
})

