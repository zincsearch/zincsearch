import { mount } from "@vue/test-utils";
import About from "../views/About.vue";
import { useI18n, createI18n } from "vue-i18n";
import { useStore } from "vuex";
import store from "../store";
import { Quasar, Dialog, Notify, QLayout } from "quasar";
import { expect, it } from "vitest";
import i18n from "../locales";
// const { t } = useI18n();
// const store = useStore();

const wrapper = mount(About, {
  shallow: true,
  global: {
    plugins: [Quasar, i18n, store],
  },
});

it("should mount component", async () => {
  expect(About).toBeTruthy();

  console.log("HTML is: ", wrapper.html());

  // expect(wrapper.text()).toContain("About");
});
