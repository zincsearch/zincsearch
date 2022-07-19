import { mount } from "@vue/test-utils";
import DateTime from "../../../components/search/DateTime.vue";
import store from "../../../store";
import { Quasar, Dialog, Notify } from "quasar";
import { expect, it } from "vitest";
import i18n from "../../../locales";

it("should mount DateTime component", async () => {
  const wrapper = mount(DateTime, {
    shallow: false,
    components: {
      Notify,
      Dialog,
    },
    global: {
      plugins: [Quasar, i18n, store],
    },
  });
  expect(DateTime).toBeTruthy();

  console.log("DateTime is: ", wrapper.html());

  // expect(wrapper.text()).toContain("DateTime");
});
