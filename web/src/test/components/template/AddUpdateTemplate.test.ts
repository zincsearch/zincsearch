import { mount } from "@vue/test-utils";
import AddUpdateTemplate from "../../../components/template/AddUpdateTemplate.vue";
import store from "../../../store";
import { Quasar, Dialog, Notify } from "quasar";
import { expect, it } from "vitest";
import i18n from "../../../locales";

it("should mount AddUpdateTemplate component", async () => {
  const wrapper = mount(AddUpdateTemplate, {
    shallow: false,
    components: {
      Notify,
      Dialog,
    },
    global: {
      plugins: [Quasar, i18n, store],
    },
  });
  expect(AddUpdateTemplate).toBeTruthy();

  console.log("AddUpdateTemplate is: ", wrapper.html());

  // expect(wrapper.text()).toContain("AddUpdateTemplate");
});
