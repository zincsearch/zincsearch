import { act, render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { afterEach, expect, it, vi } from 'vitest';
import ThemeSelect from './ThemeSelect';
import { setTheme } from '../theme';

afterEach(() => {
  vi.restoreAllMocks();
  act(() => setTheme('system'));
  localStorage.clear();
});

it('defaults to system and persists each explicit setting', async () => {
  const user = userEvent.setup();
  render(<ThemeSelect />);
  const select = screen.getByRole('combobox', { name: 'Theme' });
  expect(select).toHaveTextContent('System');
  for (const [label, value] of [['Dark', 'dark'], ['Light', 'light'], ['System', 'system']]) {
    await user.click(select);
    await user.click(screen.getByRole('option', { name: label }));
    expect(select).toHaveTextContent(label);
    expect(localStorage.getItem('zincsearch-theme')).toBe(value);
    expect(document.documentElement.dataset.theme).toBe(value);
  }
});

it('synchronizes saved settings, invalid values and clearing storage across tabs', () => {
  render(<ThemeSelect />);
  for (const [value, label] of [['dark', 'Dark'], ['light', 'Light'], ['invalid', 'System']]) {
    act(() => {
      localStorage.setItem('zincsearch-theme', value);
      window.dispatchEvent(new StorageEvent('storage', { key: 'zincsearch-theme' }));
    });
    expect(screen.getByRole('combobox')).toHaveTextContent(label);
  }
  act(() => {
    setTheme('dark');
    localStorage.clear();
    window.dispatchEvent(new StorageEvent('storage', { key: null }));
  });
  expect(screen.getByRole('combobox')).toHaveTextContent('System');
});

it('still applies the theme if storage is blocked', () => {
  vi.spyOn(Storage.prototype, 'setItem').mockImplementation(() => { throw new DOMException('Blocked', 'SecurityError'); });
  render(<ThemeSelect />);
  act(() => setTheme('dark'));
  expect(screen.getByRole('combobox')).toHaveTextContent('Dark');
  expect(document.documentElement.dataset.theme).toBe('dark');
});
