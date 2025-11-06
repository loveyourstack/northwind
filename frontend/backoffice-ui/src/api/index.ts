import axios, { type AxiosInstance } from 'axios'
import { useAppStore } from '@/stores/app'

const appStore = useAppStore()

const ax: AxiosInstance = axios.create({
  baseURL: import.meta.env.VITE_API_URL,
  headers: {
    'Content-type': 'application/json',
  },
})

ax.interceptors.response.use(
  (response) => { return response },
  (error) => {
    //console.log(error.toJSON())

    // special handling for blob requests (file downloads: see file.ts) errors in order to get err_description
    // from https://github.com/axios/axios/issues/815
    if (error.request.responseType === 'blob') {
      return new Promise((resolve, reject) => {
          let reader = new FileReader();
          reader.onload = () => {
              error.response.data = JSON.parse(String(reader.result))

              appStore.apiErr = { 
                method: error.config.method.toUpperCase(), 
                url: error.config.url, 
                errMsg: error.response.status + ' - ' + error.response.data.err_description 
              }

              resolve(Promise.reject(error))
          }

          reader.onerror = () => {
              reject(error)
          }
          reader.readAsText(error.response.data)
      })
    }

    // regular requests

    appStore.apiErr = { 
      method: error.config.method.toUpperCase(), 
      url: error.config.url, 
      errMsg: error.response.status + ' - ' + error.response.data.err_description 
    }

    return Promise.reject(error)
  }
)

export default ax