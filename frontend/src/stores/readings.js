import { defineStore } from 'pinia'

export const useReadingsStore = defineStore('readings', {
  state: () => ({
    cardOfDay: null,
    spreads: [],
    history: [],
  }),
})
