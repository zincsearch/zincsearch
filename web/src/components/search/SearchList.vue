<template>
  <div class="col column q-my-md q-ml-md">
    <div style="display: none">histogram</div>
    <div class="search-list">
      <q-table
        v-model:expanded="searchResult._source"
        :rows="searchResult"
        :columns="resultColumns"
        :loading="searchLoading"
        :pagination="pagination"
        wrap-cells
        title="Search Results"
        row-key="_id"
      >
        <template #top-right>
          <div class="text-subtitle1">{{ resultCount }}</div>
        </template>

        <template #header="props">
          <q-tr :props="props">
            <q-th auto-width />
            <q-th v-for="col in props.cols" :key="col.name" :props="props">
              {{ col.label }}
            </q-th>
          </q-tr>
        </template>

        <template #body="props">
          <q-tr :props="props">
            <q-td auto-width>
              <q-btn
                size="sm"
                color="secondary"
                round
                dense
                :icon="props.expand ? 'remove' : 'add'"
                @click="props.expand = !props.expand"
              />
            </q-td>
            <q-td v-for="col in props.cols" :key="col.name" :props="props">
              {{ col.value }}
            </q-td>
          </q-tr>
          <q-tr v-show="props.expand" :props="props">
            <q-td colspan="100%">
              <pre class="expanded">{{
                JSON.stringify(props.row, null, 2)
              }}</pre>
            </q-td>
          </q-tr>
        </template>
      </q-table>
    </div>
  </div>
</template>

<script>
import { defineComponent, nextTick, ref } from "vue";
import { date } from "quasar";

import searchService from "../../services/search";

