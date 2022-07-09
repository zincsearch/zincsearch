<template>
  <q-page class="q-pa-md">
    <q-table
      v-model:selected="selectedIndexes"
      :title="t('index.header')"
      :rows="indexes"
      :columns="resultColumns"
      row-key="name"
      :pagination="pagination"
      selection="multiple"
      :loading="loading"
      :filter="filterQuery"
      :filter-method="filterData"
    >
      <template #top-right>
        <q-input
          v-model="filterQuery"
          filled
          borderless
          dense
          :placeholder="t('index.search')"
        >
          <template #append>
            <q-icon name="search" class="cursor-pointer" />
          </template>
        </q-input>
        <q-btn
          class="q-ml-sm"
          color="primary"
          icon="add"
          :label="t('index.add')"
          @click="addIndex"
        />
        <q-btn
          class="q-ml-sm"
          color="negative"
          icon="delete"
          :label="t('index.delete')"
          @click="deleteSelectedIndexes"
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
import { defineComponent, nextTick, ref } from "vue";
import { useStore } from "vuex";
import { useQuasar } from "quasar";
import { useI18n } from "vue-i18n";

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
    const { t } = useI18n();
    const loading = ref(false);
    const selectedIndexes = ref([]);
    const indexes = ref([]);
    const getIndexes = () => {
      loading.value = true;
      indexService.list().then((res) => {
        var counter = 1;
        indexes.value = res.data.map((data) => {
          let storage_size = (data.storage_size / 1024).toFixed(2) + " KB";
          if (data.storage_size > 1024 * 1024) {
            storage_size = (data.storage_size / 1024 / 1024).toFixed(2) + " MB";
          }
          return {
            no: counter++,
            name: data.name,
            doc_num: data.doc_num,
            shard_num: data.shard_num,
            storage_size: storage_size,
            storage_type: data.storage_type,
            actions: {
              settings: data.settings,
              mappings: data.mappings,
            },
          };
        });
        loading.value = false;
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
        name: "doc_num",
        field: (row) => row.doc_num,
        label: "DOC_NUM",
        align: "right",
        sortable: true,
      },
      {
        name: "shard_num",
        field: (row) => row.shard_num,
        label: "SHARD_NUM",
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
        doc_num: props.row.doc_num,
        shard_num: props.row.shard_num,
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
          nextTick(getIndexes);
        });
      });
    };

    const deleteSelectedIndexes = () => {
      if (!selectedIndexes.value || selectedIndexes.value.length == 0) {
        $q.notify({
          position: "top",
          color: "warning",
          textColor: "white",
          icon: "warning",
          message: "Please select index for deletion",
        });
        return;
      }

      const showText = selectedIndexes.value
        .map((r) => "<li>" + r.name + "</li>")
        .join("");
      $q.dialog({
        title: "Delete indexes",
        message:
          "You are about to delete these indexes: <ul>" + showText + "</ul>",
        cancel: true,
        persistent: true,
        html: true,
      }).onOk(() => {
        const indexNames = selectedIndexes.value.map((r) => r.name).join(",");
        indexService.delete(indexNames).then((res) => {
          selectedIndexes.value = [];
          nextTick(getIndexes);
        });
      });
    };

    return {
      t,
      showAddIndexDialog,
      showPreviewIndexDialog,
      resultColumns,
      index,
      loading,
      selectedIndexes,
      deleteSelectedIndexes,
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
