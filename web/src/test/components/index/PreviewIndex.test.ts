import { mount } from "@vue/test-utils";
import PreviewIndex from "../../../components/index/PreviewIndex.vue";
import store from "../../../store";
import { Quasar, Dialog, Notify } from "quasar";
import { expect, it } from "vitest";
import i18n from "../../../locales";

it("should mount component", async () => {
  const wrapper = mount(PreviewIndex, {
    shallow: false,
    components: {
      Notify,
      Dialog,
    },
    global: {
      plugins: [Quasar, i18n, store],
    },
  });
  expect(PreviewIndex).toBeTruthy();

  console.log("PreviewIndex is: ", wrapper.html());

  // expect(wrapper.text()).toContain("PreviewIndex");
});
