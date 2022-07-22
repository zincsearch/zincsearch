import { mount } from "@vue/test-utils";
import { expect, it } from "vitest";
import { Quasar, Dialog, Notify } from "quasar";

import i18n from "../../../locales";
import AddUpdateUser from "../../../components/user/AddUpdateUser.vue";
import JsonEditor from "../../../components/JsonEditor.vue";
import store from "../../../store";

it("should mount component", async () => {
  const wrapper = mount(AddUpdateUser, {
    shallow: false,
    components: {
      Notify,
      Dialog,
      JsonEditor,
    },
    global: {
      plugins: [Quasar, i18n, store],
    },
  });
  expect(AddUpdateUser).toBeTruthy();
  // expect(AddUpdateUser).toBeTruthy();

  // console.log("AddUpdateUser is: ", wrapper.html());
  // console.log("AddUpdateUser is: ", wrapper.html());

  // expect(wrapper.text()).toContain("AddUpdateUser");
  // expect(wrapper.text()).toContain("AddUpdateUser");
});
