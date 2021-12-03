import axios from "axios";
import store from "./store";

// a wrapper around axios to pass credentials.
// standard interceptors was not working as the interceptor was asynchronous
var http = {
  get: async function (url) {
    const base64encoded = store.state.user.base64encoded;

    const instance = axios.create({
      // timeout: 10000,
      headers: { Authorization: "Basic " + base64encoded },
    });

    return instance.get(url);
  },
  post: async function (url, formData) {
    const base64encoded = store.state.user.base64encoded;

    const instance = axios.create({
      // timeout: 10000,
      headers: { Authorization: "Basic " + base64encoded },
    });

    return instance.post(url, formData);
  },
  put: async function (url, formData) {
    const base64encoded = store.state.user.base64encoded;

    const instance = axios.create({
      // timeout: 10000,
      headers: { Authorization: "Basic " + base64encoded },
    });

    return instance.put(url, formData);
  },
  delete: async function (url) {
    // var creds = await Auth.currentAuthenticatedUser()
    // const idToken = creds.signInUserSession.idToken.jwtToken

    const instance = axios.create({
      // timeout: 10000,
      //   headers: { 'Authorization': 'Bearer ' + idToken }
    });

    return instance.delete(url);
  },
  putPresignedFile: async function (url, file, config) {
    const instance = axios.create({});

    return instance.put(url, file, config);
  },
};

export default http;
