import { defineStore } from 'pinia'
import { ApiError } from '@/types/system'

export const useAppStore = defineStore('app', () => {
  const apiErr = null as ApiError | null

  const booleanOptions = [
    {value: true, title: 'Yes'},
    {value: false, title: 'No'}
  ]

  const projectTitle = 'Northwind'

  return { apiErr, booleanOptions, projectTitle }
})

