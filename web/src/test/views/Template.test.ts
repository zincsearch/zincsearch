import { it, expect } from "vitest";
import { mount } from "@vue/test-utils";
import { Quasar, Notify } from "quasar";

import i18n from "../../locales";
import AddUpdateTemplate from "../../components/template/AddUpdateTemplate.vue";
import PreviewTemplate from "../../components/template/PreviewTemplate.vue";
import Template from "../../views/Template.vue";
import store from "../../store";

it("should mount Template view", async () => {
  const wrapper = mount(Template, {
    shallow: true,
    components: {
      Notify,
      AddUpdateTemplate,
      PreviewTemplate,
    },
    global: {
      plugins: [Quasar, i18n, store],
    },
  });
  expect(Template).toBeTruthy();

  // console.log("Template is", wrapper.html());
});
