<template>
  <div class="column index-menu">
    <div>
      <q-select
        v-model="selectedIndex"
        :options="options"
        filled
        dense
        use-input
        input-debounce="0"
        label="Select Index"
        behavior="menu"
        class="q-mt-md q-mb-sm"
        @filter="filterFn"
        @update:model-value="selectFn"
      >
        <template #no-option>
          <q-item>
            <q-item-section class="text-grey"> No results </q-item-section>
          </q-item>
        </template>
      </q-select>
    </div>
    <div>
      <q-table
        v-model:selected="selectedFields"
        :rows="indexFields"
        row-key="name"
        :filter="filterField"
        :filter-method="filterFieldFn"
        dense
        hide-header
        hide-bottom
        selection="multiple"
        @row-click="clickFieldFn"
        @update:selected="selectedFieldFn"
      >
        <template #top-right>
          <q-input
            v-model="filterField"
            filled
            borderless
            dense
            clearable
            debounce="1"
            placeholder="Search for a field"
          >
            <template #append>
              <q-icon name="search" />
            </template>
          </q-input>
        </template>
      </q-table>
    </div>
  </div>
</template>

<script>
import { defineComponent, ref } from "vue";
import indexService from "../../services/index";

export default defineComponent({
  name: "ComponentSearchIndexSelect",
  props: {
    data: {
      type: Object,
      default: () => ({}),
    },
  },
  emits: ["updated"],
  setup(props, { emit }) {
    const getIndexData = (field) => props.data[field];
    const selectedIndex = ref(getIndexData("name"));
    const selectedFields = ref(getIndexData("columns"));
    const indexList = ref([]);
    const indexFields = ref([]);
    const mappingList = ref({});
    const options = ref([]);

    // index operation
    const filterFn = (val, update) => {
      if (val === "") {
        update(() => {
          options.value = indexList.value;
        });
        return;
      }

      update(() => {
        const needle = val.toLowerCase();
        options.value = indexList.value.filter((v) =>
          v.value.toLowerCase().includes(needle)
        );
      });
    };

    const selectFn = (index) => {
      if (!index || !index.value) {
        return;
      }
      selectedFields.value = [];
      indexFields.value = [];
      for (var k in mappingList.value[index.value]) {
        if (k == "_id" || k == "@timestamp") {
          continue;
        }
        indexFields.value.push({ name: k });
      }

      emit("updated", {
        name: index.value,
        columns: [],
      });
    };

    // fields operation
    const filterField = ref("");
    const filterFieldFn = (rows, terms) => {
      var filtered = [];
      terms = terms.toLowerCase();
      for (var i = 0; i < rows.length; i++) {
        if (rows[i]["name"].toLowerCase().includes(terms)) {
          filtered.push(rows[i]);
        }
      }
      return filtered;
    };
    const clickFieldFn = (evt, row, index) => {
      if (selectedFields.value.includes(row)) {
        selectedFields.value = selectedFields.value.filter(
          (v) => v.name !== row.name
        );
      } else {
        selectedFields.value.push(row);
      }
      emit("updated", {
        name: selectedIndex.value.value,
        columns: selectedFields.value.map((v) => v.name),
      });
    };
    const selectedFieldFn = () => {
      emit("updated", {
        name: selectedIndex.value.value,
        columns: selectedFields.value.map((v) => v.name),
      });
    };

    // get the list of indices from server when the component is mounted
    const getIndexList = () => {
      indexService.list().then((res) => {
        res.data.map((item) => {
          indexList.value.push({
            label: item.name,
            value: item.name,
          });
          mappingList.value[item.name] = item.mappings
            ? item.mappings.properties
            : [];
        });
        selectedIndex.value = indexList.value[0];
        selectFn(selectedIndex.value);
      });
    };

    getIndexList();

    return {
      selectedIndex,
      selectedFields,
      mappingList,
      options,
      filterFn,
      selectFn,
      indexFields,
      filterField,
      filterFieldFn,
      clickFieldFn,
      selectedFieldFn,
    };
  },
});
</script>

<style lang="scss">
.index-menu {
  width: 220px;
}
</style>
