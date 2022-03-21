<template>
  <q-layout view="hHh lpR fFf">
    <q-header>
      <q-toolbar>
        <q-btn
          flat
          dense
          round
          icon="menu"
          aria-label="Menu"
          @click="leftDrawerOpen = !leftDrawerOpen"
        />

        <q-toolbar-title>Zinc Search</q-toolbar-title>

        <div>
          <q-btn-dropdown outline rounded no-caps icon-right="manage_accounts">
            <template #label>
              <div class="row items-center no-wrap">{{ user.name }}</div>
            </template>
            <q-list>
              <q-item-label header>Account</q-item-label>

              <q-item v-ripple v-close-popup clickable @click="signout">
                <q-item-section avatar>
                  <q-avatar
                    size="md"
                    icon="exit_to_app"
                    color="red"
                    text-color="white"
                  />
                </q-item-section>
                <q-item-section>
                  <q-item-label>Sign Out</q-item-label>
                </q-item-section>
              </q-item>
            </q-list>
          </q-btn-dropdown>
        </div>
      </q-toolbar>
    </q-header>

    <q-drawer
      v-model="leftDrawerOpen"
      :width="200"
      :breakpoint="500"
      bordered
      class="bg-grey-3"
    >
      <q-list>
        <menu-link
          v-for="link in essentialLinks"
          :key="link.title"
          v-bind="link"
        />
      </q-list>
    </q-drawer>

    <q-page-container>
      <router-view v-slot="{ Component }">
        <keep-alive>
          <component
            :is="Component"
            v-if="$route.meta.keepAlive"
            :key="$route.name"
          />
        </keep-alive>
        <component
          :is="Component"
          v-if="!$route.meta.keepAlive"
          :key="$route.name"
        />
      </router-view>
    </q-page-container>
  </q-layout>
</template>

<script>
import MenuLink from "../components/MenuLink.vue";

const linksList = [
  {
    title: "Search",
    icon: "manage_search",
    link: "/search",
  },
  {
    title: "Index",
    icon: "list",
    link: "/index",
  },
  {
    title: "Template",
    icon: "apps",
    link: "/template",
  },
  {
    title: "User",
    icon: "people",
    link: "/user",
  },
  {
    title: "About",
    icon: "info",
    link: "/about",
  },
];

import { ref } from "vue";
import { useStore } from "vuex";
import { useRouter } from "vue-router";

export default {
  name: "MainLayout",

  components: {
    MenuLink,
  },

  setup() {
    const store = useStore();
    const router = useRouter();

    const signout = () => {
      store.dispatch("logout");
      localStorage.setItem("creds", "");
      router.push("/login");
    };

    return {
      essentialLinks: linksList,
      leftDrawerOpen: ref(false),
      user: store.state.user,
      signout,
    };
  },
};
</script>

<style lang="scss">
@import "../styles/app.scss";
</style>
