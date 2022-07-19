import { it, expect } from "vitest";
import { mount } from "@vue/test-utils";
import { Quasar } from "quasar";

import i18n from "../../locales";
import Search from "../../views/Search.vue";
import store from "../../store";

it("mount Search", async () => {
  const wrapper = mount(Search, {
    shallow: true,
    global: {
      plugins: [Quasar, i18n, store],
    },
  });
  expect(Search).toBeTruthy();

  console.log("Search is", wrapper.html());
});
