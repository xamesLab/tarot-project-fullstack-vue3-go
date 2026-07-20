import i18n from 'i18next';
import { initReactI18next } from 'react-i18next';
import commonEn from './locales/en/common.json';
import homeEn from './locales/en/home.json';
import profileEn from './locales/en/profile.json';
import adminEn from './locales/en/admin.json';
import spreadsEn from './locales/en/spreads.json';
import celestialEn from './locales/en/celestial.json';
import commonRu from './locales/ru/common.json';
import homeRu from './locales/ru/home.json';
import profileRu from './locales/ru/profile.json';
import adminRu from './locales/ru/admin.json';
import spreadsRu from './locales/ru/spreads.json';
import celestialRu from './locales/ru/celestial.json';

void i18n.use(initReactI18next).init({
  resources: {
    ru: { common: commonRu, home: homeRu, profile: profileRu, admin: adminRu, spreads: spreadsRu, celestial: celestialRu },
    en: { common: commonEn, home: homeEn, profile: profileEn, admin: adminEn, spreads: spreadsEn, celestial: celestialEn },
  },
  lng: 'ru',
  fallbackLng: 'ru',
  defaultNS: 'common',
  interpolation: { escapeValue: false },
});

export { i18n };
