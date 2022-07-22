import { mount } from "@vue/test-utils";
import { Quasar, Dialog, Notify } from "quasar";
import { expect, it } from "vitest";

import i18n from "../../../locales";
import store from "../../../store";

import AddUpdateIndex from "../../../components/index/AddUpdateIndex.vue";
import JsonEditor from "../../../components/JsonEditor.vue";

it("should mount component", async () => {
  const wrapper = mount(AddUpdateIndex, {
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
  expect(AddUpdateIndex).toBeTruthy();

  // console.log("AddUpdateIndex is: ", wrapper.html());

  // expect(wrapper.text()).toContain("AddUpdateIndex");
});
