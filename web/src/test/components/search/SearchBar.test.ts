import { mount } from "@vue/test-utils";
import { Quasar, Dialog, Notify } from "quasar";
import { expect, it } from "vitest";

import i18n from "../../../locales";
import store from "../../../store";

import SearchBar from "../../../components/search/SearchBar.vue";
import DateTime from "../../../components/search/DateTime.vue";
import SyntaxGuide from "../../../components/search/SyntaxGuide.vue";

it("should mount SearchBar component", async () => {
  const wrapper = mount(SearchBar, {
    shallow: false,
    components: {
      Notify,
      Dialog,
      DateTime,
      SyntaxGuide,
    },
    global: {
      plugins: [Quasar, i18n, store],
    },
  });
  expect(SearchBar).toBeTruthy();

  // console.log("SearchBar is: ", wrapper.html());

  // expect(wrapper.text()).toContain("SearchBar");
});
