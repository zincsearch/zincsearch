import { mount } from "@vue/test-utils";
import MainLayout from "../../layouts/MainLayout.vue";
import store from "../../store";
import router from "../../router";
import { Quasar, Dialog, Notify } from "quasar";
import { expect, it } from "vitest";
import i18n from "../../locales";
// import { useStore } from "vuex";
import { useRouter, RouterView } from "vue-router";

import { installQuasar } from "../helpers/install-quasar-plugin";

installQuasar();

it("should mount MainLayout view", async () => {
  // const router = useRouter();

  const wrapper = mount(MainLayout, {
    shallow: true,
    global: {
      plugins: [i18n, store, router],
    },
  });
  expect(MainLayout).toBeTruthy();

  // console.log("MainLayout is: ", wrapper.html());

  // expect(wrapper.text()).toContain("MainLayout");
});
