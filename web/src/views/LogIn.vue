<template>
  <q-form @submit="signIn">
    <q-page
      class="window-height window-width row justify-center items-center"
      style="background: gray"
    >
      <div class="column q-pa-lg">
        <div class="row">
          <q-card square class="shadow-24" style="width: 300px; height: 485px">
            <q-card-section class="bg-indigo-5">
              <h4 class="text-h5 text-white q-my-md">Zinc Search</h4>
              <div
                class="absolute-bottom-right q-pr-md"
                style="transform: translateY(50%)"
              ></div>
            </q-card-section>
            <q-card-section>
              <!-- <q-form class="q-px-sm q-pt-xl"> -->
              <q-input square v-model="id" type="text" label="User ID">
                <template v-slot:prepend>
                  <q-icon name="email" />
                </template>
              </q-input>
              <q-input
                square
                v-model="password"
                type="password"
                label="Password"
              >
                <template v-slot:prepend>
                  <q-icon name="lock" />
                </template>
              </q-input>
              <!-- </q-form> -->
            </q-card-section>
            <q-card-actions class="q-px-lg">
              <q-btn
                type="submit"
                unelevated
                size="lg"
                color="indigo-4"
                class="full-width text-white"
                label="Sign In"
              />
            </q-card-actions>
            <!-- <q-card-section class="text-center q-pa-sm">
              <p class="text-grey-6">Forgot your password?</p>
            </q-card-section> -->
          </q-card>
        </div>
      </div>
    </q-page>
  </q-form>
</template>

<script>
import { ref } from "vue";
import { useStore } from "vuex";
import axios from "../axios";
import router from "../router";
// import { mapMutations } from 'vuex'

export default {
  name: "Login",
  setup() {
    const store = useStore();
    const id = ref("");
    const password = ref("");

    const signIn = () => {
      const base64encoded = btoa(id.value + ":" + password.value);
      var creds = {
        _id: id.value,
        password: password.value,
        base64encoded: base64encoded,
      };

      axios.post(store.state.API_ENDPOINT + "api/login", creds).then((res) => {
        if (res.data.validated) {
          creds.name = res.data.user.name;
          creds.role = res.data.user.role;

          store.dispatch("login", creds);

          localStorage.setItem("base64encoded", base64encoded);
          localStorage.setItem("_id", creds._id);
          localStorage.setItem("name", creds.name);
          localStorage.setItem("role", creds.role);

          router.push({ path: "/search" });
        } else {
          alert("Invalid credentials");
          store.dispatch("logout");
        }
      });
    };

    return {
      id,
      password,
      signIn,
    };
  },
};
</script>

<style></style>
