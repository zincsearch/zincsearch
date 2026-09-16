import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { beforeEach, expect, it, vi } from 'vitest';
import PasswordInput from './PasswordInput';
import { changeLanguage, translate } from '../locales';

beforeEach(() => changeLanguage('en'));

it('toggles with the keyboard without changing the value or submitting', async () => {
  const submit = vi.fn(event => event.preventDefault());
  render(<form onSubmit={submit}><label htmlFor="password">Password</label><PasswordInput id="password" defaultValue="example1" /></form>);
  const user = userEvent.setup();
  const input = screen.getByLabelText('Password');
  const toggle = screen.getByRole('button', { name: 'Show password' });
  expect(input).toHaveAttribute('type', 'password');
  expect(toggle).toHaveAttribute('aria-controls', input.id);
  await user.click(screen.getByText('Password'));
  expect(input).toHaveFocus();
  await user.tab();
  expect(toggle).toHaveFocus();
  await user.keyboard('{Enter}');
  expect(input).toHaveAttribute('type', 'text');
  expect(screen.getByRole('button', { name: 'Hide password' })).toHaveAttribute('title', 'Hide password');
  await user.keyboard(' ');
  expect(input).toHaveAttribute('type', 'password');
  expect(input).toHaveValue('example1');
  expect(submit).not.toHaveBeenCalled();
});

it('disables both the input and toggle', async () => {
  render(<PasswordInput aria-label="Password" disabled />);
  const input = screen.getByLabelText('Password');
  const toggle = screen.getByRole('button', { name: 'Show password' });
  expect(input).toBeDisabled();
  expect(toggle).toBeDisabled();
  await userEvent.setup().click(toggle);
  expect(input).toHaveAttribute('type', 'password');
});

it('generates distinct control IDs and localizes both toggle labels', async () => {
  changeLanguage('zh-cn');
  render(<><PasswordInput aria-label="First" /><PasswordInput aria-label="Second" /></>);
  const first = screen.getByLabelText('First');
  const second = screen.getByLabelText('Second');
  expect(first.id).not.toBe(second.id);
  const toggles = screen.getAllByRole('button', { name: translate('passwordInput.show') });
  expect(toggles[0]).toHaveAttribute('aria-controls', first.id);
  expect(toggles[1]).toHaveAttribute('aria-controls', second.id);
  await userEvent.setup().click(toggles[0]);
  expect(screen.getByRole('button', { name: translate('passwordInput.hide') })).toHaveAttribute('aria-controls', first.id);
  expect(second).toHaveAttribute('type', 'password');
});
