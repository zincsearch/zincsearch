import { mount } from "@vue/test-utils";
import { expect, it } from "vitest";
import { Quasar, Dialog, Notify } from "quasar";
// import ApexCharts from "apexcharts";

import SearchResult from "../../../components/search/SearchResult.vue";
import HighLight from "../../../components/HighLight.vue";
import store from "../../../store";
import i18n from "../../../locales";

it("should mount SearchResult component", async () => {
  const wrapper = mount(SearchResult, {
    shallow: false,
    components: {
      Notify,
      Dialog,
      HighLight,
    },
    global: {
      plugins: [Quasar, i18n, store],
    },
  });
  expect(SearchResult).toBeTruthy();

  // console.log("SearchResult is: ", wrapper.html());

  // expect(wrapper.text()).toContain("SearchResult");
});
