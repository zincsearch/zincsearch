<template>
  <q-card>
    <q-card-section>
      <div v-if="beingUpdated" class="text-h6">Update user</div>
      <div v-else class="text-h6">Add user</div>
    </q-card-section>
    <q-card-section class="q-w-md">
      <q-form ref="addUserForm" @submit="onSubmit">
        <q-input
          v-model="userData._id"
          dense
          borderless
          filled
          :readonly="beingUpdated"
          :disabled="beingUpdated"
          :bg-color="disableColor"
          :label="t('user.id')"
          :rules="[validateUserID]"
        />
        <q-input
          v-model="userData.name"
          dense
          borderless
          filled
          :label="t('user.name')"
          :rules="[validateUserName]"
        />
        <q-select
          v-model="userData.role"
          :options="roles"
          dense
          borderless
          filled
          emit-value
          :label="t('user.role')"
          :rules="[validateUserRole]"
        />
        <q-input
          v-model="userData.password"
          borderless
          dense
          filled
          :type="isPwd ? 'password' : 'text'"
          :label="t('user.password')"
          :rules="[validatePassword]"
        >
          <template #append>
            <q-icon
              :name="isPwd ? 'visibility_off' : 'visibility'"
              class="cursor-pointer"
              @click="isPwd = !isPwd"
            />
          </template>
        </q-input>
        <q-input
          v-model="userData.confirmPassword"
          borderless
          dense
          filled
          :type="isPwd ? 'password' : 'text'"
          :label="t('user.repassword')"
          :rules="[validateConfirmPassword]"
        >
          <template #append>
            <q-icon
              :name="isPwd ? 'visibility_off' : 'visibility'"
              class="cursor-pointer"
              @click="isPwd = !isPwd"
            />
          </template>
        </q-input>

        <q-btn
          no-caps
          class="q-mb-md"
          color="primary"
          type="submit"
          icon="add"
          label="Save User"
        />
      </q-form>
    </q-card-section>
  </q-card>
</template>

<script>
import { defineComponent, ref } from "vue";
import userService from "../../services/user";
import roleService from "../../services/role";
import { useI18n } from "vue-i18n";

const defaultValue = () => {
  return {
    _id: "",
    name: "",
    role: "",
    password: "",
    confirmPassword: "",
  };
};

export default defineComponent({
  name: "ComponentAddUpdateUser",
  props: {
    modelValue: {
      type: Object,
      default: () => defaultValue(),
    },
  },
  emits: ["update:modelValue", "updated"],
  setup() {
    const beingUpdated = ref(false);
    const addUserForm = ref(null);
    const disableColor = ref("");
    const userData = ref(defaultValue());
    const { t } = useI18n();
    const roles = ref([]);
    const getRoles = () => {
      roleService.list().then((res) => {
        roles.value = res.data.map((data) => {
          return {
            value: data._id,
            label: data.name,
          };
        });
        roles.value.push({ value: "admin", label: "admin" });
      });
    };

    getRoles();

    return {
      t,
      disableColor,
      isPwd: ref(true),
      beingUpdated,
      roles,
      userData,
      addUserForm,
    };
  },
  created() {
    if (this.modelValue && this.modelValue.id) {
      this.beingUpdated = true;
      this.disableColor = "grey-5";
      this.userData = {
        _id: this.modelValue.id,
        name: this.modelValue.name,
        role: this.modelValue.role,
        password: "",
        confirmPassword: "",
      };
    }
  },
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
      if (data.length < 8) {
        return "Your password must be at least 8 characters";
      }
      if (data.search(/[a-z]/i) < 0) {
        return "Your password must contain at least one letter.";
      }
      if (data.search(/[0-9]/) < 0) {
        return "Your password must contain at least one digit.";
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
        if (!valid) {
          // console.log("Form is invalid");
          return false;
        }
        // console.log("Form is valid");
        userService.update(this.userData).then((res) => {
          var data = res.data;
          this.userData = {
            _id: "",
            name: "",
            password: "",
            confirmPassword: "",
            role: "",
          };

          this.$emit("update:modelValue", data);
          this.$emit("updated", data);
          this.addUserForm.resetValidation();
        });
      });
    },
  },
});
</script>