export default defineComponent({
  name: "ComponentSearchSearchList",
  setup() {
    // Accessing nested JavaScript objects and arrays by string path
    // https://stackoverflow.com/questions/6491463/accessing-nested-javascript-objects-and-arrays-by-string-path
    Object.byString = function (o, s) {
      if (s == undefined) {
        return "";
      }
      s = s.replace(/\[(\w+)\]/g, ".$1"); // convert indexes to properties
      s = s.replace(/^\./, ""); // strip a leading dot
      var a = s.split(".");
      for (var i = 0, n = a.length; i < n; ++i) {
        var k = a[i];
        if (typeof o == "object" && k in o) {
          o = o[k];
        } else {
          return;
        }
      }
      return o;
    };

    const defaultColumns = () => {
      return [
        {
          name: "@timestamp",
          field: (row) =>
            date.formatDate(row["@timestamp"], "MMM DD, YYYY HH:mm:ss.SSS Z"),
          label: "@timestamp",
          align: "left",
          sortable: true,
        },
        {
          name: "_source",
          field: (row) => JSON.stringify(row),
          label: "_source",
          align: "left",
          sortable: true,
        },
      ];
    };

    const searchResult = ref([]);
    const resultCount = ref("");
    const resultColumns = ref(defaultColumns());

    // get the normalized date and time from the dateVal object
    const getDateConsumableDateTime = function (dateVal) {
      if (dateVal.tab == "relative") {
        var period = "";
        var periodValue = 0;

        // quasar does not support arithmetic on weeks. convert to days.
        if (dateVal.selectedRelativePeriod.toLowerCase() == "weeks") {
          period = "days";
          periodValue = dateVal.selectedRelativeValue * 7;
        } else {
          period = dateVal.selectedRelativePeriod.toLowerCase();
          periodValue = dateVal.selectedRelativeValue;
        }
        var subtractObject = '{"' + period + '":' + periodValue + "}";

        var endTimeStamp = new Date();
        var startTimeStamp = date.subtractFromDate(
          endTimeStamp,
          JSON.parse(subtractObject)
        );

        return {
          start_time: startTimeStamp,
          end_time: endTimeStamp,
        };
      } else {
        var start, end;

        if (dateVal.startDate == "" && dateVal.startTime == "") {
          start = new Date();
        } else {
          start = new Date(dateVal.startDate + " " + dateVal.startTime);
        }

        if (dateVal.endDate == "" && dateVal.endTime == "") {
          end = new Date();
        } else {
          end = new Date(dateVal.endDate + " " + dateVal.endTime);
        }

        var rVal = {
          start_time: start,
          end_time: end,
        };
        return rVal;
      }
    };

    // whether enable histogram or not
    const getHistogram = true;
    const buildSearch = (queryData) => {
      var req = {
        query: {
          bool: {
            must: [],
          },
        },
        sort: ["-@timestamp"],
        form: 0,
        size: 100,
      };

      var timestamps = getDateConsumableDateTime(queryData.time);
      if (timestamps.start_time || timestamps.end_time) {
        if (!queryData.time.selectedFullTime) {
          req.query.bool.must.push({
            range: {
              "@timestamp": {
                gte: timestamps.start_time.toISOString(),
                lt: timestamps.end_time.toISOString(),
              },
            },
          });
        }
      }

      if (getHistogram) {
        req.aggs = {
          histogram: {
            date_histogram: {
              field: "@timestamp",
              calendar_interval: "1s",
            },
          },
        };
        console.log(
          timestamps.end_time,
          timestamps.start_time,
          timestamps.end_time - timestamps.start_time > 1000 * 60 * 60 * 24
        );
        if (timestamps.end_time - timestamps.start_time > 1000 * 60 * 60 * 1) {
          req.aggs.histogram.date_histogram.calendar_interval = "1m";
        }
        if (timestamps.end_time - timestamps.start_time > 1000 * 60 * 60 * 24) {
          req.aggs.histogram.date_histogram.calendar_interval = "1h";
        }
        if (
          timestamps.end_time - timestamps.start_time >
          1000 * 60 * 60 * 24 * 7
        ) {
          req.aggs.histogram.date_histogram.calendar_interval = "1d";
        }
      }

      if (queryData.query == "") {
        req.query.bool.must.push({
          match_all: {},
        });
        return req;
      }

      req.query.bool.must.push({
        query_string: {
          query: queryData.query,
        },
      });

      return req;
    };

    const pagination = ref({
      rowsPerPage: 20,
      // rowsNumber: 100,
    });

    let lastIndexName = "";
    const searchLoading = ref(false);
    const searchData = (indexData, queryData) => {
      if (searchLoading.value) {
        return false;
      }
      searchLoading.value = true;
      const query = buildSearch(queryData);
      searchService
        .search({ index: indexData.name, query: query })
        .then((res) => {
          if (lastIndexName != "" && lastIndexName != indexData.name) {
            resetColumns(indexData);
          }
          lastIndexName = indexData.name;

          var results = [];
          if (res.data.hits.hits) {
            results = res.data.hits.hits;
          }

          nextTick(() => {
            searchResult.value = results;
            resultCount.value =
              "Found " +
              res.data.hits.total.value.toLocaleString() +
              " records in " +
              res.data.took +
              " milliseconds";
            searchLoading.value = false;
          });
        });
    };

    const resetColumns = (indexData) => {
      resultColumns.value = defaultColumns();
      if (indexData.columns.length == 0) {
        return;
      }

      // remove _source column
      resultColumns.value.splice(1);

      // add all the selected fields one by one
      for (let i = 0; i < indexData.columns.length; i++) {
        var newCol = {
          name: indexData.columns[i],
          label: indexData.columns[i],
          field: (row) => Object.byString(row._source, indexData.columns[i]),
          align: "left",
          sortable: true,
        };

        resultColumns.value.push(newCol);
      }
    };

    return {
      searchData,
      resetColumns,
      resultColumns,
      searchResult,
      resultCount,
      searchLoading,
      pagination,
    };
  },
});
</script>

<style lang="scss">
.search-list {
  width: 100%;
  .q-table thead tr,
  .q-table tbody td {
    height: 38px;
  }
  .q-table__bottom {
    min-height: 40px;
    padding-top: 0;
    padding-bottom: 0;
  }
}
</style>
