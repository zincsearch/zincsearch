import axios from 'axios';
import { getCredentials, setCredentials } from '../auth';
import { apiEndpoint } from './endpoint';

export { apiEndpoint };
const http = ({ authenticated = true } = {}) => {
  const credentials = authenticated ? getCredentials() : null;
  const instance = axios.create({
    baseURL: apiEndpoint,
    headers: credentials ? { Authorization: `Basic ${credentials.base64encoded}` } : {},
  });
  instance.interceptors.response.use(response => response, error => {
    const current = getCredentials();
    if (error.response?.status === 401 && credentials && current === credentials) {
      setCredentials(null);
    }
    return Promise.reject(error);
  });
  return instance;
};
export default http;
