<template>
  <router-view></router-view>
</template>

<script lang="ts">
import { useStore } from "vuex";
import { useRouter } from "vue-router";

export default {
  setup() {
    const store = useStore();
    if (window.location.origin != "http://localhost:8080") {
      store.dispatch(
        "endpoint",
        window.location.origin +
          window.location.pathname.split("/").slice(0, -2).join("/")
      );
    }

    const router = useRouter();
    const creds = localStorage.getItem("creds");
    if (creds) {
      const credsInfo = JSON.parse(creds);
      store.dispatch("login", credsInfo);
      router.push("/search");
    }
  },
};
</script>
