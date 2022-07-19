import { mount } from "@vue/test-utils";
import PreviewTemplate from "../../../components/template/PreviewTemplate.vue";
import store from "../../../store";
import { Quasar, Dialog, Notify } from "quasar";
import { expect, it } from "vitest";
import i18n from "../../../locales";

it("should mount component", async () => {
  const wrapper = mount(PreviewTemplate, {
    shallow: false,
    components: {
      Notify,
      Dialog,
    },
    global: {
      plugins: [Quasar, i18n, store],
    },
  });
  expect(PreviewTemplate).toBeTruthy();

  console.log("PreviewTemplate is: ", wrapper.html());

  // expect(wrapper.text()).toContain("PreviewTemplate");
});
