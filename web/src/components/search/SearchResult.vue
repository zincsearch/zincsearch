<template>
  <div class="col column q-my-md q-ml-md">
    <div class="search-list">
      <q-table
        ref="searchTable"
        v-model:expanded="searchResult._source"
        data-cy="search-result-area"
        :rows="searchResult"
        :columns="resultColumns"
        :loading="searchLoading"
        :pagination="pagination"
        wrap-cells
        :title="t('search.searchResult')"
        row-key="_id"
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
                  :content="col.value + ''"
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

        <template #bottom="scope">
          <div class="q-table__control full-width row justify-between">
            <div class="max-result">
              <q-input
                v-model="maxRecordToReturn"
                :label="t('search.maxRecords')"
                dense
                filled
                square
                type="search"
                class="search-field"
              />
            </div>
            <div class="q-table__control">
              <span class="q-table__bottom-item">Records per page:</span>
              <q-select
                v-model="pagination.rowsPerPage"
                borderless
                :options="perPageOptions"
                @update:modelValue="changePagination"
              />

              <span class="q-table__bottom-item"
                >{{
                  (scope.pagination.page - 1) * scope.pagination.rowsPerPage +
                  1
                }}-{{ scope.pagination.page * scope.pagination.rowsPerPage }} of
                {{ resultTotal }}</span
              >
              <q-btn
                icon="first_page"
                color="grey-8"
                size="sm"
                round
                dense
                flat
                :disable="scope.isFirstPage"
                @click="scope.firstPage"
              />
              <q-btn
                icon="chevron_left"
                color="grey-8"
                size="sm"
                round
                dense
                flat
                :disable="scope.isFirstPage"
                @click="scope.prevPage"
              />
              <q-btn
                icon="chevron_right"
                color="grey-8"
                size="sm"
                round
                dense
                flat
                :disable="scope.isLastPage"
                @click="scope.nextPage"
              />
              <q-btn
                icon="last_page"
                color="grey-8"
                size="sm"
                round
                dense
                flat
                :disable="scope.isLastPage"
                @click="scope.lastPage"
              />
            </div>
          </div>
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
import { byString, deepKeys } from "../../utils/json";

export default defineComponent({
  name: "SearchResult",
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
    const resultTotal = ref(0);
    const resultColumns = ref(defaultColumns());
    const queryString = ref("");

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
        text: "No data available",
      },
    };

    const buildSearch = (queryData) => {
      var req = {
        query: {
          bool: {
            must: [],
          },
        },
        sort: ["-@timestamp"],
        from: 0,
        size: parseInt(maxRecordToReturn.value, 10),
      };

      var timestamps = getDateConsumableDateTime(queryData.time);
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

    const maxRecordToReturn = ref(100);
    const selectedPerPage = ref("20");
    const perPageOptions = [
      { label: "5", value: 5 },
      { label: "10", value: 10 },
      { label: "20", value: 20 },
      { label: "50", value: 50 },
      { label: "100", value: 100 },
      { label: "All", value: 0 },
    ];
    const pagination = ref({
      rowsPerPage: 20,
    });
    const changePagination = (val) => {
      selectedPerPage.value = val.label;
      pagination.value.rowsPerPage = val.value;
      searchTable.value.setPagination(pagination.value);
    };

    let lastIndexName = "";
    const searchLoading = ref(false);
    const searchData = (indexData, queryData) => {
      if (searchLoading.value) {
        return false;
      }
      searchLoading.value = true;
      try {
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

            var results = [];
            if (res.data.hits.hits) {
              results = res.data.hits.hits;
              // update index fields
              let fields = {};
              results.forEach((row) => {
                let keys = deepKeys(row._source);
                for (let i in keys) {
                  fields[keys[i]] = {};
                }
              });
              emit("updated:fields", Object.keys(fields));
            }

            nextTick(() => {
              searchResult.value = results;
              resultTotal.value = results.length;
              resultCount.value =
                "Found " +
                res.data.hits.total.value.toLocaleString() +
                " hits in " +
                res.data.took +
                " ms";
              searchLoading.value = false;

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
                            x: date.formatDate(
                              bucket.key,
                              chartKeyFormat.value
                            ),
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
      } catch (e) {
        console.log(e.message);
        searchLoading.value = false;
      }
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
        let newCol = {
          name: indexData.columns[i],
          label: indexData.columns[i],
          field: (row) => {
            if (["_id", "_index", "_score"].includes(indexData.columns[i])) {
              return row[indexData.columns[i]];
            } else {
              return byString(row._source, indexData.columns[i]);
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
      resultTotal,
      resultCount,
      searchLoading,
      selectedPerPage,
      maxRecordToReturn,
      perPageOptions,
      pagination,
      changePagination,
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
