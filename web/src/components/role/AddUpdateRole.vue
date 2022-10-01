<template>
  <q-card class="my-card">
    <q-card-section>
      <div v-if="beingUpdated" class="text-h6">Update role</div>
      <div v-else class="text-h6">Add role</div>
    </q-card-section>
    <q-card-section>
      <q-form ref="addRoleForm" @submit="onSubmit">
        <q-input
          v-model="roleData._id"
          dense
          borderless
          filled
          :readonly="beingUpdated"
          :disabled="beingUpdated"
          :bg-color="disableColor"
          :label="t('role.id')"
          :rules="[validateRoleID]"
        />
        <q-input
          v-model="roleData.name"
          dense
          borderless
          filled
          :label="t('role.name')"
          :rules="[validateRoleName]"
        />
        <q-select
          class="q-field--with-bottom"
          v-model="roleData.permission"
          :options="permissions"
          dense
          multiple
          use-chips
          borderless
          filled
          :label="t('role.permission')"
        />

        <q-btn
          no-caps
          class="q-mb-md"
          color="primary"
          type="submit"
          icon="add"
          label="Save Role"
        />
      </q-form>
    </q-card-section>
  </q-card>
</template>

<script>
import { defineComponent, ref } from "vue";
import roleService from "../../services/role";
import permissionService from "../../services/permission";
import { useI18n } from "vue-i18n";

const defaultValue = () => {
  return {
    _id: "",
    name: "",
    role: "",
    permission: [],
  };
};

export default defineComponent({
  name: "ComponentAddUpdateRole",
  props: {
    modelValue: {
      type: Object,
      default: () => defaultValue(),
    },
  },
  emits: ["update:modelValue", "updated"],
  setup() {
    const beingUpdated = ref(false);
    const addRoleForm = ref(null);
    const disableColor = ref("");
    const roleData = ref(defaultValue());
    const { t } = useI18n();
    const permissions = ref([]);
    const getPermissions = () => {
      permissionService.list().then((res) => {
        permissions.value = res.data;
      });
    };

    getPermissions();

    return {
      t,
      disableColor,
      beingUpdated,
      permissions,
      roleData,
      addRoleForm,
    };
  },
  created() {
    if (this.modelValue && this.modelValue.id) {
      this.beingUpdated = true;
      this.disableColor = "grey-5";
      this.roleData = {
        _id: this.modelValue.id,
        name: this.modelValue.name,
        permission: this.modelValue.permission,
      };
    }
  },
  methods: {
    validateRoleID(data) {
      if (data.length < 3) {
        return "Role ID must be at least 3 characters long";
      }
    },
    validateRoleName(data) {
      if (data.length < 3) {
        return "Role name must be at least 3 characters long";
      }
    },
    onSubmit() {
      this.addRoleForm.validate().then((valid) => {
        if (!valid) {
          // console.log("Form is invalid");
          return false;
        }
        // console.log("Form is valid");
        roleService.update(this.roleData).then((res) => {
          var data = res.data;
          this.roleData = {
            _id: "",
            name: "",
            permission: [],
          };

          this.$emit("update:modelValue", data);
          this.$emit("updated", data);
          this.addRoleForm.resetValidation();
        });
      });
    },
  },
});
</script>

<style lang="scss" scoped>
.my-card {
  width: 100%;
  max-width: 800px;
}
</style>
