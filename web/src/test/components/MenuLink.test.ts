import { mount } from "@vue/test-utils";
import MenuLink from "../../components/MenuLink.vue";
import store from "../../store";
import { Quasar, Dialog, Notify } from "quasar";
import { expect, it } from "vitest";
import i18n from "../../locales";

it("should mount MenuLink component", async () => {
  const wrapper = mount(MenuLink, {
    shallow: false,
    components: {
      Notify,
      Dialog,
    },
    global: {
      plugins: [Quasar, i18n, store],
    },
  });
  expect(MenuLink).toBeTruthy();

  console.log("MenuLink is: ", wrapper.html());

  // expect(wrapper.text()).toContain("MenuLink");
});
