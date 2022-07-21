import { mount } from "@vue/test-utils";
import SearchBar from "../../../components/search/SearchBar.vue";
import store from "../../../store";
import { Quasar, Dialog, Notify } from "quasar";
import { expect, it } from "vitest";
import i18n from "../../../locales";

it("should mount SearchBar component", async () => {
  const wrapper = mount(SearchBar, {
    shallow: false,
    components: {
      Notify,
      Dialog,
    },
    global: {
      plugins: [Quasar, i18n, store],
    },
  });
  expect(SearchBar).toBeTruthy();

  console.log("SearchBar is: ", wrapper.html());

  // expect(wrapper.text()).toContain("SearchBar");
});
