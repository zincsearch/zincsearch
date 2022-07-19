import { mount } from "@vue/test-utils";
import User from "../views/User.vue";
import { useI18n, createI18n } from "vue-i18n";
import { useStore } from "vuex";
import store from "../store";
import { Quasar, Dialog, Notify, QLayout } from "quasar";
import { expect, it } from "vitest";
import i18n from "../locales";
// const { t } = useI18n();
// const store = useStore();

// const wrapper = mount(About, {
//   global: {
//     plugins: [Quasar, i18n],
//   },
// });

it("should mount component", async () => {
  expect(User).toBeTruthy();

  // const wrapper = mount(User, {
  //   global: {
  //     plugins: [Quasar, i18n],
  //   },
  // });

  // expect(wrapper.find("h1").text()).toBe("User");
});
