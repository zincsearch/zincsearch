import { mount } from "@vue/test-utils";
import { expect, it } from "vitest";
import { Quasar, Dialog, Notify } from "quasar";

import HighLight from "../../components/HighLight.vue";
import store from "../../store";
import i18n from "../../locales";

it("should mount HighLight component", async () => {
  const wrapper = mount(HighLight, {
    shallow: false,
    props: {
      content: "This is a test",
    },
    components: {
      Notify,
      Dialog,
    },
    global: {
      plugins: [Quasar, i18n, store],
    },
  });
  expect(HighLight).toBeTruthy();

  // console.log("HighLight is: ", wrapper.html());

  // expect(wrapper.text()).toContain("HighLight");
});
