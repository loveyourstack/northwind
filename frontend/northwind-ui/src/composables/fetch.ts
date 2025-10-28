import { ref, type Ref, watchEffect } from 'vue'
import ax from '@/api'
import { type GetMetadata } from '@/types/system'

export function useFetch(url: string, result: Ref<any>, successFn?: Function, loading?: Ref<boolean>) {

  const fetchData = () => {

    if (loading) { 
      loading.value = true
    }

    ax.get(url)
      .then(resp => {
        result.value = resp.data.data

        if (successFn) {
          successFn()
        }
      })
      .catch() // handled by interceptor
      .finally(() => { 
        if (loading) { loading.value = false } 
      })
  }

  watchEffect(() => {
    fetchData()
  })
}

// useFetchDt gets data using the supplied url and changes the supplied items ref and other metadata refs
// adapted from https://vuejs.org/guide/reusability/composables.html#accepting-reactive-state
export function useFetchDt(url: string, items: Ref<any>, totalItems: Ref<number>, totalItemsIsEstimate: Ref<boolean>, totalItemsEstimated: Ref<number>, successFn?: Function) {

  const metadata = ref<GetMetadata | null>()
  
  const fetchData = () => {
  
    ax.get(url)
      .then(resp => {
        items.value = resp.data.data
        metadata.value = resp.data.metadata

        totalItemsIsEstimate.value = metadata.value!.total_count_is_estimated
        if (totalItemsIsEstimate.value) {
          totalItemsEstimated.value = metadata.value!.total_count
          totalItems.value = 101 // workaround so that next page button is enabled even if estimate is too low
        } else {
          totalItems.value = metadata.value!.total_count
        }

        if (successFn) {
          successFn()
        }
      })
      .catch() // handled by interceptor
  }

  watchEffect(() => {
    fetchData()
  })

  return
}
