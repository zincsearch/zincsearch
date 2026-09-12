import axios from 'axios';
import { getCredentials, setCredentials } from '../auth';

export const apiEndpoint = import.meta.env.VITE_API_ENDPOINT || import.meta.env.BASE_URL.replace(/\/ui\/$/, '').replace(/\/$/, '');
const http = () => {
  const credentials = getCredentials();
  const instance = axios.create({
    baseURL: apiEndpoint,
    headers: credentials ? { Authorization: `Basic ${credentials.base64encoded}` } : {},
  });
  instance.interceptors.response.use(response => response, error => {
    if (error.response?.status === 401) setCredentials(null);
    return Promise.reject(error);
  });
  return instance;
};
export default http;
