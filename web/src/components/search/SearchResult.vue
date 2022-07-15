<template>
  <div class="col column q-my-md q-ml-md">
    <div class="search-list">
      <q-table
        ref="searchTable"
        v-model:expanded="searchResult._source"
        v-model:pagination="pagination"
        data-cy="search-result-area"
        :rows="searchResult"
        :columns="resultColumns"
        :loading="searchLoading"
        :rows-per-page-options="rowsPerPageOptions"
        wrap-cells
        :title="t('search.searchResult')"
        row-key="_id"
        @request="onRequest"
      >
        <template #top>
          <div class="chart">
            <apexchart
              ref="chartHistogram"
              width="100%"
              height="170"
              type="bar"
              :options="chartOptions"
              :series="chartOptions.series"
            ></apexchart>
          </div>
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
            <q-td width="30">
              <q-btn
                size="sm"
                color="secondary"
                round
                dense
                :icon="props.expand ? 'remove' : 'add'"
                @click="props.expand = !props.expand"
              />
            </q-td>
            <template v-for="col in props.cols" :key="col.name" :props="props">
              <q-td v-if="col.name == '@timestamp'" width="238">
                <span v-text="col.value"></span>
              </q-td>
              <q-td v-else>
                <high-light
                  :content="col.value"
                  :query-string="queryString"
                ></high-light>
              </q-td>
            </template>
          </q-tr>
          <q-tr v-show="props.expand" :props="props">
            <q-td colspan="100%">
              <pre class="expanded">
                 <high-light
                   :content="JSON.stringify(props.row, null, 2)"
                   :query-string="queryString"
                 ></high-light>
              </pre>
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
import { useI18n } from "vue-i18n";

import searchService from "../../services/search";
import HighLight from "../HighLight.vue";

