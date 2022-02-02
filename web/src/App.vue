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
      <q-page-sticky position="bottom-right" :offset="[30, 80]">
        <div class="feedback">
          <q-btn
            color="green"
            label="Feedback"
            class="feedback-button"
            @click="showFeedback"
          />
        </div>
        <!-- <q-btn fab icon="add" color="accent" /> -->
      </q-page-sticky>
    </q-page-container>
  </q-layout>
</template>

<script>
import { ref } from "vue";
import { mapState } from "vuex";
import { useStore } from "vuex";
import { useQuasar } from "quasar";
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
    const $q = useQuasar();

    if (window.location.origin != "http://localhost:8080") {
      store.dispatch("endpoint", window.location.origin + "/");
    }

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

      showFeedback() {
        $q.dialog({
          title: "",
          message:
            '<iframe src="https://docs.google.com/forms/d/e/1FAIpQLSeYvo97TMoJ6qWZKXLndQbimxcP0sCdHBrUOeg8D785rHVv1g/viewform?embedded=true" width="840" height="1822" frameborder="0" marginheight="0" marginwidth="0">Loading…</iframe>',
          html: true,
          fullWidth: true,
        })
          .onOk(() => {
            // console.log('OK')
          })
          .onCancel(() => {
            // console.log('Cancel')
          })
          .onDismiss(() => {
            // console.log('I am triggered on both OK and Cancel')
          });
      },

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

<style scoped>
.feedback {
  /* transform: rotate(-90deg); */
  margin-right: 0%;
  position: relative;
}

.feedback-button {
  border-radius: 50px;
}
</style>
