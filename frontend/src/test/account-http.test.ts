import axios, { AxiosError, AxiosHeaders, type InternalAxiosRequestConfig } from 'axios';
import { afterEach, beforeEach, describe, expect, it } from 'vitest';
import { getCredentials, setCredentials } from '../auth';
import auth from '../services/auth';
import http from '../services/http';

const adapter = axios.defaults.adapter;
const user = { _id: 'admin', name: 'Admin', role: 'admin', base64encoded: 'old-credentials' };
function unauthorized(config: InternalAxiosRequestConfig) {
  return new AxiosError('Unauthorized', 'ERR_BAD_REQUEST', config, null, {
    status: 401, statusText: 'Unauthorized', data: { error: 'invalid credentials' }, headers: new AxiosHeaders(), config,
  });
}
beforeEach(() => setCredentials(user));
afterEach(() => { axios.defaults.adapter = adapter; setCredentials(null); });

describe('Account HTTP authentication', () => {
  it('keeps the session on a wrong current password and allows retry', async () => {
    axios.defaults.adapter = async config => { throw unauthorized(config); };
    await expect(auth.updateAccount({ _id: user._id, password: 'wrong', name: 'Root' })).rejects.toThrow('Unauthorized');
    expect(getCredentials()).toEqual(user);
    expect(JSON.parse(localStorage.getItem('creds')!)).toEqual(user);

    axios.defaults.adapter = async config => {
      expect(config.headers.has('Authorization')).toBe(false);
      expect(JSON.parse(config.data)).toEqual({ _id: user._id, password: 'correct', name: 'Root' });
      return { status: 200, statusText: 'OK', data: { _id: user._id, name: 'Root', role: 'admin' }, headers: new AxiosHeaders(), config };
    };
    await expect(auth.updateAccount({ _id: user._id, password: 'correct', name: 'Root' })).resolves.toMatchObject({ status: 200 });
  });

  it('does not log out new credentials when an old request returns 401', async () => {
    const client = http();
    const updated = { ...user, base64encoded: 'new-credentials' };
    client.defaults.adapter = async config => {
      expect(config.headers.get('Authorization')).toBe('Basic old-credentials');
      setCredentials(updated);
      throw unauthorized(config);
    };
    await expect(client.get('/api/index')).rejects.toThrow('Unauthorized');
    expect(getCredentials()).toEqual(updated);
  });

  it('keeps a new session with identical credentials when an old request returns 401', async () => {
    const client = http();
    const replacement = { ...user };
    client.defaults.adapter = async config => {
      setCredentials(null);
      setCredentials(replacement);
      throw unauthorized(config);
    };
    await expect(client.get('/api/index')).rejects.toThrow('Unauthorized');
    expect(getCredentials()).toBe(replacement);
    expect(JSON.parse(localStorage.getItem('creds')!)).toEqual(replacement);
  });

  it('still clears rejected current credentials on protected endpoints', async () => {
    const client = http();
    client.defaults.adapter = async config => { throw unauthorized(config); };
    await expect(client.get('/api/index')).rejects.toThrow('Unauthorized');
    expect(getCredentials()).toBeNull();
  });
});
