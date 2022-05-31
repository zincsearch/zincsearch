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

        <q-toolbar-title>{{ t("menu.zincSearch") }}</q-toolbar-title>

        <div class="q-mr-xs">
          <q-btn-dropdown outline rounded no-caps icon-right="manage_accounts">
            <template #label>
              <div class="row items-center no-wrap">{{ user.name }}</div>
            </template>
            <q-list>
              <q-item-label header>{{ t("menu.account") }}</q-item-label>

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
                  <q-item-label>{{ t("menu.signOut") }}</q-item-label>
                </q-item-section>
              </q-item>
            </q-list>
          </q-btn-dropdown>
        </div>

        <div>
          <q-btn-dropdown outline rounded no-caps icon-right="language">
            <q-list>
              <q-item v-ripple v-close-popup clickable @click="changeLanguage('en')">
                <q-item-section>
                  <q-item-label>English</q-item-label>
                </q-item-section>
              </q-item>
              <q-item v-ripple v-close-popup clickable @click="changeLanguage('zh-cn')">
                <q-item-section>
                  <q-item-label>简体中文</q-item-label>
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

<script lang="ts">
import MenuLink from "../components/MenuLink.vue";
import { useI18n } from "vue-i18n";
import { setLanguage } from "../utils/cookies";


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
    const { t } = useI18n();

    const linksList = [
      {
        title: t("menu.search"),
        icon: "manage_search",
        link: "/search",
      },
      {
        title: t("menu.index"),
        icon: "list",
        link: "/index",
      },
      {
        title:  t("menu.template"),
        icon: "apps",
        link: "/template",
      },
      {
        title:  t("menu.user"),
        icon: "people",
        link: "/user",
      },
      {
        title:  t("menu.about"),
        icon: "info",
        link: "/about",
      },
    ];
    const changeLanguage = (language: string) => {
      setLanguage(language);
      router.go(0);
    };
    const signout = () => {
      store.dispatch("logout");
      localStorage.setItem("creds", "");
      router.push("/login");
    };

    return {
      t,
      essentialLinks: linksList,
      leftDrawerOpen: ref(false),
      user: store.state.user,
      changeLanguage,
      signout,
    };
  },
};
</script>

<style lang="scss">
@import "../styles/app.scss";
</style>
