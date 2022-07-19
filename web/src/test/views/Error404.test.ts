import { it, expect } from "vitest";
import { mount } from "@vue/test-utils";
import { Quasar, Notify, Dialog } from "quasar";

import Error404 from "../../views/Error404.vue";
import store from "../../store";
import i18n from "../../locales";

it("mount Error404", async () => {
  const wrapper = mount(Error404, {
    shallow: true,
    components: {
      Notify,
      Dialog,
    },
    global: {
      plugins: [Quasar, store, i18n],
    },
  });
  expect(Error404).toBeTruthy();

  expect(wrapper.html()).toContain("Oops");
  expect(wrapper.html()).toContain("404");
});
