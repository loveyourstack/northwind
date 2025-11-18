import ax from '@/api'

export function callDelete(deleteUrl: string, successFn?: Function) {

  if (!confirm('Are you sure?')) {
    return
  }

  ax.delete(deleteUrl)
    .then(() => {
      if (successFn) {
        successFn()
      }
    })
    .catch() // handled by interceptor
}