import { mount } from "@vue/test-utils";
import { Quasar, Dialog, Notify } from "quasar";
import { expect, it } from "vitest";

import JsonEditor from "../../components/JsonEditor.vue";
import store from "../../store";
import i18n from "../../locales";

it("should mount JsonEditor component", async () => {
  const wrapper = mount(JsonEditor, {
    shallow: false,
    components: {
      Notify,
      Dialog,
    },
    global: {
      plugins: [Quasar, i18n, store],
    },
  });
  expect(JsonEditor).toBeTruthy();

  // console.log("JsonEditor is: ", wrapper.html());

  // expect(wrapper.text()).toContain("JsonEditor");
});
