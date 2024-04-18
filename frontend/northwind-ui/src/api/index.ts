import axios, { AxiosInstance } from 'axios'
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

    appStore.apiErr = { 
      method: error.config.method.toUpperCase(), 
      url: error.config.url, 
      errMsg: error.response.status + ' - ' + error.response.data.err_description 
    }

    return Promise.reject(error)
  }
)

export default ax