import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { MemoryRouter } from 'react-router';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import Cookies from 'js-cookie';
import Keys from '../constant/key';
import { setCredentials } from '../auth';
import { AppRoutes } from '../App';
import { changeLanguage, getLocale, languages, translate, type Locale } from '../locales';
import en from '../locales/en';
import { getLanguage, setLanguage } from '../utils/cookies';

beforeEach(() => {
  setCredentials(null);
  changeLanguage('en');
  Cookies.remove(Keys.languageKey, { path: '/' });
  vi.spyOn(navigator, 'languages', 'get').mockReturnValue(['en-US']);
  vi.spyOn(navigator, 'language', 'get').mockReturnValue('en-US');
});
afterEach(() => {
  vi.restoreAllMocks();
  changeLanguage('en');
});

describe('locale matching', () => {
  it.each([
    ['es-MX', 'es'], ['FR-ca', 'fr'], ['de-AT', 'de'], ['pt-BR', 'pt'],
    ['pt-PT', 'pt'], ['ja-JP', 'ja'], ['ko-KR', 'ko'], ['ru-RU', 'ru'],
    ['tr-TR', 'tr'], ['en-GB', 'en'], ['zh', 'zh-cn'], ['zh-SG', 'zh-cn'],
    ['zh-Hans', 'zh-cn'], ['zh-Hant', 'zh-tw'], ['zh-HK', 'zh-tw'],
    ['zh-MO', 'zh-tw'], ['ZH_TW', 'zh-tw'], ['zh-Hans-HK', 'zh-cn'],
    ['zh-Hant-CN', 'zh-tw'], ['xx-XX', 'en'], ['english', 'en'],
  ])('matches browser preference %s to %s', (tag, expected) => {
    vi.spyOn(navigator, 'languages', 'get').mockReturnValue([tag]);
    expect(getLocale()).toBe(expected);
  });
  it('uses the first supported browser preference', () => {
    vi.spyOn(navigator, 'languages', 'get').mockReturnValue(['xx', 'fr-CA', 'de']);
    expect(getLocale()).toBe('fr');
  });
  it('uses navigator.language when the preference list is empty', () => {
    vi.spyOn(navigator, 'languages', 'get').mockReturnValue([]);
    vi.spyOn(navigator, 'language', 'get').mockReturnValue('ja-JP');
    expect(getLocale()).toBe('ja');
  });
  it('prefers a normalized saved choice and ignores unsupported cookies', () => {
    vi.spyOn(navigator, 'languages', 'get').mockReturnValue(['de-DE']);
    setLanguage('FR_ca');
    expect(getLocale()).toBe('fr');
    setLanguage('french');
    expect(getLocale()).toBe('de');
    setLanguage('constructor');
    expect(getLocale()).toBe('de');
  });
  it('sets the document language on initialization', async () => {
    setLanguage('ko-KR');
    vi.resetModules();
    await import('../locales');
    expect(document.documentElement.lang).toBe('ko');
  });
});

describe('translations', () => {
  const dictionaries = import.meta.glob<{ default: typeof en }>('../locales/{de,es,fr,pt,ja,ko,ru}.ts', { eager: true });
  it.each(Object.entries(dictionaries))('has every English key in %s', (_path, { default: dictionary }) => {
    expect(Object.keys(dictionary).sort()).toEqual(Object.keys(en).sort());
    for (const section of Object.keys(en) as (keyof typeof en)[]) {
      expect(Object.keys(dictionary[section]).sort()).toEqual(Object.keys(en[section]).sort());
      for (const value of Object.values(dictionary[section])) expect(value.trim()).not.toBe('');
    }
  });
  it.each(Object.keys(languages) as Locale[])('persists and translates %s', locale => {
    changeLanguage(locale);
    expect(getLanguage()).toBe(locale);
    expect(getLocale()).toBe(locale);
    expect(document.documentElement.lang).toBe(locale);
    expect(translate('menu.search')).not.toBe('menu.search');
    if (locale !== 'en') expect(translate('menu.search')).not.toBe(en.menu.search);
    expect(translate('missing.key')).toBe('missing.key');
  });
  it('retains English fallback and interpolation', () => {
    changeLanguage('zh-cn');
    expect(translate('user.CREATED')).toBe(en.user.CREATED);
    expect(translate('Found {count} {kind}', { count: 0, kind: 'items' })).toBe('Found 0 items');
    expect(translate('Found {count}')).toBe('Found {count}');
  });
});

it('switches language before login and in the console without a reload', async () => {
  const user = userEvent.setup();
  const login = render(<MemoryRouter initialEntries={['/login']}><AppRoutes /></MemoryRouter>);
  await user.click(screen.getByRole('combobox', { name: 'Language' }));
  expect(screen.getAllByRole('option')).toHaveLength(11);
  await user.click(screen.getByRole('option', { name: languages.es }));
  expect(screen.getByRole('button', { name: 'Iniciar sesión' })).toBeInTheDocument();
  expect(screen.getByRole('combobox', { name: 'Idioma' })).toHaveValue('es');
  login.unmount();
  setCredentials({ _id: 'admin', name: 'Admin', role: 'admin', base64encoded: 'YQ==' });
  render(<MemoryRouter initialEntries={['/missing']}><AppRoutes /></MemoryRouter>);
  expect(screen.getByRole('combobox', { name: 'Idioma' })).toHaveValue('es');
  await user.click(screen.getByRole('combobox', { name: 'Idioma' }));
  await user.click(screen.getByRole('option', { name: languages.ja }));
  expect(screen.getByRole('combobox', { name: '言語' })).toHaveValue('ja');
  expect(screen.getByRole('link', { name: '検索' })).toBeInTheDocument();
  expect(getLanguage()).toBe('ja');
});
