import { createApp } from "vue";
import App from "./App.vue";
// import "./registerServiceWorker";
import router from "./router";
import store from "./store";
import i18n from "./locales";
import { Quasar } from "quasar";
import quasarUserOptions from "./quasar-user-options";

const app = createApp(App);

app.use(Quasar, quasarUserOptions).use(i18n).use(store).use(router).mount("#app");
