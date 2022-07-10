import { createApp } from "vue";
import { Quasar, Dialog, Notify } from "quasar";
import VueApexCharts from "vue3-apexcharts";

// Import icon libraries
import "@quasar/extras/roboto-font/roboto-font.css";
import "@quasar/extras/material-icons/material-icons.css";
// import "@quasar/extras/material-icons-outlined/material-icons-outlined.css";
// import "@quasar/extras/material-icons-round/material-icons-round.css";

// Import Quasar css
import "quasar/src/css/index.sass";
import i18n from "./locales";
import App from "./App.vue";
import router from "./router";
import store from "./store";

const myApp = createApp(App);

myApp
  .use(Quasar, {
    plugins: [Dialog, Notify], // import Quasar plugins and add here
  })
  .use(i18n);

myApp.use(VueApexCharts);

// createApp(App).mount('#app')
myApp.use(store).use(router).mount("#app");
