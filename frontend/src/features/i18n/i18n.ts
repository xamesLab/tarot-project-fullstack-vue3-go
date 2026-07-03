import i18n from 'i18next';
import { initReactI18next } from 'react-i18next';
import commonEn from './locales/en/common.json';
import homeEn from './locales/en/home.json';
import commonRu from './locales/ru/common.json';
import homeRu from './locales/ru/home.json';

void i18n.use(initReactI18next).init({
  resources: {
    ru: { common: commonRu, home: homeRu },
    en: { common: commonEn, home: homeEn },
  },
  lng: 'ru',
  fallbackLng: 'ru',
  defaultNS: 'common',
  interpolation: { escapeValue: false },
});

export { i18n };
