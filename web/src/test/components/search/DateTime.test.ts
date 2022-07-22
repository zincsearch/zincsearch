import { mount } from "@vue/test-utils";
import { Quasar, Dialog } from "quasar";
import { expect, it, describe } from "vitest";

import i18n from "../../../locales";
import DateTime from "../../../components/search/DateTime.vue";

describe("DateTime", () => {
  const wrapper = mount(DateTime, {
    shallow: false,
    components: { Dialog },
    global: {
      plugins: [Quasar, i18n],
    },
  });

  it("should mount DateTime component", async () => {
    expect(DateTime).toBeTruthy();
  });

  it("should show Minutes", async () => {
    const dateTimeButton = wrapper.find("#date-time-button");

    // console.log("item is: ", dateTimeButton);

    // console.log("displayValue: ", wrapper.vm.displayValue);
    // dateTimeButton.trigger("click");
    // console.log(wrapper.html());
    expect(wrapper.text()).toContain("-");
  });
});
