import { mount } from "@vue/test-utils";
import Index from "../../views/Index.vue";
import store from "../../store";
import { Quasar, Dialog, Notify } from "quasar";
import { expect, it } from "vitest";
import i18n from "../../locales";
import AddUpdateIndex from "../../components/index/AddUpdateIndex.vue";
import PreviewIndex from "../../components/index/PreviewIndex.vue";

import { installQuasar } from "../helpers/install-quasar-plugin";

installQuasar();

it("should mount Index view", async () => {
  const wrapper = mount(Index, {
    shallow: false,
    components: {
      Notify,
      Dialog,
      AddUpdateIndex,
      PreviewIndex,
    },
    global: {
      plugins: [i18n, store],
    },
  });
  expect(Index).toBeTruthy();

  // console.log("Index is: ", wrapper.html());

  // expect(wrapper.text()).toContain("Index");
});
