import { act, fireEvent, render, screen } from '@testing-library/react';
import { MemoryRouter } from 'react-router';
import { afterEach, beforeEach, expect, it, vi } from 'vitest';
import Login from '../views/Login';
import Account from '../components/Account';
import auth from '../services/auth';
import { setCredentials } from '../auth';
import { changeLanguage, languages, translate, type Locale } from '../locales';
import en from '../locales/en';

vi.mock('../services/auth', () => ({ default: { login: vi.fn(), updateAccount: vi.fn() } }));
const dictionaries = import.meta.glob('../locales/*.ts', { eager: true, import: 'default' });
beforeEach(() => { vi.resetAllMocks(); setCredentials(null); changeLanguage('zh-cn'); });
afterEach(() => { setCredentials(null); changeLanguage('en'); });

it.each(['invalid', '401', 'network', 'server'])('localizes login %s errors and updates them when language changes', async kind => {
  if (kind === 'invalid') vi.mocked(auth.login).mockResolvedValue({ data: { validated: false } } as Awaited<ReturnType<typeof auth.login>>);
  else vi.mocked(auth.login).mockRejectedValue(kind === 'network' ? new Error('Network Error') : { response: { status: kind === '401' ? 401 : 500 } });
  render(<MemoryRouter><Login /></MemoryRouter>);
  fireEvent.change(screen.getByLabelText('用户名'), { target: { value: 'admin' } });
  fireEvent.change(screen.getByLabelText('密码'), { target: { value: 'wrong' } });
  fireEvent.click(screen.getByRole('button', { name: '登录' }));
  const invalid = kind === 'invalid' || kind === '401';
  expect(await screen.findByRole('alert')).toHaveTextContent(invalid ? '用户名或密码错误。' : '登录失败，请重试。');
  act(() => changeLanguage('en'));
  expect(screen.getByRole('alert')).toHaveTextContent(invalid ? 'Invalid user ID or password.' : 'Unable to sign in. Please try again.');
});

it('localizes empty login fields without sending a request', () => {
  render(<MemoryRouter><Login /></MemoryRouter>);
  fireEvent.click(screen.getByRole('button', { name: '登录' }));
  expect(screen.getByRole('alert')).toHaveTextContent('请输入用户名和密码。');
  expect(auth.login).not.toHaveBeenCalled();
});

it('localizes missing and incorrect current passwords, including live language changes', async () => {
  const user = { _id: 'admin', name: 'Admin', role: 'admin', base64encoded: 'test' };
  setCredentials(user);
  vi.mocked(auth.updateAccount).mockRejectedValue({ response: { status: 401 } });
  render(<Account user={user} onClose={vi.fn()} />);
  fireEvent.change(screen.getByLabelText('显示名称'), { target: { value: 'Root' } });
  fireEvent.click(screen.getByRole('button', { name: '保存修改' }));
  expect(screen.getByRole('alert')).toHaveTextContent('请输入当前密码。');
  expect(auth.updateAccount).not.toHaveBeenCalled();
  fireEvent.change(screen.getByLabelText('当前密码'), { target: { value: 'wrong' } });
  fireEvent.click(screen.getByRole('button', { name: '保存修改' }));
  expect(await screen.findByText(translate('account.wrongPassword'))).toHaveAttribute('role', 'alert');
  act(() => changeLanguage('en'));
  expect(screen.getByRole('alert')).toHaveTextContent('Current password is incorrect.');
});

it.each(Object.keys(languages) as Locale[])('provides authentication messages in %s without fallback', locale => {
  changeLanguage(locale);
  const dictionary = dictionaries[`../locales/${locale}.ts`] as typeof en;
  for (const key of Object.keys(en.authErrors) as (keyof typeof en.authErrors)[]) {
    expect(dictionary.authErrors[key]).toBeTruthy();
    expect(translate(`authErrors.${key}`)).toBe(dictionary.authErrors[key]);
    if (locale !== 'en') expect(dictionary.authErrors[key]).not.toBe(en.authErrors[key]);
  }
});
