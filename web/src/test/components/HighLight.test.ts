import { mount } from "@vue/test-utils";
import HighLight from "../../components/HighLight.vue";
import store from "../../store";
import { Quasar, Dialog, Notify } from "quasar";
import { expect, it } from "vitest";
import i18n from "../../locales";

it("should mount component", async () => {
  const wrapper = mount(HighLight, {
    shallow: false,
    components: {
      Notify,
      Dialog,
    },
    global: {
      plugins: [Quasar, i18n, store],
    },
  });
  expect(HighLight).toBeTruthy();

  console.log("HighLight is: ", wrapper.html());

  // expect(wrapper.text()).toContain("HighLight");
});
