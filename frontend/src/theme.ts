import { useSyncExternalStore } from 'react';

export type Theme = 'system' | 'light' | 'dark';
const storageKey = 'zincsearch-theme';
function readTheme(): Theme {
  try {
    const value = localStorage.getItem(storageKey);
    return value === 'light' || value === 'dark' ? value : 'system';
  } catch {
    return 'system';
  }
}
let theme = readTheme();
const listeners = new Set<() => void>();
function applyTheme() {
  document.documentElement.dataset.theme = theme;
  listeners.forEach(listener => listener());
}
applyTheme();
export function setTheme(value: Theme) {
  theme = value;
  try {
    localStorage.setItem(storageKey, value);
  } catch {
    // Keep the setting usable when browser storage is unavailable.
  }
  applyTheme();
}
window.addEventListener('storage', event => {
  if (event.key === storageKey || event.key === null) {
    theme = readTheme();
    applyTheme();
  }
});
function subscribe(listener: () => void) {
  listeners.add(listener);
  return () => listeners.delete(listener);
}
export function useTheme() {
  return useSyncExternalStore(subscribe, () => theme);
}
