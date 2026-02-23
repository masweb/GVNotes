import { createI18n } from 'vue-i18n'
import en from '@/locales/en.json'
import es from '@/locales/es.json'

let savedLang = localStorage.getItem('lang')
if (savedLang == null) {
  savedLang = 'es'
  localStorage.setItem('lang', savedLang)
}

export const i18n = createI18n({
 globalInjection: true,
 locale: savedLang,
 fallbackLocale: 'en',
 messages: {
  en,
  es
 }
})
