import { mount } from "@vue/test-utils";
import { expect, it, describe } from "vitest";
import { Quasar } from "quasar";

import i18n from "../../locales";
import MenuLink from "../../components/MenuLink.vue";
import store from "../../store";

describe("MenuLink", () => {
  it("should be able to mount Menulink component", async () => {
    const wrapper = mount(MenuLink, {
      shallow: false,
      props: {
        title: "Search",
      },
      global: {
        plugins: [Quasar, i18n, store],
      },
    });
    expect(MenuLink).toBeTruthy();
  });

  it("should be able to display the titel in Menulink", async () => {
    const wrapper = mount(MenuLink, {
      shallow: false,
      props: {
        title: "Search",
      },
      global: {
        plugins: [Quasar, i18n, store],
      },
    });

    expect(wrapper.text()).toContain("Search"); // check title
  });
});
