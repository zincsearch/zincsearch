import { act, render, screen, within } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { MemoryRouter } from 'react-router';
import { AxiosError, AxiosHeaders } from 'axios';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { AppRoutes } from '../App';
import { getCredentials, setCredentials } from '../auth';
import auth from '../services/auth';
import { changeLanguage } from '../locales';
import { validatePassword } from '../utils/password';

vi.mock('../services/auth', () => ({ default: { login: vi.fn(), updateAccount: vi.fn() } }));
vi.mock('../views/Search', () => ({ default: () => <h1>Search page</h1> }));
const ok = { data: { _id: 'admin', name: 'Admin', role: 'admin' }, status: 200, statusText: 'OK', headers: new AxiosHeaders(), config: { headers: new AxiosHeaders() } };
beforeEach(() => {
  setCredentials({ _id: 'admin', name: 'Admin', role: 'admin', base64encoded: 'YWRtaW46b2xk' });
  changeLanguage('en');
  vi.mocked(auth.updateAccount).mockReset();
});
const dialog = () => within(screen.getByRole('dialog'));
async function openDialog(via: 'avatar' | 'name' = 'avatar') {
  render(<MemoryRouter initialEntries={['/search']}><AppRoutes /></MemoryRouter>);
  const user = userEvent.setup();
  await user.click(via === 'avatar' ? screen.getByRole('button', { name: 'Account' }) : screen.getByRole('button', { name: 'Admin' }));
  expect(screen.getByRole('dialog', { name: 'My Account' })).toBeInTheDocument();
  expect(dialog().getByLabelText('Display name')).toHaveValue('Admin');
  return user;
}
async function fill(user: ReturnType<typeof userEvent.setup>, { name, current, next = '', confirm = '' }: { name?: string; current: string; next?: string; confirm?: string }) {
  if (name !== undefined) {
    await user.clear(dialog().getByLabelText('Display name'));
    if (name) await user.type(dialog().getByLabelText('Display name'), name);
  }
  if (current) await user.type(dialog().getByLabelText('Current password'), current);
  if (next) await user.type(dialog().getByLabelText('New password'), next);
  if (confirm) await user.type(dialog().getByLabelText('Confirm new password'), confirm);
  await user.click(dialog().getByRole('button', { name: 'Save Changes' }));
}
describe('Account dialog', () => {
  it('toggles each password independently without saving and resets on reopening', async () => {
    const user = await openDialog();
    const fields = ['Current password', 'New password', 'Confirm new password'].map(label => dialog().getByLabelText(label));
    for (const field of fields) {
      expect(field).toHaveAttribute('type', 'password');
      await user.type(field, 'example1');
    }
    for (const field of fields) {
      const toggle = dialog().getAllByRole('button', { name: 'Show password' }).find(button => button.getAttribute('aria-controls') === field.id)!;
      expect(toggle.closest('label')).toBeNull();
      expect(field.closest('label')).toBeNull();
      await user.click(toggle);
      expect(field).toHaveAttribute('type', 'text');
      expect(field).toHaveValue('example1');
      for (const other of fields.filter(input => input !== field)) expect(other).toHaveAttribute('type', 'password');
      await user.click(dialog().getByRole('button', { name: 'Hide password' }));
      expect(field).toHaveAttribute('type', 'password');
    }
    expect(auth.updateAccount).not.toHaveBeenCalled();
    await user.click(dialog().getAllByRole('button', { name: 'Show password' })[0]);
    await user.click(dialog().getByRole('button', { name: 'Close' }));
    await user.click(screen.getByRole('button', { name: 'Account' }));
    expect(dialog().getByLabelText('Current password')).toHaveAttribute('type', 'password');
    expect(dialog().getByLabelText('Current password')).toHaveValue('');
  });
  it('opens from the avatar icon and the header name', async () => {
    const user = await openDialog('avatar');
    await user.click(dialog().getByRole('button', { name: 'Close' }));
    expect(screen.queryByRole('dialog')).not.toBeInTheDocument();
    await user.click(screen.getByRole('button', { name: 'Admin' }));
    expect(screen.getByRole('dialog', { name: 'My Account' })).toBeInTheDocument();
  });
  it('validates locally before calling the API', async () => {
    const user = await openDialog();
    await fill(user, { current: 'oldpass1' });
    expect(await dialog().findByRole('alert')).toHaveTextContent('Nothing to change');
    await fill(user, { name: 'ab', current: '' });
    expect(dialog().getByRole('alert')).toHaveTextContent('at least 3 characters');
    await fill(user, { name: 'Admin', current: '', next: 'newpass1', confirm: 'other' });
    expect(dialog().getByRole('alert')).toHaveTextContent('should match');
    expect(auth.updateAccount).not.toHaveBeenCalled();
    expect(validatePassword('short1', 'short1')).toContain('at least 8');
    expect(validatePassword('12345678', '12345678')).toContain('letter');
    expect(validatePassword('abcdefgh', 'abcdefgh')).toContain('digit');
    expect(validatePassword('newpass1', 'newpass1')).toBe('');
  });
  it('updates the name-validation error when the language changes', async () => {
    const user = await openDialog();
    await fill(user, { name: 'ab', current: 'oldpass1' });
    expect(dialog().getByRole('alert')).toHaveTextContent('at least 3 characters');
    act(() => changeLanguage('zh-cn'));
    expect(dialog().getByRole('alert')).toHaveTextContent('3');
    expect(dialog().getByRole('alert')).not.toHaveTextContent('at least 3 characters');
    expect(auth.updateAccount).not.toHaveBeenCalled();
  });
  it('renames and refreshes stale stored credentials with the verified current password', async () => {
    vi.mocked(auth.updateAccount).mockResolvedValue(ok);
    const user = await openDialog();
    await fill(user, { name: 'Root', current: 'oldpass1' });
    expect(await dialog().findByRole('status')).toHaveTextContent('Account updated');
    expect(auth.updateAccount).toHaveBeenCalledWith({ _id: 'admin', password: 'oldpass1', name: 'Root', new_password: undefined });
    expect(getCredentials()).toMatchObject({ name: 'Root', base64encoded: btoa('admin:oldpass1') });
    expect(screen.getByRole('button', { name: 'Root' })).toBeInTheDocument();
  });
  it('re-encodes Basic auth credentials when the password changes', async () => {
    vi.mocked(auth.updateAccount).mockResolvedValue(ok);
    const user = await openDialog();
    await fill(user, { current: 'oldpass1', next: 'newpass1', confirm: 'newpass1' });
    expect(await dialog().findByRole('status')).toHaveTextContent('Account updated');
    expect(auth.updateAccount).toHaveBeenCalledWith({ _id: 'admin', password: 'oldpass1', name: 'Admin', new_password: 'newpass1' });
    expect(getCredentials()?.base64encoded).toBe(btoa('admin:newpass1'));
  });
  it.each(['logout', 'new session', 'same credentials'])('does not overwrite a %s while saving', async (change) => {
    let finish!: (value: typeof ok) => void;
    vi.mocked(auth.updateAccount).mockReturnValue(new Promise(resolve => { finish = resolve; }));
    const user = await openDialog();
    await fill(user, { name: 'Root', current: 'oldpass1' });
    const replacement = change === 'logout' ? null : { _id: 'admin', name: 'Admin', role: 'admin', base64encoded: change === 'same credentials' ? 'YWRtaW46b2xk' : 'new-session' };
    await act(async () => {
      setCredentials(null);
      setCredentials(replacement);
      finish(ok);
    });
    expect(getCredentials()).toBe(replacement);
  });
  it('reports a wrong current password and keeps the session', async () => {
    vi.mocked(auth.updateAccount).mockRejectedValue(new AxiosError('Unauthorized', 'ERR_BAD_REQUEST', undefined, null, { status: 401, statusText: 'Unauthorized', data: {}, headers: new AxiosHeaders(), config: { headers: new AxiosHeaders() } }));
    const user = await openDialog();
    await fill(user, { name: 'Root', current: 'wrong' });
    expect(await dialog().findByRole('alert')).toHaveTextContent('Current password is incorrect');
    expect(getCredentials()).toMatchObject({ name: 'Admin', base64encoded: 'YWRtaW46b2xk' });
    expect(dialog().getByRole('button', { name: 'Save Changes' })).toBeEnabled();
  });
});
