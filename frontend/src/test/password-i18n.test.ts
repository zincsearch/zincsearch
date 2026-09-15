import { afterEach, describe, expect, it } from 'vitest';
import { changeLanguage, languages, translate, type Locale } from '../locales';
import en from '../locales/en';
import { validatePassword } from '../utils/password';

const dictionaries = import.meta.glob('../locales/*.ts', { eager: true, import: 'default' });
const cases = [
  ['short1', 'short1', 'tooShort'],
  ['12345678', '12345678', 'missingLetter'],
  ['abcdefgh', 'abcdefgh', 'missingDigit'],
  ['newpass1', 'newpass2', 'mismatch'],
] as const;

afterEach(() => changeLanguage('en'));
describe('Localized password validation', () => {
  it.each(Object.keys(languages) as Locale[])('translates every password error in %s without fallback', locale => {
    changeLanguage(locale);
    const dictionary = dictionaries[`../locales/${locale}.ts`] as typeof en;
    for (const [password, confirm, key] of cases) {
      const message = dictionary.passwordValidation[key];
      expect(message).toBeTruthy();
      expect(validatePassword(password, confirm)).toBe(message);
      if (locale !== 'en') expect(message).not.toBe(en.passwordValidation[key]);
    }
    expect(validatePassword('newpass1', 'newpass1')).toBe('');
    expect(translate('account.wrongPassword')).toBe(dictionary.account.wrongPassword);
  });

  it('uses the newly selected language on subsequent validation', () => {
    changeLanguage('en');
    expect(validatePassword('short1', 'short1')).toContain('at least 8');
    changeLanguage('zh-cn');
    expect(validatePassword('short1', 'short1')).toBe('密码至少需要 8 个字符。');
    expect(validatePassword('newpass1', 'newpass2')).toBe('密码与确认密码必须一致。');
  });
});
