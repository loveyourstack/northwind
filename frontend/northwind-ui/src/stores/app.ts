import { defineStore } from 'pinia'
import { ApiError } from '@/types/system'

export const useAppStore = defineStore('app', () => {
  const apiErr = null as ApiError | null
  const projectTitle = 'Northwind'

  return { apiErr, projectTitle }
})

