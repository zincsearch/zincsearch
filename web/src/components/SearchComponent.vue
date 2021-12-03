<template>
  <div class="search">
    <q-form @submit="searchData">
      <div class="row">
        <div class="col-7">
          <q-input
            v-model="search_query"
            label="Type your search query here"
            :dense="true"
            filled
            type="search"
            class="search-field"
          >
            <template v-slot:append>
              <q-icon name="search" />
            </template>
          </q-input>
        </div>
        <div class="col-4 search-controls">
          <SyntaxGuide />
          <DateTime v-model="dateVal" />
        </div>
        <div class="col-1">
          <q-btn
            color="secondary search-button"
            label="Search"
            type="submit"
            class="search-submit-button"
            @submit="searchData"
          >
          </q-btn>
        </div>
      </div>
      <div class="row">
        <div class="col-2">
          <q-select
            filled
            :dense="true"
            v-model="selectedIndex"
            use-input
            input-debounce="0"
            label="Select Index"
            :options="options"
            @filter="filterFn"
            @update:model-value="selectFn"
            behavior="menu"
            class="index-field"
          >
            <template v-slot:no-option>
              <q-item>
                <q-item-section class="text-grey"> No results </q-item-section>
              </q-item>
            </template>
          </q-select>

          <q-table
            :rows="index_fields"
            row-key="fields"
            :filter="filter_query"
            :filter-method="filter_method"
            dense
            selection="multiple"
            v-model:selected="selectedFields"
            class="field-list"
            :pagination="pagination_fields"
          >
            <template v-slot:top-right>
              <q-input
                filled
                borderless
                dense
                debounce="1"
                v-model="filter_query"
                placeholder="Search for a field"
              >
                <template v-slot:append>
                  <q-icon name="search" />
                </template>
              </q-input>
            </template>
          </q-table>
        </div>

        <div class="col-10">
          <q-table
            :rows="search_result"
            :columns="result_columns"
            title="Search Results"
            v-model:expanded="search_result._source"
            row-key="_id"
            dense
            class="table-class"
            :pagination="pagination"
          >
            <template v-slot:top-right>
              <div class="text-subtitle1">{{ resultCount }}</div>
            </template>

            <template v-slot:header="props">
              <q-tr :props="props">
                <q-th auto-width />
                <q-th v-for="col in props.cols" :key="col.name" :props="props">
                  {{ col.label }}
                </q-th>
              </q-tr>
            </template>

            <template v-slot:body="props">
              <q-tr :props="props">
                <q-td auto-width>
                  <q-btn
                    size="sm"
                    color="accent"
                    round
                    dense
                    @click="props.expand = !props.expand"
                    :icon="props.expand ? 'remove' : 'add'"
                  />
                </q-td>
                <q-td v-for="col in props.cols" :key="col.name" :props="props">
                  {{ col.value }}
                </q-td>
              </q-tr>
              <q-tr v-show="props.expand" :props="props">
                <q-td colspan="100%">
                  <pre class="expanded"
                    >{{ JSON.stringify(props.row, null, 2) }}
                    </pre
                  >
                </q-td>
              </q-tr>
            </template>
          </q-table>
        </div>
      </div>
    </q-form>

    <div class="row">
      <div class="col-2"></div>
    </div>
  </div>
</template>

<script>
import { ref } from "vue";
import axios from "../axios";
import { date } from "quasar";
import { useStore } from "vuex";
import router from "../router";

// @ is an alias to /src
import SyntaxGuide from "@/components/SyntaxGuide.vue";
import DateTime from "@/components/DateTime.vue";

