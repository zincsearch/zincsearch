import { mount } from "@vue/test-utils";
import { expect, it } from "vitest";
import { Quasar, Dialog, Notify } from "quasar";

import i18n from "../../../locales";
import SyntaxGuide from "../../../components/search/SyntaxGuide.vue";
import store from "../../../store";

it("should mount SyntaxGuide component", async () => {
  const wrapper = mount(SyntaxGuide, {
    shallow: false,
    components: {
      Notify,
      Dialog,
    },
    global: {
      plugins: [Quasar, i18n, store],
    },
  });
  expect(SyntaxGuide).toBeTruthy();

  // console.log("SyntaxGuide is: ", wrapper.html());

  // expect(wrapper.text()).toContain("SyntaxGuide");
});
