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
          <q-btn
            unelevated
            no-caps
            padding="xs sm"
            :label="t('menu.openapi')"
            href="/swagger/index.html"
            target="_blank"
          />
        </div>

        <div class="q-mr-xs">
          <q-btn
            unelevated
            no-caps
            padding="xs sm"
            :label="t('menu.documentation')"
            href="https://docs.zincsearch.com"
            target="_blank"
          />
        </div>

        <div>
          <q-btn-dropdown unelevated no-caps padding="xs sm">
            <template #label>
              <div class="row no-wrap">{{ selectedLanguage.label }}</div>
            </template>
            <q-list>
              <q-item
                v-for="lang in langList"
                :key="lang.value"
                v-ripple
                v-close-popup
                clickable
                v-bind="lang"
                @click="changeLanguage(lang)"
              >
                <q-item-section>
                  <q-item-label>{{ lang.label }}</q-item-label>
                </q-item-section>
              </q-item>
            </q-list>
          </q-btn-dropdown>
        </div>

        <div class="q-mr-xs">
          <q-btn-dropdown unelevated no-caps padding="xs sm">
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
import { getLocale } from "../locales";
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
        title: t("menu.template"),
        icon: "apps",
        link: "/template",
      },
      {
        title: t("menu.user"),
        icon: "people",
        link: "/user",
      },
      {
        title: t("menu.role"),
        icon: "manage_accounts",
        link: "/role",
      },
      {
        title: t("menu.about"),
        icon: "info",
        link: "/about",
      },
    ];

    const langList = [
      {
        label: "English",
        value: "en",
      },
      {
        label: "Türkçe",
        value: "tr",
      },
      {
        label: "简体中文",
        value: "zh-cn",
      },
    ];

    const local = ref(getLocale());
    const selectedLanguage = ref(langList.find((l) => l.value == local.value));

    if (!selectedLanguage.value && langList.length > 0) {
      selectedLanguage.value = langList[0];
    }

    const changeLanguage = (item: any) => {
      setLanguage(item.value);
      selectedLanguage.value = item;
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
      langList,
      selectedLanguage,
      changeLanguage,
      signout,
    };
  },
};
</script>

<style lang="scss">
@import "../styles/app.scss";

.q-header .q-btn-dropdown__arrow {
  margin-left: -4px;
}
</style>
