<template>
  <router-view></router-view>
</template>

<script lang="ts">
import { useStore } from "vuex";
import { useRouter } from "vue-router";
import user from "./services/user";

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
    user
      .isLoggedIn()
      .then((r) => {
        store.dispatch("login", {
          _id: r.data._id,
          password: r.data.password,
          name: r.data.name,
          role: r.data.role,
        });
        router.push("/search");
      })
      .catch((_) => {
        store.dispatch("logout");
      });
  },
};
</script>
