import { it, expect } from "vitest";
import { mount } from "@vue/test-utils";
import { Quasar, Notify, Dialog } from "quasar";
import { useRouter } from "vue-router";

import i18n from "../../locales";
import Login from "../../views/Login.vue";
import store from "../../store";

const router = useRouter();

it("mount Login", async () => {
  const wrapper = mount(Login, {
    shallow: true,
    components: {
      Notify,
      Dialog,
    },
    global: {
      plugins: [Quasar, i18n, store, router],
    },
  });
  expect(Login).toBeTruthy();
  // const wrapper = wrapperFactory();

  console.log("Login is", wrapper.html());
});
