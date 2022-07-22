import { mount } from "@vue/test-utils";
import { expect, it } from "vitest";

import { Quasar, Dialog, Notify } from "quasar";

import i18n from "../../../locales";
import store from "../../../store";
import PreviewIndex from "../../../components/index/PreviewIndex.vue";
import JsonEditor from "../../../components/JsonEditor.vue";

it("should mount component", async () => {
  const wrapper = mount(PreviewIndex, {
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
  expect(PreviewIndex).toBeTruthy();

  // console.log("PreviewIndex is: ", wrapper.html());

  // expect(wrapper.text()).toContain("PreviewIndex");
});
