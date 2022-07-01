<template>
  <q-layout view="hHh lpR fFf">
    <q-page-container>
      <q-page class="fullscreen bg-grey-7 flex flex-center">
        <q-card square class="my-card shadow-24 bg-white text-white">
          <q-card-section class="bg-primary">
            <div class="text-h5 q-my-md">Zinc Search</div>
          </q-card-section>
          <q-card-section class="bg-white">
            <q-form class="q-gutter-md" @submit="onSubmit">
              <q-input data-cy="login-user-id" v-model="id" label="User ID">
                <template #prepend>
                  <q-icon name="email" />
                </template>
              </q-input>

              <q-input data-cy="login-password" v-model="password" type="password" :label="t('login.password')">
                <template #prepend>
                  <q-icon name="lock" />
                </template>
              </q-input>

              <q-card-actions class="q-px-lg q-mt-md q-mb-xl">
                <q-btn
                  data-cy="login-sign-in"
                  unelevated
                  size="lg"
                  class="full-width"
                  color="primary"
                  type="submit"
                  :label="t('login.signIn')"
                  :loading="submitting"
                />
              </q-card-actions>
            </q-form>
          </q-card-section>
        </q-card>
      </q-page>
    </q-page-container>
  </q-layout>
</template>

<script>
import { defineComponent, ref } from "vue";
import { useStore } from "vuex";
import { useQuasar } from "quasar";
import { Buffer } from "buffer";
import { useRouter } from "vue-router";
import authapi from "../services/auth";
import { useI18n } from "vue-i18n";

export default defineComponent({
  name: "PageLogin",

  setup() {
    const store = useStore();
    const router = useRouter();
    const $q = useQuasar();
    const { t } = useI18n();
    const id = ref("");
    const password = ref("");
    const submitting = ref(false);

    const onSubmit = () => {
      if (id.value == "" || password.value == "") {
        $q.notify({
          position: "top",
          color: "warning",
          textColor: "white",
          icon: "warning",
          message: "Please input",
        });
      } else {
        submitting.value = true;
        let creds = {
          _id: id.value,
          password: password.value,
          base64encoded: Buffer.from(id.value + ":" + password.value).toString(
            "base64"
          ),
        };

        authapi.login(creds).then((res) => {
          if (res.data.validated) {
            creds.name = res.data.user.name;
            creds.role = res.data.user.role;

            localStorage.setItem("creds", JSON.stringify(creds));
            store.dispatch("login", creds);
            router.replace({ path: "/search" });
          } else {
            $q.notify({
              position: "bottom-right",
              progress: true,
              multiLine: true,
              color: "red-5",
              textColor: "white",
              icon: "warning",
              message: "Invalid credentials",
            });
            store.dispatch("logout");
            submitting.value = false;
          }
        });
      }
    };

    return {
      t,
      id,
      password,
      submitting,
      onSubmit,
    };
  },
});
</script>

<style lang="scss">
.my-card {
  width: 300px;
}
</style>
