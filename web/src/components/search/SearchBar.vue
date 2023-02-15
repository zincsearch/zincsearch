<template>
  <div class="row justify-between">
    <div class="col row justify-between">
      <div class="col">
        <q-input
          v-model="searchQuery"
          data-cy="search-bar-input"
          :placeholder="t('search.typeSearch')"
          dense
          filled
          type="search"
          class="search-field"
          @change="searchData"
          @keycode.enter="searchData"
        >
          <template #prepend>
            <q-icon name="search" />
          </template>
        </q-input>
      </div>
      <div><syntax-guide /></div>
      <div class="q-ml-sm"><date-time v-model="dateVal" /></div>
    </div>
    <div class="q-ml-sm search-time">
      <q-btn-group spread>
        <q-btn
          data-cy="search-bar-refresh-button"
          dense
          color="primary"
          label=""
          title="Search"
          :icon="refreshIcon"
          :loading="searching"
          class="q-pa-sm"
          @click="searchData"
        />
        <q-btn-dropdown
          data-cy="search-bar-button-dropdown"
          color="primary"
          class="q-pa-sm search-dropdown"
          no-caps
          :label="refreshTime"
        >
          <q-list>
            <q-item
              v-for="item in refreshTimes"
              :key="item.value"
              v-close-popup
              dense
              clickable
              @click="refreshTimeChange(item)"
            >
              <q-item-section>
                <q-item-label>{{ item.label }}</q-item-label>
              </q-item-section>
            </q-item>
          </q-list>
        </q-btn-dropdown>
      </q-btn-group>
    </div>
  </div>
</template>

<script>
import { defineComponent, nextTick, watch, ref } from "vue";
import { debounce } from "quasar";
import { useI18n } from "vue-i18n";

import DateTime from "./DateTime.vue";
import SyntaxGuide from "./SyntaxGuide.vue";

export default defineComponent({
  name: "ComponentSearchSearchBar",
  components: {
    DateTime,
    SyntaxGuide,
  },
  props: {
    data: {
      type: Object,
      default: () => ({}),
    },
  },
  emits: ["updated", "refresh"],
  setup(props, { emit }) {
    const searching = ref(false);
    const refreshIcon = ref("refresh");
    const refreshTime = ref("Off");
    const { t } = useI18n();
    const refreshTimes = [
      { label: "Off", value: 0 },
      { label: "5s", value: 5 },
      { label: "10s", value: 10 },
      { label: "15s", value: 15 },
      { label: "30s", value: 30 },
      { label: "1m", value: 60 },
      { label: "5m", value: 300 },
      { label: "15m", value: 900 },
      { label: "30m", value: 1800 },
      { label: "1h", value: 3600 },
      { label: "2h", value: 7200 },
      { label: "1d", value: 86400 },
    ];

    const refreshTimer = ref(null);
    const refreshTimeChange = (time) => {
      refreshTime.value = time.label;
      if (time.value === 0) {
        clearInterval(refreshTimer.value);
        refreshTimer.value = null;
      } else {
        if (refreshTimer.value != null) {
          clearInterval(refreshTimer.value);
        }
        refreshTimer.value = setInterval(() => {
          emit("refresh", "");
        }, time.value * 1000);
      }
    };

    const searchQuery = ref(props.data.query);

    const dateVal = ref({
      tab: "relative",
      startDate: "",
      startTime: "",
      endDate: "",
      endTime: "",

      selectedRelativePeriod: "Minutes",
      selectedRelativeValue: 30,
      selectedFullTime: false,
    });

    const searchData = () => {
      if (searching.value) {
        return;
      }
      refreshIcon.value = "refresh";
      searching.value = true;
      emit("updated", {
        query: searchQuery.value,
        time: dateVal.value,
      });
      nextTick(() => {
        searching.value = false;
      });
    };

    const dSearchData = debounce(searchData, 1000);

    // when the datetime filter changes then update the results
    watch(dateVal.value, (v) => {
      refreshIcon.value = "search";
      if (v.startDate != "" && v.endDate != "") {
        dSearchData();
      }
    });

    const setSearchQuery = (query) => {
      searchQuery.value = query;
    };

    return {
      t,
      searching,
      searchQuery,
      dateVal,
      refreshIcon,
      refreshTime,
      refreshTimes,
      refreshTimer,
      refreshTimeChange,
      searchData,
      setSearchQuery,
    };
  },
  beforeUnmount() {
    clearInterval(this.refreshTimer);
    this.refreshTimer = null;
  },
});
</script>

<style lang="scss">
.search-field .q-field__control {
  border-radius: 4px 0 0 0;
}
.search-time {
  width: 120px;
}
.search-dropdown {
  width: 60px !important;
}
</style>
