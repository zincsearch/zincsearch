import { mount } from "@vue/test-utils";
import About from "../../views/About.vue";
import store from "../../store";
import { Quasar, Dialog, Notify } from "quasar";
import { expect, it } from "vitest";
import i18n from "../../locales";
// const { t } = useI18n();
// const store = useStore();

it("should mount component", async () => {
  const wrapper = mount(About, {
    shallow: true,
    components: {
      Notify,
      Dialog,
    },
    global: {
      plugins: [Quasar, i18n, store],
    },
  });
  expect(About).toBeTruthy();

  console.log("About is: ", wrapper.html());

  // expect(wrapper.text()).toContain("About");
});
