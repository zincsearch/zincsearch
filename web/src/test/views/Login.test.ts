import { it, expect } from "vitest";
import { mount } from "@vue/test-utils";
import { Quasar, Notify, Dialog } from "quasar";
import { useRouter } from "vue-router";

import i18n from "../../locales";
import Login from "../../views/Login.vue";
import store from "../../store";
// import router from "../../router";

import { installQuasar } from "../helpers/install-quasar-plugin";

installQuasar();

it("should mount Login view", async () => {
  const router = useRouter();
  const wrapper = mount(Login, {
    shallow: false,
    components: {
      // Notify,
      // Dialog,
    },
    global: {
      plugins: [i18n, store],
    },
  });
  expect(Login).toBeTruthy();

  // console.log("Login is", wrapper.html());
});
