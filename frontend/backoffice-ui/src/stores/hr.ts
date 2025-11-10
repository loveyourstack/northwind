import { ref } from 'vue'
import { defineStore } from 'pinia'
import { type Employee } from '@/types/hr'
import ax from '@/api'

export const useHRStore = defineStore('hr', () => {

  const employeesList = ref<Employee[]>([])
  const mandEmployeesList = ref<Employee[]>([])

  function loadEmployeesList() {
    var myURL = '/a/hr/employees?xfields=id,name&xsort=family_name,given_name&xper_page=5000'
    ax.get(myURL)
    .then(response => {
      employeesList.value = response.data.data
      mandEmployeesList.value = employeesList.value.filter((el: Employee) => el.id > 0)
    })
    .catch() // handled by interceptor
  }

  return { employeesList, mandEmployeesList, loadEmployeesList }
})

