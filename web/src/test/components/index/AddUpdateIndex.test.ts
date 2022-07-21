import { mount } from "@vue/test-utils";
import AddUpdateIndex from "../../../components/index/AddUpdateIndex.vue";
import store from "../../../store";
import { Quasar, Dialog, Notify } from "quasar";
import { expect, it } from "vitest";
import i18n from "../../../locales";

it("should mount component", async () => {
  const wrapper = mount(AddUpdateIndex, {
    shallow: false,
    components: {
      Notify,
      Dialog,
    },
    global: {
      plugins: [Quasar, i18n, store],
    },
  });
  expect(AddUpdateIndex).toBeTruthy();

  console.log("AddUpdateIndex is: ", wrapper.html());

  // expect(wrapper.text()).toContain("AddUpdateIndex");
});
