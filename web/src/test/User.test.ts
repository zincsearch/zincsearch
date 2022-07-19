import User from "../views/User.vue";
import store from "../store";

import { test, expect, describe } from "vitest";
import { mount } from "@vue/test-utils";
import { Quasar, Notify, Dialog, useQuasar } from "quasar";
import AddUpdateUser from "../components/user/AddUpdateUser.vue";

import i18n from "../locales";
// import { useI18n } from "vue-i18n";

const wrapper = mount(User, {
  shallow: true,
  components: {
    AddUpdateUser,
  },
  global: {
    plugins: [Quasar, i18n, store],
  },
});

test("mount User", () => {
  expect(User).toBeTruthy();
  // const wrapper = wrapperFactory();

  console.log(wrapper.html());
});
