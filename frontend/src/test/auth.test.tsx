import { act, render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { MemoryRouter, Route, Routes } from 'react-router';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import Login from '../views/Login';
import { encodeCredentials, getCredentials, readCredentials, setCredentials } from '../auth';
import auth from '../services/auth';
import http from '../services/http';
import { AxiosError, AxiosHeaders } from 'axios';
import { changeLanguage, translate } from '../locales';

vi.mock('../services/auth', () => ({ default: { login: vi.fn() } }));
beforeEach(() => { setCredentials(null); changeLanguage('en'); });
function login() {
  render(<MemoryRouter initialEntries={['/login']}><Routes><Route path="/login" element={<Login />} /><Route path="/search" element={<h1>Search page</h1>} /></Routes></MemoryRouter>);
}
describe('Authentication', () => {
  it('accepts legacy credentials and safely handles malformed storage', () => {
    localStorage.setItem('creds', '{broken');
    expect(readCredentials()).toBeNull();
    localStorage.setItem('creds', JSON.stringify({ _id: 'admin', base64encoded: 'YQ==' }));
    expect(readCredentials()?._id).toBe('admin');
    expect(encodeCredentials('用户', 'secret')).toBe('55So5oi3OnNlY3JldA==');
  });
  it('logs in and persists credentials without the plaintext password', async () => {
    vi.mocked(auth.login).mockResolvedValue({ data: { validated: true, user: { name: 'Admin', role: 'admin' } }, status: 200, statusText: 'OK', headers: new AxiosHeaders(), config: { headers: new AxiosHeaders() } });
    login();
    const user = userEvent.setup();
    await user.type(screen.getByLabelText('User ID'), 'admin');
    await user.type(screen.getByLabelText('Password'), 'secret');
    await user.click(screen.getByRole('button', { name: 'Sign In' }));
    expect(await screen.findByText('Search page')).toBeInTheDocument();
    expect(auth.login).toHaveBeenCalledWith({ _id: 'admin', password: 'secret', base64encoded: 'YWRtaW46c2VjcmV0' });
    expect(JSON.parse(localStorage.getItem('creds')!)).not.toHaveProperty('password');
  });
  it('allows retry after failed network requests', async () => {
    vi.mocked(auth.login).mockRejectedValue(new Error('Network Error'));
    login();
    const user = userEvent.setup();
    await user.type(screen.getByLabelText('User ID'), 'admin');
    await user.type(screen.getByLabelText('Password'), 'secret');
    await user.click(screen.getByRole('button', { name: 'Sign In' }));
    expect(await screen.findByRole('alert')).toHaveTextContent('Network Error');
    expect(screen.getByRole('button', { name: 'Sign In' })).toBeEnabled();
    expect(getCredentials()).toBeNull();
  });
  it('attaches Basic auth and clears the session on unauthorized responses', async () => {
    setCredentials({ _id: 'admin', name: 'Admin', role: 'admin', base64encoded: 'YQ==' });
    const client = http();
    expect(client.defaults.headers.Authorization).toBe('Basic YQ==');
    client.defaults.adapter = async config => {
      throw new AxiosError('Unauthorized', 'ERR_BAD_REQUEST', config, null, { status: 401, statusText: 'Unauthorized', data: {}, headers: new AxiosHeaders(), config });
    };
    await expect(client.get('/api/user')).rejects.toThrow('Unauthorized');
    expect(getCredentials()).toBeNull();
    expect(localStorage.getItem('creds')).toBeNull();
  });
  it('synchronizes logout from another tab', async () => {
    setCredentials({ _id: 'admin', name: 'Admin', role: 'admin', base64encoded: 'YQ==' });
    localStorage.removeItem('creds');
    act(() => window.dispatchEvent(new StorageEvent('storage', { key: 'creds' })));
    await waitFor(() => expect(getCredentials()).toBeNull());
  });
});
it('switches all original locales with English fallback', () => {
  changeLanguage('zh-cn');
  expect(translate('menu.search')).not.toBe('Search');
  changeLanguage('en');
  expect(translate('menu.search')).toBe('Search');
  expect(translate('unknown.key')).toBe('unknown.key');
});
