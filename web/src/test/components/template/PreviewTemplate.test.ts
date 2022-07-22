import { mount } from "@vue/test-utils";
import { expect, it } from "vitest";
import { Quasar, Dialog, Notify } from "quasar";

import store from "../../../store";
import i18n from "../../../locales";
import JsonEditor from "../../../components/JsonEditor.vue";
import PreviewTemplate from "../../../components/template/PreviewTemplate.vue";

it("should mount PreviewTemplate component", async () => {
  const wrapper = mount(PreviewTemplate, {
    shallow: false,
    components: {
      Notify,
      Dialog,
      JsonEditor,
    },
    global: {
      plugins: [Quasar, i18n, store],
    },
  });
  expect(PreviewTemplate).toBeTruthy();

  // console.log("PreviewTemplate is: ", wrapper.html());

  // expect(wrapper.text()).toContain("PreviewTemplate");
});
