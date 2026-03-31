import { createApp } from 'vue'
import { createPinia } from 'pinia'
import { createI18n } from 'vue-i18n'

import App from './App.vue'
import router from './router'
import en from './locales/en.json'
import ru from './locales/ru.json'
import './styles/main.scss'

const browserLocale = typeof navigator !== 'undefined' ? navigator.language.toLowerCase() : 'en'
const initialLocale = browserLocale.startsWith('ru') ? 'ru' : 'en'

const i18n = createI18n({
  legacy: false,
  locale: initialLocale,
  fallbackLocale: 'en',
  messages: {
    en,
    ru,
  },
})

const app = createApp(App)

app.use(createPinia())
app.use(router)
app.use(i18n)
app.mount('#app')
