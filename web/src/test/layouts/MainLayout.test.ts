import { mount } from "@vue/test-utils";
import MainLayout from "../../layouts/MainLayout.vue";
import store from "../../store";
import { Quasar, Dialog, Notify } from "quasar";
import { expect, it } from "vitest";
import i18n from "../../locales";

it("should mount MainLayout view", async () => {
  const wrapper = mount(MainLayout, {
    shallow: true,
    global: {
      plugins: [Quasar, i18n, store],
    },
  });
  expect(MainLayout).toBeTruthy();

  // console.log("MainLayout is: ", wrapper.html());

  // expect(wrapper.text()).toContain("MainLayout");
});
