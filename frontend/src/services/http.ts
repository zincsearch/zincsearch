import axios from 'axios';
import { getCredentials, setCredentials } from '../auth';
import { apiEndpoint } from './endpoint';

export { apiEndpoint };
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
