<template>
  <div class="add-user">
    <q-form @submit="onSubmit" ref="addUserForm">
      <div>
        <q-input
          dense
          borderless
          filled
          v-model="userData._id"
          :readonly="beingUpdated"
          :disabled="beingUpdated"
          :bg-color="uidbgColor"
          type="text"
          label="User ID"
          class="form-field"
          :rules="[validateUserID]"
        />
      </div>
      <div>
        <q-input
          dense
          borderless
          filled
          v-model="userData.name"
          type="text"
          label="User name"
          class="form-field"
          :rules="[validateUserName]"
        />
      </div>
      <div>
        <q-select
          :options="roles"
          dense
          borderless
          filled
          v-model="userData.role"
          label="Role"
          class="form-field"
          :rules="[validateUserRole]"
        />
      </div>
      <div>
        <q-input
          borderless
          dense
          v-model="userData.password"
          filled
          :type="isPwd ? 'password' : 'text'"
          label="Password"
          class="form-field"
          :rules="[validatePassword]"
        >
          <template v-slot:append>
            <q-icon
              :name="isPwd ? 'visibility_off' : 'visibility'"
              class="cursor-pointer"
              @click="isPwd = !isPwd"
            />
          </template>
        </q-input>
      </div>
      <div>
        <q-input
          borderless
          dense
          v-model="userData.confirmPassword"
          filled
          :type="isPwd ? 'password' : 'text'"
          label="Reconfirm Password"
          class="form-field"
          :rules="[validateConfirmPassword]"
        >
          <template v-slot:append>
            <q-icon
              :name="isPwd ? 'visibility_off' : 'visibility'"
              class="cursor-pointer"
              @click="isPwd = !isPwd"
            />
          </template>
        </q-input>
      </div>
      <div>
        <q-btn class="add-button form-field" color="secondary" type="submit">
          <q-icon name="add" />
          Add/Update User
        </q-btn>
      </div>
    </q-form>
  </div>
</template>

<script>
import { ref } from "vue";
import axios from "../axios";
// import { useStore } from "vuex";

export default {
  emits: ["userAdded", "userUpdated"],
  props: ["user"],
  methods: {
    validateUserID(data) {
      if (data.length < 3) {
        return "User ID must be at least 3 characters long";
      }
    },
    validateUserName(data) {
      if (data.length < 3) {
        return "User name must be at least 3 characters long";
      }
    },
    validateUserRole(data) {
      if (data.length < 3) {
        return "You must select a role";
      }
    },
    validatePassword(data) {
      if (this.beingUpdated && data.length == 0) {
        return true;
      }
      var errors = [];
      if (data.length < 8) {
        errors.push("Your password must be at least 8 characters");
      }
      if (data.search(/[a-z]/i) < 0) {
        errors.push("Your password must contain at least one letter.");
      }
      if (data.search(/[0-9]/) < 0) {
        errors.push("Your password must contain at least one digit.");
      }
      if (errors.length > 0) {
        return errors.join(". ");
      }
      return true;
    },

    validateConfirmPassword(data) {
      if (data !== this.userData.password) {
        return "Password and Confirmation password should match.";
      }
    },
    onSubmit() {
      this.addUserForm.validate().then((valid) => {
        if (valid) {
          console.log("Form is valid");
          axios
            .put(this.$store.state.API_ENDPOINT + "api/user", this.userData)
            .then((response) => {
              var data = response.data;
              this.userData = {
                _id: "",
                name: "",
                password: "",
                confirmPassword: "",
                role: "",
              };

              if (this.beingUpdated) {
                this.$emit("userUpdated", data);
              } else {
                this.$emit("userAdded", data);
              }
              this.$emit("userAdded", data);
              this.addUserForm.resetValidation();
            });
        } else {
          console.log("Form is invalid");
        }
      });
    },
  },
  data: function () {
    return {
      userInfo: this.user,
    };
  },
  created() {
    console.log("User Info", this.user);
    if (this.user && this.user.id) {
      this.beingUpdated = true;
      this.uidbgColor = "grey-5";
      this.userData = {
        _id: this.user.id,
        name: this.user.name,
        role: this.user.role,
        password: "",
        confirmPassword: "",
      };
    }
  },
  setup() {
    // const store = useStore();
    const beingUpdated = ref(false);
    const roles = ref(["admin", "user"]);
    const addUserForm = ref(null);
    const uidbgColor = ref("");
    const userData = ref({
      _id: "",
      name: "",
      password: "",
      confirmPassword: "",
      role: "",
    });

    return {
      uidbgColor,
      isPwd: ref(true),
      beingUpdated,
      roles,
      userData,
      // onSubmit,
      addUserForm,
    };
  },
};
</script>

<style scoped>
.add-user {
  margin: 10px;
}

.form-field {
  margin: 10px;
  /* width: 100%; */
}
</style>
