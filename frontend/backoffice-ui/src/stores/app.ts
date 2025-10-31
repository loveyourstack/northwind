import { defineStore } from 'pinia'
import { type ApiError } from '@/types/system'

export const useAppStore = defineStore('app', () => {
  const apiErr = null as ApiError | null

  const booleanOptions = [
    {value: true, title: 'Yes'},
    {value: false, title: 'No'}
  ]

  const company = 'Northwind'
  const projectTitle = 'Backoffice'

  return { apiErr, booleanOptions, company, projectTitle }
})

