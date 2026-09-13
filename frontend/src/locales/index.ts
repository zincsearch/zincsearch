import { useSyncExternalStore } from 'react';
import en from './en';
import tr from './tr';
import de from './de';
import es from './es';
import fr from './fr';
import pt from './pt';
import ja from './ja';
import ko from './ko';
import ru from './ru';
import zhCN from './zh-cn';
import zhTW from './zh-tw';
import { getLanguage, setLanguage } from '../utils/cookies';

export const languages = {
  en: 'English', de: 'Deutsch', es: 'Español', fr: 'Français',
  pt: 'Português', ja: '日本語', ko: '한국어', ru: 'Русский',
  tr: 'Türkçe', 'zh-cn': '简体中文', 'zh-tw': '繁體中文',
};
export type Locale = keyof typeof languages;
const messages = { en, de, es, fr, pt, ja, ko, ru, tr, 'zh-cn': zhCN, 'zh-tw': zhTW };
function matchLocale(value: string): Locale | undefined {
  const tag = value.trim().toLowerCase().replace(/_/g, '-');
  if (Object.hasOwn(languages, tag)) return tag as Locale;
  const [language, ...parts] = tag.split('-');
  if (language === 'zh') {
    if (parts.includes('hant')) return 'zh-tw';
    if (parts.includes('hans')) return 'zh-cn';
    return parts.some(part => ['tw', 'hk', 'mo'].includes(part)) ? 'zh-tw' : 'zh-cn';
  }
  if (Object.hasOwn(languages, language)) return language as Locale;
}
export function getLocale(): Locale {
  const saved = matchLocale(getLanguage() || '');
  if (saved) return saved;
  for (const language of [...(navigator.languages || []), navigator.language]) {
    const matched = matchLocale(language);
    if (matched) return matched;
  }
  return 'en';
}
let locale = getLocale();
document.documentElement.lang = locale;
const listeners = new Set<() => void>();
export function changeLanguage(value: Locale) {
  locale = value;
  setLanguage(value);
  document.documentElement.lang = value;
  listeners.forEach(listener => listener());
}
function lookup(dictionary: unknown, key: string): string | undefined {
  let value = dictionary;
  for (const part of key.split('.')) {
    if (!value || typeof value !== 'object') return;
    value = (value as Record<string, unknown>)[part];
  }
  return typeof value === 'string' ? value : undefined;
}
export function translate(key: string, values?: Record<string, string | number>) {
  const text = lookup(messages[locale], key) || lookup(en, key) || key;
  return text.replace(/\{(\w+)\}/g, (match, name: string) => String(values?.[name] ?? match));
}
export function useTranslation() {
  const current = useSyncExternalStore(listener => {
    listeners.add(listener);
    return () => listeners.delete(listener);
  }, () => locale);
  return { t: translate, locale: current, changeLanguage };
}