export default defineComponent({
  name: "ComponentSearchSearchList",
  components: {
    HighLight,
  },
  props: {
    data: {
      type: Object,
      default: () => ({}),
    },
  },
  emits: ["updated:fields"],
  setup(props, { emit }) {
    // Accessing nested JavaScript objects and arrays by string path
    // https://stackoverflow.com/questions/6491463/accessing-nested-javascript-objects-and-arrays-by-string-path
    const { t } = useI18n();
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

    Object.deepKeys = function (o) {
      if (!(o instanceof Object)) {
        return [];
      }
      let results = [];
      let keys = Object.keys(o);
      for (var i in keys) {
        if (o[keys[i]] == undefined || o[keys[i]].length) {
          results.push(keys[i]);
        } else {
          let subKeys = Object.deepKeys(o[keys[i]]);
          if (subKeys.length > 0) {
            subKeys.forEach((key) => {
              results.push(keys[i] + "." + key);
            });
          } else {
            results.push(keys[i]);
          }
        }
      }
      return results;
    };

    const defaultColumns = () => {
      return [
        {
          name: "@timestamp",
          field: (row) =>
            date.formatDate(row["@timestamp"], "MMM DD, YYYY HH:mm:ss.SSS Z"),
          label: t("search.timestamp"),
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

    const searchTable = ref(null);
    const searchResult = ref([]);
    const resultCount = ref("");
    const resultColumns = ref(defaultColumns());
    const queryString = ref("");
    const rowsPerPageOptions = ref([5, 10, 20, 50, 100]);
    const pagination = ref({
      rowsPerPage: 20,
      sortBy: "desc",
      descending: false,
      page: 1,
      rowsNumber: 0,
    });
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
    const chartKeyFormat = ref("HH:mm:ss");
    const chartHistogram = ref(null);
    const chartOptions = {
      chart: {
        id: "search-summary",
        toolbar: {
          show: true,
        },
      },
      grid: {
        borderColor: "#eee",
        strokeDashArray: 5,
      },
      colors: ["#26A69A", "#9C27B0"],
      series: [],
      xaxis: {
        type: "numeric",
        labels: {
          show: false,
          rotateAlways: false,
          rotate: 0,
          hideOverlappingLabels: true,
        },
      },
      dataLabels: {
        enabled: false,
      },
      plotOptions: {
        bar: {
          columnWidth: "90%",
        },
      },
      title: {
        text: "",
      },
      noData: {
        text: "Loading...",
      },
    };

    const buildSearch = (queryData) => {
      let req = {
        query: {
          bool: {
            must: [],
          },
        },
        sort: ["-@timestamp"],
        from: pagination.value.page - 1,
        size: pagination.value.rowsPerPage,
      };

      let timestamps = getDateConsumableDateTime(queryData.time);
      if (timestamps.start_time || timestamps.end_time) {
        if (queryData.time.selectedFullTime) {
          chartKeyFormat.value = "HH:mm:ss";
          req.aggs = {
            histogram: {
              auto_date_histogram: {
                field: "@timestamp",
                buckets: 100,
              },
            },
          };
        } else {
          req.query.bool.must.push({
            range: {
              "@timestamp": {
                gte: timestamps.start_time.toISOString(),
                lt: timestamps.end_time.toISOString(),
                format: "2006-01-02T15:04:05Z07:00",
              },
            },
          });

          req.aggs = {
            histogram: {
              date_histogram: {
                field: "@timestamp",
                calendar_interval: "1s",
              },
            },
          };

          if (timestamps.end_time - timestamps.start_time >= 1000 * 60 * 5) {
            req.aggs.histogram.date_histogram.calendar_interval = "";
            req.aggs.histogram.date_histogram.fixed_interval = "5s";
            chartKeyFormat.value = "HH:mm:ss";
          }
          if (timestamps.end_time - timestamps.start_time >= 1000 * 60 * 10) {
            req.aggs.histogram.date_histogram.calendar_interval = "";
            req.aggs.histogram.date_histogram.fixed_interval = "10s";
            chartKeyFormat.value = "HH:mm:ss";
          }
          if (timestamps.end_time - timestamps.start_time >= 1000 * 60 * 30) {
            req.aggs.histogram.date_histogram.calendar_interval = "";
            req.aggs.histogram.date_histogram.fixed_interval = "30s";
            chartKeyFormat.value = "HH:mm:ss";
          }
          if (timestamps.end_time - timestamps.start_time >= 1000 * 60 * 60) {
            req.aggs.histogram.date_histogram.calendar_interval = "1m";
            req.aggs.histogram.date_histogram.fixed_interval = "";
            chartKeyFormat.value = "HH:mm";
          }
          if (timestamps.end_time - timestamps.start_time >= 1000 * 3600 * 3) {
            req.aggs.histogram.date_histogram.calendar_interval = "";
            req.aggs.histogram.date_histogram.fixed_interval = "5m";
            chartKeyFormat.value = "MM-DD HH:mm";
          }
          if (timestamps.end_time - timestamps.start_time >= 1000 * 3600 * 24) {
            req.aggs.histogram.date_histogram.calendar_interval = "1h";
            req.aggs.histogram.date_histogram.fixed_interval = "";
            chartKeyFormat.value = "MM-DD HH:mm";
          }
          if (timestamps.end_time - timestamps.start_time >= 1000 * 86400 * 7) {
            req.aggs.histogram.date_histogram.calendar_interval = "1d";
            req.aggs.histogram.date_histogram.fixed_interval = "";
            chartKeyFormat.value = "YYYY-MM-DD";
          }
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

    let lastIndexName = "";
    const searchLoading = ref(false);
    let lastIndexData = {};
    let lastQueryData = {};
    const searchData = (indexData, queryData) => {
      if (searchLoading.value) {
        return false;
      }
      lastIndexData = indexData;
      lastQueryData = queryData;
      searchLoading.value = true;
      const query = buildSearch(queryData);

      if (!indexData.name) {
        indexData.name = "";
      }

      queryString.value = queryData.query;
      searchService
        .search({ index: indexData.name, query: query })
        .then((res) => {
          if (lastIndexName != "" && lastIndexName != indexData.name) {
            resetColumns(indexData);
          }
          lastIndexName = indexData.name;

          let results = [];
          const hits = res.data.hits;
          if (hits.hits) {
            results = hits.hits;
            // update index fields
            let fields = {};
            results.forEach((row) => {
              let keys = Object.deepKeys(row._source);
              for (let i in keys) {
                fields[keys[i]] = {};
              }
            });
            emit("updated:fields", Object.keys(fields));
          }

          nextTick(() => {
            searchResult.value = results;
            resultCount.value = `Found ${hits.total.value} hits in ${res.data.took} ms`;
            searchLoading.value = false;

            pagination.value.rowsNumber = hits.total.value;

            // rerender the chart
            nextTick(() => {
              if (!res.data.aggregations) {
                console.log("res.data.aggregations is null");
                return;
              }
              const interval = res.data.aggregations.histogram["interval"];
              if (interval) {
                if (interval.includes("s")) {
                  chartKeyFormat.value = "HH:mm:ss";
                } else if (interval.includes("m")) {
                  chartKeyFormat.value = "HH:mm";
                } else if (interval.includes("h")) {
                  chartKeyFormat.value = "MM-DD HH:mm";
                } else if (interval.includes("d")) {
                  chartKeyFormat.value = "YYYY-MM-DD";
                }
              }
              chartHistogram.value.updateOptions({
                title: {
                  text: resultCount.value,
                },
                xaxis: {
                  type: "numeric",
                  labels: {
                    show: res.data.hits.total.value > 0 ? true : false,
                  },
                },
                series: [
                  {
                    name: "Count",
                    data: res.data.aggregations.histogram.buckets.map(
                      (bucket) => {
                        return {
                          x: date.formatDate(bucket.key, chartKeyFormat.value),
                          y: parseInt(bucket.doc_count, 10),
                        };
                      }
                    ),
                  },
                ],
              });
            });
          });
        })
        .catch((err) => {
          // handle the errors so as to continue using the applications
          console.log(err.message);
          searchLoading.value = false;
        });
    };

    const onRequest = (props) => {
      const { page, rowsPerPage, sortBy, descending } = props.pagination;
      pagination.value.page = page;
      pagination.value.rowsPerPage = rowsPerPage;
      pagination.value.sortBy = sortBy;
      pagination.value.descending = descending;
      searchData(lastIndexData, lastQueryData);
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
          field: (row) => {
            if (["_id", "_index", "_score"].includes(indexData.columns[i])) {
              return row[indexData.columns[i]];
            } else {
              return Object.byString(row._source, indexData.columns[i]);
            }
          },
          align: "left",
          sortable: true,
        };

        resultColumns.value.push(newCol);
      }
    };

    return {
      t,
      searchTable,
      searchData,
      resetColumns,
      resultColumns,
      searchResult,
      resultCount,
      searchLoading,
      rowsPerPageOptions,
      pagination,
      onRequest,
      chartHistogram,
      chartOptions,
      queryString,
    };
  },
});
</script>

<style lang="scss">
.max-result {
  width: 170px;
}
.search-list {
  width: 100%;
  .chart {
    width: 100%;
    border-bottom: 1px solid rgba(0, 0, 0, 0.12);
  }
  .q-table__top {
    padding: 5px 0 0 0;
  }
  .q-table thead tr,
  .q-table tbody td {
    height: 38px;
    padding: 6px 12px;
  }
  .q-table__bottom {
    width: 100%;
  }
  .q-table__bottom {
    min-height: 40px;
    padding-top: 0;
    padding-bottom: 0;
  }
  .q-td {
    word-wrap: break-word;
    word-break: break-all;
    white-space: pre-wrap;
    .expanded {
      margin: 0;
      white-space: pre-wrap;
      word-wrap: break-word;
      word-break: break-all;
    }
  }
  .highlight {
    background-color: rgb(255, 213, 0);
  }
}
</style>
