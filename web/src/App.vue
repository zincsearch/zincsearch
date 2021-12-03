<template>
  <q-layout view="hHh lpR fFf">
    <q-header bordered class="bg-primary text-white" v-show="isLoggedIn">
      <q-toolbar>
        <q-btn
          flat
          dense
          round
          @click="leftDrawerOpen = !leftDrawerOpen"
          aria-label="Menu"
          icon="menu"
        />

        <q-toolbar-title> Zinc Search </q-toolbar-title>

        <!-- <div>Quasar v{{ $q.version }}</div> -->
        <q-btn-dropdown outline rounded icon-right="manage_accounts">
          <template v-slot:label>
            <div class="row items-center no-wrap">
              {{ loggedInUserName }}
              <!-- <q-avatar size="40px">
                {{ loggedInUserName }}
              </q-avatar> -->
            </div>
          </template>

          <q-item-label header>Account</q-item-label>

          <q-item clickable v-close-popup @click="signout">
            <q-item-section avatar>
              <q-avatar icon="exit_to_app" color="red" text-color="white" />
            </q-item-section>
            <q-item-section>
              <q-item-label>Sign Out </q-item-label>
              <!-- <q-item-label style="font-size: 0.9rem">{{
                loggedInUser
              }}</q-item-label> -->
            </q-item-section>
          </q-item>
        </q-btn-dropdown>
      </q-toolbar>
    </q-header>

    <q-drawer
      v-model="leftDrawerOpen"
      bordered
      class="bg-grey-2"
      v-show="isLoggedIn"
    >
      <q-list>
        <!-- <q-item-label header>Essential Links</q-item-label> -->

        <q-item clickable to="/">
          <q-item-section avatar>
            <q-icon name="manage_search" />
          </q-item-section>
          <q-item-section>
            <q-item-label>Search</q-item-label>
            <q-item-label caption>search</q-item-label>
          </q-item-section>
        </q-item>

        <q-item clickable to="/indexmanagement">
          <q-item-section avatar>
            <q-icon name="list" />
          </q-item-section>
          <q-item-section>
            <q-item-label>Index Management</q-item-label>
            <q-item-label caption>Index management</q-item-label>
          </q-item-section>
        </q-item>

        <q-item clickable to="/users">
          <q-item-section avatar>
            <q-icon name="people" />
          </q-item-section>
          <q-item-section>
            <q-item-label>Users</q-item-label>
            <q-item-label caption>Users</q-item-label>
          </q-item-section>
        </q-item>

        <q-item clickable to="/about">
          <q-item-section avatar>
            <q-icon name="info" />
          </q-item-section>
          <q-item-section>
            <q-item-label>About</q-item-label>
            <q-item-label caption>about</q-item-label>
          </q-item-section>
        </q-item>
      </q-list>
    </q-drawer>

    <q-page-container>
      <router-view></router-view>
    </q-page-container>
  </q-layout>
</template>

<script>
import { ref } from "vue";
import { mapState } from "vuex";
import { useStore } from "vuex";
import router from "./router";
// import HelloWorld from "./components/HelloWorld.vue";

export default {
  name: "LayoutDefault",

  components: {
    // HelloWorld,
  },
  computed: mapState({
    isLoggedIn: (state) => state.user.isLoggedIn,
    loggedInUser: (state) => state.user._id,
    loggedInUserName: (state) => state.user.name,
  }),

  setup() {
    // const user = ref({});
    const store = useStore();

    const _id = localStorage.getItem("_id");
    const base64encoded = localStorage.getItem("base64encoded");
    const name = localStorage.getItem("name");
    const role = localStorage.getItem("role");

    if (_id && base64encoded) {
      store.dispatch("login", {
        _id: localStorage.getItem("_id"),
        base64encoded: localStorage.getItem("base64encoded"),
        name: localStorage.getItem("name"),
        role: localStorage.getItem("role"),
      });
      // user.value = store.state.user;
      router.push("/");
    }
    return {
      leftDrawerOpen: ref(false),
      name,
      role,
      // user,

      signout() {
        store.dispatch("logout");
        localStorage.setItem("_id", "");
        localStorage.setItem("base64encoded", "");
        localStorage.setItem("name", "");
        localStorage.setItem("role", "");
        router.push("/login");
      },
    };
  },
};
</script>
