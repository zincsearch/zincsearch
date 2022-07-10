<template>
  <q-page class="q-pa-md">
    <div class="column">
      <search-bar
        ref="searchBar"
        :data="queryData"
        @updated="queryUpdated"
        @refresh="searchData"
      />
      <div class="row">
        <index-list
          ref="indexListRef"
          :data="indexData"
          @updated="indexUpdated"
        />
        <search-result
          ref="searchResultRef"
          @updated:fields="updateIndexFields"
        />
      </div>
    </div>
  </q-page>
</template>

<script>
import { defineComponent, ref } from "vue";

import SearchBar from "../components/search/SearchBar.vue";
import IndexList from "../components/search/IndexList.vue";
import SearchResult from "../components/search/SearchResult.vue";

export default defineComponent({
  name: "PageSearch",
  components: {
    SearchBar,
    IndexList,
    SearchResult,
  },

  setup() {
    const indexData = ref({
      name: "",
      columns: [],
    });
    const queryData = ref({
      query: "",
      time: {
        tab: "relative",
        startDate: "",
        startTime: "",
        endDate: "",
        endTime: "",
        selectedRelativePeriod: "Minutes",
        selectedRelativeValue: 30,
        selectedFullTime: false,
      },
    });

    const searchBar = ref(null);
    const indexListRef = ref(null);
    const searchResultRef = ref(null);
    const searchData = () => {
      searchResultRef.value.searchData(indexData.value, queryData.value);
    };

    const resetColumns = () => {
      searchResultRef.value.resetColumns(indexData.value);
    };

    const indexUpdated = ({ name, columns }) => {
      if (indexData.value.name != name) {
        indexData.value.name = name;
        indexData.value.columns = columns;
        queryData.value.query = "";
        searchBar.value.setSearchQuery("");
        searchData();
      } else {
        indexData.value.columns = columns;
        resetColumns();
      }
    };

    const queryUpdated = ({ query, time }) => {
      queryData.value.query = query;
      queryData.value.time = time;
      searchData();
    };

    const updateIndexFields = (fields) => {
      indexListRef.value.setIndexFields(fields);
    };

    return {
      indexData,
      queryData,
      searchBar,
      indexListRef,
      searchResultRef,
      searchData,
      indexUpdated,
      queryUpdated,
      updateIndexFields,
    };
  },
});
</script>
