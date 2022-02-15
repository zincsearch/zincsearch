<template>
  <div class="indexmanagement">
    <q-form>
      <div class="row">
        <div class="col-3">
          <q-table
            :rows="indexList"
            row-key="name"
            dense
            @row-click="indexNameClicked"
            class="index-table"
            :pagination="pagination_indexes"
          >
          </q-table>
        </div>
        <div class="col-9">
          <q-table
            dense
            :rows="currentMapping"
            row-key="field"
            class="mapping-table"
            :pagination="pagination_fields"
          >
          </q-table>
        </div>
      </div>
    </q-form>
  </div>
</template>

<script>
import { ref } from "vue";
import axios from "../axios";
import store from "../store";

export default {
  setup() {
    const indexStruct = ref([]);
    const indexList = ref([]);
    const mappingStruct = ref({});
    const currentMapping = ref([]);

    // get the list of indices from server when the component is mounted
    const getIndexList = async function () {
      const response = await axios.get(store.state.API_ENDPOINT + "api/index");
      var data = response.data;
      indexStruct.value = data;
      for (var index in indexStruct.value) {
        indexList.value.push({ index: indexStruct.value[index].name });
        var mappingArray = [];
        for (var field in indexStruct.value[index].mappings) {
          mappingArray.push({
            field: field,
            type: indexStruct.value[index].mappings[field],
          });
        }
        mappingStruct.value[index] = mappingArray;
      }
    };

    getIndexList();

    return {
      // variables
      indexStruct,
      indexList,
      mappingStruct,
      currentMapping,
      pagination_indexes: {
        rowsPerPage: 10, // current rows per page being displayed in the fields section
      },
      pagination_fields: {
        rowsPerPage: 20, // current rows per page being displayed in the fields section
      },

      // methods
      indexNameClicked: function (evt, row) {
        currentMapping.value = mappingStruct.value[row.index];
        //
      },
    };
  },
};
</script>

<style scoped>
.indexmanagement {
  display: flex;
  flex-direction: column;
  align-content: center;
}

.index-table {
  margin: 10px;
}

.mapping-table {
  margin: 10px;
}
</style>
