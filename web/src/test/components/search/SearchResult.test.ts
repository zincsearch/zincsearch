import { mount } from "@vue/test-utils";
import SearchResult from "../../../components/search/SearchResult.vue";
import store from "../../../store";
import { Quasar, Dialog, Notify } from "quasar";
import { expect, it } from "vitest";
import i18n from "../../../locales";

it("should mount component", async () => {
  const wrapper = mount(SearchResult, {
    shallow: false,
    components: {
      Notify,
      Dialog,
    },
    global: {
      plugins: [Quasar, i18n, store],
    },
  });
  expect(SearchResult).toBeTruthy();

  console.log("SearchResult is: ", wrapper.html());

  // expect(wrapper.text()).toContain("SearchResult");
});
