import { mount } from "@vue/test-utils";
import AddUpdateUser from "../../../components/user/AddUpdateUser.vue";
import store from "../../../store";
import { Quasar, Dialog, Notify } from "quasar";
import { expect, it } from "vitest";
import i18n from "../../../locales";

it("should mount component", async () => {
  const wrapper = mount(AddUpdateUser, {
    shallow: false,
    components: {
      Notify,
      Dialog,
    },
    global: {
      plugins: [Quasar, i18n, store],
    },
  });
  expect(AddUpdateUser).toBeTruthy();

  console.log("AddUpdateUser is: ", wrapper.html());

  // expect(wrapper.text()).toContain("AddUpdateUser");
});
