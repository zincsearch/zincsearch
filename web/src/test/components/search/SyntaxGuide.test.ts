import { mount } from "@vue/test-utils";
import SyntaxGuide from "../../../components/search/SyntaxGuide.vue";
import store from "../../../store";
import { Quasar, Dialog, Notify } from "quasar";
import { expect, it } from "vitest";
import i18n from "../../../locales";

it("should mount component", async () => {
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

  console.log("SyntaxGuide is: ", wrapper.html());

  // expect(wrapper.text()).toContain("SyntaxGuide");
});
