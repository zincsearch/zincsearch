<template>
  <div>
    <q-dialog v-model="showUserAddDialog">
      <q-card>
        <q-card-section>
          <div class="text-h6">
            {{ t("user.addOrUpdate") }}
          </div>
        </q-card-section>

        <q-card-section class="add-user-dialog">
          <AddUpdateUserComponent @userAdded="userAdded" />
        </q-card-section>
      </q-card>
    </q-dialog>

    <q-dialog v-model="showUserUpdateDialog">
      <q-card>
        <q-card-section>
          <div class="text-h6">
            {{ t("user.addOrUpdate") }}
          </div>
        </q-card-section>

        <q-card-section class="add-user-dialog">
          <AddUpdateUserComponent
            @userUpdated="userUpdated"
            v-bind:user="user"
          />
        </q-card-section>
      </q-card>
    </q-dialog>

    <q-table
      dense
      :rows="users"
      row-key="UserID"
      class="users-table"
      :title="t('user.header')"
      :pagination="pagination"
      :filter="filter_query"
      :filter-method="filter_method"
    >
      <template v-slot:top-right>
        <q-input
          filled
          borderless
          dense
          debounce="1"
          v-model="filter_query"
          :placeholder="t('user.search')"
          class="search-user"
        >
          <template v-slot:append>
            <q-icon name="search" />
          </template>
        </q-input>
        <q-btn class="add-button" color="secondary" @click="addUser">
          <q-icon name="add" />
          {{ t("user.add") }}
        </q-btn>
      </template>

      <template v-slot:body="props">
        <q-tr :props="props">
          <q-td v-for="col in props.cols" :key="col.name" :props="props">
            {{ col.value }}
          </q-td>
          <q-td auto-width>
            <q-btn
              size="sm"
              color="secondary"
              class="action-button"
              dense
              @click="startEditing(props)"
              :icon="'edit'"
            />
            <q-btn
              size="sm"
              color="negative"
              class="action-button"
              dense
              @click="confirmDelete(props.key)"
              :icon="'delete'"
            />
          </q-td>
        </q-tr>
      </template>
    </q-table>
  </div>
</template>

<script>
import { ref } from "@vue/reactivity";
import { useQuasar } from "quasar";
import axios from "../axios";
import store from "../store";
import AddUpdateUserComponent from "./AddUpdateUserComponent";
import { useI18n } from "vue-i18n";

export default {
  components: {
    AddUpdateUserComponent,
  },
  created() {},
  methods: {},
  setup() {
    const $q = useQuasar();
    const users = ref([]);
    const columns = ref([]);
    const filter_query = ref("");
    const showUserAddDialog = ref(false);
    const showUserUpdateDialog = ref(false);
    const user = ref({});
    const { t } = useI18n();

    function startEditing(u) {
      console.log("startEditing", u);
      user.value = {
        id: u.row.id,
        name: u.row.name,
        role: u.row.role,
      };

      showUserUpdateDialog.value = true;
    }

    function confirmDelete(id) {
      $q.dialog({
        title: "Confirm User Delete",
        message:
          "Do you want to delete user " +
          id +
          "?" +
          " This action cannot be undone.",
        cancel: true,
        persistent: true,
      })
        .onOk(() => {
          // console.log('>>>> OK')
          axios.delete(store.state.API_ENDPOINT + "api/user/" + id).then(() => {
            getUsers();
          });
        })
        .onOk(() => {
          // console.log('>>>> second OK catcher')
        })
        .onCancel(() => {
          // console.log('>>>> Cancel')
        })
        .onDismiss(() => {
          // console.log('I am triggered on both OK and Cancel')
        });
    }

    function getUsers() {
      axios.get(store.state.API_ENDPOINT + "api/users").then((response) => {
        var counter = 1;
        users.value = response.data.hits.hits.map((data) => {
          return {
            "#": counter++,
            id: data._source._id,
            name: data._source.name,
            role: data._source.role,
            Created: data._source.created_at,
            Updated: data._source["@timestamp"],
          };
        });
      });
    }

    getUsers();

    function userAdded() {
      showUserAddDialog.value = false;
      getUsers();
    }

    function userUpdated() {
      showUserUpdateDialog.value = false;
      getUsers();
    }

    function deleteUser(id) {
      axios.delete(store.state.API_ENDPOINT + "api/user/" + id).then(() => {
        getUsers();
      });
    }

    return {
      user,
      showUserAddDialog,
      showUserUpdateDialog,
      users,
      pagination: {
        rowsPerPage: 20, // current rows per page being displayed
      },
      columns,
      filter_query,
      confirmDelete,
      startEditing,
      getUsers,
      userAdded,
      userUpdated,
      deleteUser,

      // filter for username
      filter_method(rows, terms) {
        terms = terms.toLowerCase();
        var filtered_rows = [];
        for (var i = 0; i < rows.length; i++) {
          if (rows[i]["name"].toLowerCase().includes(terms)) {
            filtered_rows.push(rows[i]);
          }
        }

        return filtered_rows;
      },
      addUser() {
        showUserAddDialog.value = true;
      },
      t,
    };
  },
};
</script>

<style scoped>
.users-table {
  margin: 10px;
}

.action-button {
  margin-left: 5px;
}

.search-user {
  margin-right: 5%;
  width: 50%;
}

.add-button {
  /* margin-right: 100px; */
  width: 40%;
}

.add-user-dialog {
  width: 400px;
}
</style>
