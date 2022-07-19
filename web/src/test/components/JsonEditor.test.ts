import { mount } from "@vue/test-utils";
import JsonEditor from "../../components/JsonEditor.vue";
import store from "../../store";
import { Quasar, Dialog, Notify } from "quasar";
import { expect, it } from "vitest";
import i18n from "../../locales";

it("should mount component", async () => {
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

  console.log("JsonEditor is: ", wrapper.html());

  // expect(wrapper.text()).toContain("JsonEditor");
});
