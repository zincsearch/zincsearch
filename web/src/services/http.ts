import store from "../store";
import axios from "axios";
import { Notify } from "quasar";
import auth from "./auth";
import router from "../router";

function notifyAndLogout(error: any) {
  Notify.create({
    position: "bottom-right",
    progress: true,
    multiLine: true,
    color: "red-5",
    textColor: "white",
    icon: "warning",
    message: error.response.data["error"] || "Token expired",
  });
  store.dispatch("logout");
  router.replace({ name: "login" });
}

const http = () => {
  var instance = axios.create({
    // timeout: 10000,
    baseURL: store.state.API_ENDPOINT,
  });

  instance.interceptors.request.use(
    function (config) {
      config.withCredentials = true;
      return config;
    },
    function (error) {
      return Promise.reject(error);
    }
  );

  instance.interceptors.response.use(
    function (response) {
      return response;
    },
    async function (error) {
      if (error && error.response && error.response.status) {
        switch (error.response.status) {
          case 400:
            Notify.create({
              position: "bottom-right",
              progress: true,
              multiLine: true,
              color: "red-5",
              textColor: "white",
              icon: "warning",
              message: JSON.stringify(
                error.response.data["error"] || "Bad Request"
              ),
            });
            break;
          case 403:
            Notify.create({
              position: "bottom-right",
              progress: true,
              multiLine: true,
              color: "red-5",
              textColor: "white",
              icon: "warning",
              message: error.response.data["error"] || "No Permission",
            });
            break;
          case 404:
            Notify.create({
              position: "bottom-right",
              progress: true,
              multiLine: true,
              color: "red-5",
              textColor: "white",
              icon: "warning",
              message: error.response.data["error"] || "Not Found",
            });
            break;
          case 500:
            Notify.create({
              position: "bottom-right",
              progress: true,
              multiLine: true,
              color: "red-5",
              textColor: "white",
              icon: "warning",
              message: JSON.stringify(
                error.response.data["error"] || "Internal ServerError"
              ),
            });
            break;
          default:
          // noop
        }
      }
      if (error.config.url === "/api/login/refresh") {
        notifyAndLogout(error);
        return Promise.reject(error);
      } else if (error.config.url !== "/api/login/verify") {
        const config = error.config;
        if (error.response.status === 401 && !config?.retried) {
          config.retried = true;
          const response = await auth.refresh();
          if (response.status == 200) {
            return instance(error.config);
          }
        }
        return Promise.reject(config);
      } else {
        const response = await auth.refresh();
        if (response.status == 200) {
          return instance(error.config);
        }
        notifyAndLogout(error);
        return Promise.reject(error);
      }
    }
  );

  return instance;
};

export default http;
