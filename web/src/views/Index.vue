<template>
  <q-page class="q-pa-md">
    <q-table
      title="Indexes"
      :rows="indexes"
      :columns="resultColumns"
      row-key="id"
      :pagination="pagination"
      :filter="filterQuery"
      :filter-method="filterData"
    >
      <template #top-right>
        <q-input
          v-model="filterQuery"
          filled
          borderless
          dense
          placeholder="Search index"
        >
          <template #append>
            <q-icon name="search" class="cursor-pointer" />
          </template>
        </q-input>
        <q-btn
          class="q-ml-sm"
          color="primary"
          icon="add"
          label="Add index"
          @click="addIndex"
        />
      </template>

      <template #body-cell-no="props">
        <q-td :props="props" width="80">
          {{ props.value }}
        </q-td>
      </template>
      <template #body-cell-name="props">
        <q-td :props="props" auto-width>
          <a
            class="text-primary text-decoration-none"
            @click="previewIndex(props)"
          >
            {{ props.value }}
          </a>
        </q-td>
      </template>
      <template #body-cell-actions="props">
        <q-td :props="props" auto-width>
          <q-btn
            dense
            unelevated
            size="sm"
            color="blue-5"
            class="action-button"
            icon="description"
            @click="previewIndex(props)"
          />
          <q-btn
            dense
            unelevated
            size="sm"
            color="red-5"
            class="action-button q-ml-sm"
            icon="delete"
            @click="deleteIndex(props)"
          />
        </q-td>
      </template>
    </q-table>

    <q-dialog
      v-model="showAddIndexDialog"
      position="right"
      full-height
      seamless
      maximized
    >
      <add-update-index @updated="indexAdded" />
    </q-dialog>

    <q-dialog
      v-model="showPreviewIndexDialog"
      position="right"
      full-height
      maximized
    >
      <preview-index v-model="index" />
    </q-dialog>
  </q-page>
</template>

<script>
import { defineComponent, ref } from "vue";
import { useStore } from "vuex";
import { useQuasar } from "quasar";
import indexService from "../services/index";

import AddUpdateIndex from "../components/index/AddUpdateIndex.vue";
import PreviewIndex from "../components/index/PreviewIndex.vue";

export default defineComponent({
  name: "PageIndex",
  components: {
    AddUpdateIndex,
    PreviewIndex,
  },
  setup() {
    const store = useStore();
    const $q = useQuasar();

    const indexes = ref([]);
    const getIndexes = () => {
      indexService.list().then((res) => {
        var counter = 1;
        indexes.value = res.data.map((data) => {
          let storage_size = data.storage_size + " KB";
          if (data.storage_size > 1024) {
            storage_size = (data.storage_size / 1024).toFixed(2) + " MB";
          }
          return {
            no: counter++,
            name: data.name,
            docs_count: data.docs_count,
            storage_size: storage_size,
            storage_type: data.storage_type,
            actions: {
              settings: data.settings,
              mappings: data.mappings,
            },
          };
        });
      });
    };

    getIndexes();

    const resultColumns = [
      {
        name: "no",
        field: (row) => row.no,
        label: "#",
        align: "right",
        sortable: true,
      },
      {
        name: "name",
        field: (row) => row.name,
        label: "NAME",
        align: "left",
        sortable: true,
      },
      {
        name: "docs_count",
        field: (row) => row.docs_count,
        label: "DOCS_COUNT",
        align: "right",
        sortable: true,
      },
      {
        name: "storage_size",
        field: (row) => row.storage_size,
        label: "STORAGE_SIZE",
        align: "right",
        sortable: true,
      },
      {
        name: "storage_type",
        field: (row) => row.storage_type,
        label: "STORAGE_TYPE",
        align: "left",
        sortable: true,
      },
      {
        name: "actions",
        field: (row) => row.actions,
        label: "ACTIONS",
        align: "left",
        sortable: true,
      },
    ];

    const index = ref({});
    const showAddIndexDialog = ref(false);
    const showPreviewIndexDialog = ref(false);

    const addIndex = () => {
      showAddIndexDialog.value = true;
    };
    const previewIndex = (props) => {
      index.value = {
        name: props.row.name,
        docs_count: props.row.docs_count,
        storage_size: props.row.storage_size,
        storage_type: props.row.storage_type,
        settings: props.row.actions.settings,
        mappings: props.row.actions.mappings,
      };
      showPreviewIndexDialog.value = true;
    };
    const deleteIndex = (props) => {
      $q.dialog({
        title: "Delete index",
        message:
          "You are about to delete this index: <ul><li>" +
          props.row.name +
          "</li></ul>",
        cancel: true,
        persistent: true,
        html: true,
      }).onOk(() => {
        indexService.delete(props.row.name).then(() => {
          getIndexes();
        });
      });
    };

    return {
      showAddIndexDialog,
      showPreviewIndexDialog,
      resultColumns,
      index,
      indexes,
      pagination: {
        rowsPerPage: 20,
      },
      filterQuery: ref(""),
      filterData(rows, terms) {
        var filtered = [];
        terms = terms.toLowerCase();
        for (var i = 0; i < rows.length; i++) {
          if (rows[i]["name"].toLowerCase().includes(terms)) {
            filtered.push(rows[i]);
          }
        }
        return filtered;
      },
      addIndex,
      deleteIndex,
      previewIndex,
      indexAdded() {
        showAddIndexDialog.value = false;
        getIndexes();
      },
    };
  },
});
</script>
