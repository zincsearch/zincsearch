import { useSyncExternalStore } from 'react';

export interface Credentials {
  _id: string;
  name: string;
  role: string;
  base64encoded: string;
}

export function readCredentials(): Credentials | null {
  try {
    const value = JSON.parse(localStorage.getItem('creds') || 'null');
    return value && typeof value._id === 'string' && value._id &&
      typeof value.base64encoded === 'string' && value.base64encoded ? value : null;
  } catch {
    return null;
  }
}

let credentials = readCredentials();
const listeners = new Set<() => void>();
export function setCredentials(value: Credentials | null) {
  credentials = value;
  if (value) localStorage.setItem('creds', JSON.stringify(value));
  else localStorage.removeItem('creds');
  listeners.forEach(listener => listener());
}
export const getCredentials = () => credentials;
export function useAuth() {
  return useSyncExternalStore(listener => {
    listeners.add(listener);
    return () => listeners.delete(listener);
  }, getCredentials);
}
window.addEventListener('storage', event => {
  if (event.key === 'creds' || event.key === null) {
    credentials = readCredentials();
    listeners.forEach(listener => listener());
  }
});

export function encodeCredentials(id: string, password: string) {
  return btoa(Array.from(new TextEncoder().encode(`${id}:${password}`), byte => String.fromCharCode(byte)).join(''));
}