export default {
  components: { SyntaxGuide, DateTime },
  watch: {
    selectedFields(newVal) {
      // @timestamp should always be shown
      this.result_columns = [
        {
          name: "@timestamp",
          field: (row) =>
            date.formatDate(
              row["@timestamp"],
              "MMM DD, YYYY HH:mm:ss.SSS UTC Z"
            ),
          label: "@timestamp",
          align: "left",
          sortable: true,
        },
      ];

      // add all the selected fields one by one
      for (let i = 0; i < newVal.length; i++) {
        var newCol = {
          name: newVal[i].fields,
          field: (row) => {
            return Object.byString(row._source, newVal[i].fields);
          },
          label: newVal[i].fields,
          align: "left",
          sortable: true,
        };

        this.result_columns.push(newCol);
      }

      // show _source field if no other fields are selected
      if (newVal.length == 0) {
        var source_column = {
          name: "_source",
          field: (search_result) =>
            JSON.stringify(search_result._source).substring(0, 150) + " ...",
          label: "_source",
          align: "left",
          sortable: true,
        };
        this.result_columns.push(source_column);
      }
    },
  },
  setup() {
    // Accessing nested JavaScript objects and arrays by string path
    // https://stackoverflow.com/questions/6491463/accessing-nested-javascript-objects-and-arrays-by-string-path
    Object.byString = function (o, s) {
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

    const store = useStore();

    if (window.location.origin != "http://localhost:8080") {
      store.dispatch("endpoint", window.location.origin + "/");
    }

    const indexList = ref([]);
    const options = ref([]);
    const search_query = ref("");
    const resultCount = ref("");
    const search_result = ref([]);
    const selectedIndex = ref("");
    const index_fields = ref([]);
    const selectedFields = ref([]);
    const mapping_list = ref({});
    const filter_query = ref("");
    const dateVal = ref({
      tab: "relative",

      startDate: "",
      startTime: "",
      endDate: "",
      endTime: "",

      selectedRelativePeriod: "Minutes",
      selectedRelativeValue: 30,
    });

    // show 2 fields (@timestamp abd _source) by default
    const result_columns = ref([
      {
        name: "@timestamp",
        field: (search_result) =>
          date.formatDate(
            search_result["@timestamp"],
            "MMM DD, YYYY HH:mm:ss.SSS UTC Z"
          ),
        label: "@timestamp",
        align: "left",
        sortable: true,
      },
      {
        name: "_source",
        field: (search_result) =>
          JSON.stringify(search_result).substring(0, 150) + " ...",
        label: "_source",
        align: "left",
        sortable: true,
      },
    ]);

    // get the list of indices from server when the component is mounted
    const getIndexList = async function () {
      var response = {};

      try {
        response = await axios.get(store.state.API_ENDPOINT + "api/index");
        var data = response.data;

        for (var index in data) {
          indexList.value.push(data[index].name);
          mapping_list.value[data[index].name] = data[index].mapping;
        }
      } catch (error) {
        console.log(error);
        if (error.response.status == 401) {
          router.push("/login");
        }
      }
    };

    getIndexList();

    // get the normalized date and time from the dateVal object
    const getDateConsumableDateTime = function () {
      if (dateVal.value.tab == "relative") {
        var period = "";
        var periodValue = 0;

        // quasar does not support arithmetic on weeks. convert to days.
        if (dateVal.value.selectedRelativePeriod.toLowerCase() == "weeks") {
          period = "days";
          periodValue = dateVal.value.selectedRelativeValue * 7;
        } else {
          period = dateVal.value.selectedRelativePeriod.toLowerCase();
          periodValue = dateVal.value.selectedRelativeValue;
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

        if (dateVal.value.startDate == "" && dateVal.value.startTime == "") {
          start = new Date();
        } else {
          start = new Date(
            dateVal.value.startDate + " " + dateVal.value.startTime
          );
        }

        if (dateVal.value.endDate == "" && dateVal.value.endTime == "") {
          end = new Date();
        } else {
          end = new Date(dateVal.value.endDate + " " + dateVal.value.endTime);
        }

        var rVal = {
          start_time: start,
          end_time: end,
        };
        return rVal;
      }
    };

    // expose to template
    return {
      // variables
      dateVal,
      model: ref(null),
      indexList,
      options,
      resultCount,
      selectedIndex,
      search_result,
      mapping_list,
      selectedFields,
      search_query,
      pagination: {
        rowsPerPage: 20, // current rows per page being displayed
      },
      pagination_fields: {
        rowsPerPage: 100, // current rows per page being displayed in the fields section
      },
      result_columns,
      index_fields,
      filter_query,
      getDateConsumableDateTime,

      // methods

      handleDateUpdate(newDate) {
        // this.model.value.date = newDate;
        console.log(newDate);
      },

      // filter for fields
      filter_method(rows, terms) {
        terms = terms.toLowerCase();
        var filtered_rows = [];
        for (var i = 0; i < rows.length; i++) {
          if (rows[i]["fields"].toLowerCase().includes(terms)) {
            filtered_rows.push(rows[i]);
          }
        }

        return filtered_rows;
      },

      // filter the values when value is being typed in index field
      filterFn(val, update) {
        if (val === "") {
          update(() => {
            options.value = indexList.value;
          });
          return;
        }

        update(() => {
          const needle = val.toLowerCase();
          options.value = indexList.value.filter(
            (v) => v.toLowerCase().indexOf(needle) > -1
          );
        });
      },

      // fired when the user selects an index
      // display the fields in the index that we got from the mapping
      selectFn(index) {
        // Clear the selected fields to display
        selectedFields.value = [];

        // Clear the results
        search_result.value = [];

        // display the fields in the index that we got from the mapping
        var field_list = [];
        for (var k in mapping_list.value[index]) {
          if (k != "_id") {
            field_list.push({ fields: k });
          }
        }
        index_fields.value = field_list;
      },

      // search the index with the query
      searchData() {
        var timestamps = getDateConsumableDateTime();
        var req = {
          search_type: "querystring",
          query: {
            term: search_query.value,
            start_time: timestamps.start_time.toISOString(),
            end_time: timestamps.end_time.toISOString(),
          },
          fields: ["_all"],
        };

        var url =
          store.state.API_ENDPOINT + "api/" + selectedIndex.value + "/_search";

        axios
          .post(url, req)
          .then((res) => {
            var results = [];

            if (res.data.hits.hits) {
              results = res.data.hits.hits;
            } else {
              results = [];
            }

            search_result.value = results;
            resultCount.value =
              "Found " +
              res.data.hits.total.value.toLocaleString() +
              " records in " +
              res.data.took +
              " milliseconds";
          })
          .catch((error) => {
            if (error.response.status == 401) {
              router.push("/login");
            }
          });
      },
      calculateStartAndEndDateTime1() {
        console.log("hello");
        // var start_time = "";
        // var end_time = "";
        // var start_date = "";
        // var end_date = "";
        // var start_time_obj = {};
        // var end_time_obj = {};
        // var start_date_obj = {};
        // var end_date_obj = {};

        // if (dateVal.tab == "absolute") {
        //   start_date = dateVal.startDate;
        //   end_date = dateVal.endDate;
        //   start_time = dateVal.startTime;
        //   end_time = dateVal.endTime;

        //   start_date_obj = date.parseDate(start_date, "YYYY-MM-DD");
        //   end_date_obj = date.parseDate(end_date, "YYYY-MM-DD");
        //   start_time_obj = date.parseDate(start_time, "HH:mm:ss.SSS");
        //   end_time_obj = date.parseDate(end_time, "HH:mm:ss.SSS");

        //   start_time_obj.setFullYear(start_date_obj.getFullYear());
        //   start_time_obj.setMonth(start_date_obj.getMonth());
        //   start_time_obj.setDate(start_date_obj.getDate());

        //   end_time_obj.setFullYear(end_date_obj.getFullYear());
        //   end_time_obj.setMonth(end_date_obj.getMonth());
        //   end_time_obj.setDate(end_date_obj.getDate());

        //   start_time = date.formatDate(
        //     start_time_obj,
        //     "YYYY-MM-DDTHH:mm:ss.SSSZ"
        //   );
        //   end_time = date.formatDate(end_time_obj, "YYYY-MM-DDTHH:mm:ss.SSSZ");
        // } else {
        //   console.log("dateVal.tab == relative");
        //   console.log(dateVal.selectedRelativePeriod);
        //   console.log(dateVal.selectedRelativeValue);
        //   start_time = dateVal.startTime;
        //   end_time = dateVal.endTime;

        //   start_time_obj = date.parseDate(start_time, "HH:mm:ss.SSS");
        // }
      },
    };
  },
  name: "SearchComponent",
};
</script>

<style scoped>
.field-list {
  width: 90%;
  /* margin: 10px; */
  margin-left: 10px;
  margin-right: 10px;
  /* margin-left: 10px; */
}

.result-list {
  width: 98%;
  margin-top: 10px;
  font-family: "Roboto Mono", Consolas, Menlo, Courier, monospace;
  font-size: 10px;
}

.search-field {
  width: 98%;
  margin-bottom: 10px;
  margin-top: 10px;
  margin-left: 10px;
  /* margin: 10px; */
}
.index-field {
  width: 90%;
  margin-left: 10px;
  margin-bottom: 10px;
  margin-top: 10px;
  /* margin: 10px; */
}
.search-button {
  /* margin-left: 5px; */
  width: 30%;
  margin-bottom: 10px;
  margin-top: 10px;
  /* margin: 10px; */
}

.syntax-guide {
  /* margin-left: 5px; */
  /* width: 80%; */
  /* margin-bottom: 10px; */
  margin-top: 10px;
  height: 37px;
  margin-right: 5px;
  /* margin: 10px; */
}

.table-class {
  width: 99%;
  /* font-family: "Roboto Mono", Consolas, Menlo, Courier, monospace; */
  font-size: 8px;
  table-layout: fixed;
  white-space: normal;
  margin-top: 10px;
}

.q-table--no-wrap th,
.q-table--no-wrap td {
  white-space: normal;
}

.expanded {
  /* overflow: hidden;
  overflow-wrap: break-word;
  flex-wrap: wrap;
  word-wrap: normal; */
  white-space: pre-wrap; /* Since CSS 2.1 */
}

.result-count {
  white-space: pre-wrap;
}

.search-controls {
  display: flex;
  width: 100%;
}

.search-submit-button {
  width: 100px;
}
</style>
