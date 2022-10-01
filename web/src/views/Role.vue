<template>
  <q-page class="q-pa-md">
    <q-table
      :title="t('role.header')"
      :rows="roles"
      :columns="columns"
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
          :placeholder="t('role.search')"
        >
          <template #append>
            <q-icon name="search" class="cursor-pointer" />
          </template>
        </q-input>
        <q-btn
          class="q-ml-sm"
          color="primary"
          icon="add"
          :label="t(`role.add`)"
          @click="addRole"
        />
      </template>

      <!-- eslint-disable-next-line vue/no-lone-template -->
      <template v-slot:body-cell-#="props">
        <q-td :props="props" width="80">
          {{ props.value }}
        </q-td>
      </template>
      <template #body-cell-actions="props">
        <q-td :props="props" auto-width>
          <q-btn
            dense
            unelevated
            size="sm"
            color="teal-5"
            class="action-button"
            icon="edit"
            @click="editRole(props)"
          />
          <q-btn
            dense
            unelevated
            size="sm"
            color="red-5"
            class="action-button q-ml-sm"
            icon="delete"
            @click="deleteRole(props)"
          />
        </q-td>
      </template>
    </q-table>

    <q-dialog v-model="showAddRoleDialog">
      <add-update-role @updated="roleAdded" />
    </q-dialog>

    <q-dialog v-model="showUpdateRoleDialog">
      <add-update-role v-model="role" @updated="roleUpdated" />
    </q-dialog>
  </q-page>
</template>

<script>
import { defineComponent, ref } from "vue";
import { useStore } from "vuex";
import { useQuasar, date } from "quasar";
import { useI18n } from "vue-i18n";

import roleService from "../services/role";
import AddUpdateRole from "../components/role/AddUpdateRole.vue";

export default defineComponent({
  name: "PageRole",
  components: {
    AddUpdateRole,
  },
  setup() {
    const store = useStore();
    const $q = useQuasar();
    const { t } = useI18n();

    const columns = [
      {
        name: "#",
        label: "#",
        field: "#",
        align: "left",
      },
      {
        name: "id",
        label: "ID",
        field: "id",
        align: "left",
      },
      {
        name: "name",
        label: "NAME",
        field: "name",
        align: "left",
      },
      {
        name: "created",
        label: "CREATED",
        field: "created",
        align: "left",
      },
      {
        name: "updated",
        label: "UPDATED",
        field: "created",
        align: "left",
      },
      {
        name: "actions",
        label: "ACTIONS",
        field: "actions",
        align: "left",
      },
    ];

    const role = ref({});
    const roles = ref([]);
    const getRoles = () => {
      roleService.list().then((res) => {
        var counter = 1;
        roles.value = res.data.map((data) => {
          return {
            "#": counter++,
            id: data._id,
            name: data.name || data._id,
            permission: data.permission,
            created: date.formatDate(data.created_at, "YYYY-MM-DDTHH:mm:ssZ"),
            updated: date.formatDate(data.updated_at, "YYYY-MM-DDTHH:mm:ssZ"),
            actions: "",
          };
        });
      });
    };

    getRoles();

    const showAddRoleDialog = ref(false);
    const showUpdateRoleDialog = ref(false);

    const addRole = () => {
      showAddRoleDialog.value = true;
    };
    const editRole = (props) => {
      role.value = {
        id: props.row.id,
        name: props.row.name,
        permission: props.row.permission,
      };
      showUpdateRoleDialog.value = true;
    };
    const deleteRole = (props) => {
      $q.dialog({
        title: "Delete Role",
        message:
          "You are about to delete this role: <ul><li>" +
          escape(props.row.id) +
          "</li></ul>",
        cancel: true,
        persistent: true,
        html: true,
      }).onOk(() => {
        roleService.delete(props.row.id).then(() => {
          getRoles();
        });
      });
    };

    return {
      t,
      columns,
      role,
      showAddRoleDialog,
      showUpdateRoleDialog,
      roles,
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
      addRole,
      editRole,
      deleteRole,
      roleAdded() {
        showAddRoleDialog.value = false;
        getRoles();
      },
      roleUpdated() {
        showUpdateRoleDialog.value = false;
        getRoles();
      },
    };
  },
});
</script>
