import Error404 from "../views/Error404.vue";
import store from "../store";

import { test, expect, describe } from "vitest";
import { mount } from "@vue/test-utils";
import { Quasar, QPage, QPageContainer } from "quasar";

import i18n from "../locales";

const wrapper = mount(Error404, {
  global: {
    plugins: [Quasar, store, i18n],
  },
});

test("mount Error404", () => {
  expect(Error404).toBeTruthy();

  console.log(wrapper.html());
});
