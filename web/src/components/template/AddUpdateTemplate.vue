<template>
  <q-card class="column full-height">
    <q-card-section>
      <div class="row items-center no-wrap">
        <div class="col">
          <div class="text-h6">Add template</div>
        </div>

        <div class="col-auto">
          <q-btn v-close-popup flat round color="grey-7" icon="close" />
        </div>
      </div>
    </q-card-section>

    <q-card-section class="col q-pt-none q-w-p50">
      <q-stepper ref="stepper" v-model="step" color="primary" flat animated>
        <q-step :name="1" title="Logistics" icon="info" :done="step > 1">
          <q-form ref="step1Form" @submit="onSubmitStep1">
            <q-input
              v-model="templateData.name"
              borderless
              filled
              :readonly="beingUpdated"
              :disabled="beingUpdated"
              :bg-color="disableColor"
              :rules="[validateTemplateName]"
              label="Template Name"
              class="q-py-md"
            />
            <q-select
              v-model="templateData.index_patterns"
              label="Index Patterns"
              placeholder="Type and then hit `ENTER`"
              filled
              use-input
              use-chips
              multiple
              hide-dropdown-icon
              new-value-mode="add-unique"
              :rules="[validateTemplateIndexPatterns]"
              class="q-py-md"
            />
            <q-input
              v-model="templateData.priority"
              borderless
              filled
              label="Priority (optional)"
              :rules="[validateTemplatePriority]"
              class="q-py-md"
            />
          </q-form>
        </q-step>

        <q-step
          :name="2"
          title="Index settings"
          icon="settings"
          :done="step > 2"
        >
          <q-form ref="step2Form">
            <json-editor
              v-model="templateData.template.settings"
              name="settings"
              :height="480"
              @validation-error="onJsonError"
            ></json-editor>
          </q-form>
          <p class="text-grey q-mb-none">
            Use JSON format:
            <strong
              class="bg-grey-2 text-purple-6 q-px-sm"
              style="font-weight: normal"
              >{{ placeholderIndexSettings }}</strong
            >
          </p>
        </q-step>

        <q-step :name="3" title="Mappings" icon="assignment" :done="step > 3">
          <q-form ref="step3Form">
            <json-editor
              v-model="templateData.template.mappings"
              name="mappings"
              :height="480"
              @validation-error="onJsonError"
            ></json-editor>
          </q-form>
          <p class="text-grey q-mb-none">
            Use JSON format:
            <strong
              class="bg-grey-2 text-purple-6 q-px-sm"
              style="font-weight: normal"
              >{{ placeholderMappings }}</strong
            >
          </p>
        </q-step>

        <q-step :name="4" title="Review" icon="preview">
          <q-form ref="step4Form">
            <json-editor
              v-model="templateData"
              name="preview"
              mode="code"
              :readonly="true"
              :height="501"
            ></json-editor>
          </q-form>
        </q-step>

        <template #navigation>
          <q-stepper-navigation>
            <q-btn
              no-caps
              color="primary"
              :disable="disableBtn"
              :label="step === 4 ? 'Save Template' : 'Continue'"
              @click="nextStep"
            />
            <q-btn
              v-if="step > 1"
              flat
              no-caps
              color="primary"
              :disable="disableBtn"
              label="Back"
              class="q-ml-sm"
              @click="$refs.stepper.previous()"
            />
          </q-stepper-navigation>
        </template>
      </q-stepper>
    </q-card-section>
  </q-card>
</template>

<script>
import { defineComponent, ref } from "vue";
import templateService from "../../services/template";
import JsonEditor from "../JsonEditor.vue";

const defaultValue = () => ({
  name: "",
  index_patterns: [],
  priority: "",
  template: {
    settings: {},
    mappings: {},
  },
});

export default defineComponent({
  name: "ComponentAddUpdateTemplate",
  components: {
    JsonEditor,
  },
  props: {
    modelValue: {
      type: Object,
      default: () => defaultValue(),
    },
  },
  emits: ["update:modelValue", "updated"],
  setup() {
    const beingUpdated = ref(false);
    const step1Form = ref(null);
    const disableColor = ref("");
    const disableBtn = ref(false);
    const templateData = ref(defaultValue());

    return {
      step: ref(1),
      step1Form,
      disableColor,
      disableBtn,
      beingUpdated,
      templateData,
      placeholderIndexSettings: ref(`{
  "analysis": {
    "analyzer": {
      "default": {
        "type": "standard"
      }
    }
  }
}`),
      placeholderMappings: ref(`{
  "properties": {
    "content": {
      "type": "text"
    }
  }
}`),
    };
  },
  created() {
    if (this.modelValue && this.modelValue.name) {
      this.beingUpdated = true;
      this.disableColor = "grey-5";
      for (const indexPattern of this.modelValue.patterns.split(", ")) {
        this.templateData.index_patterns.push(indexPattern);
      }
      this.templateData.name = this.modelValue.name;
      this.templateData.priority = this.modelValue.priority;
      this.templateData.template.settings =
        this.modelValue.template.settings || {};
      this.templateData.template.mappings =
        this.modelValue.template.mappings || {};
    }
  },
  methods: {
    validateTemplateName(val) {
      if (val && val.length > 0) {
        this.disableBtn = false;
        return true;
      }
      return "Template name is required";
    },
    validateTemplateIndexPatterns(val) {
      if (val && val.length > 0) {
        this.disableBtn = false;
        return true;
      }
      return "Index patterns are required";
    },
    validateTemplatePriority(val) {
      if (val && val.length > 0) {
        this.templateData.priority = parseInt(val, 10);
      }
      return true;
    },
    onJsonError(error) {
      if (error && error.length > 0) {
        this.disableBtn = true;
      } else {
        this.disableBtn = false;
      }
    },
    nextStep() {
      if (this.step === 1) {
        this.onSubmitStep1();
        return;
      }

      if (this.step < 4) {
        this.$refs.stepper.next();
        return;
      }

      // save template
      this.onSave();
    },
    onSubmitStep1() {
      this.step1Form.validate().then((valid) => {
        if (!valid) {
          this.disableBtn = true;
        } else {
          this.$refs.stepper.next();
          this.step1Form.resetValidation();
        }
      });
    },
    onSave() {
      templateService.update(this.templateData).then((res) => {
        this.$emit("update:modelValue", this.templateData);
        this.$emit("updated", this.templateData);
        this.templateData = defaultValue();
      });
    },
  },
});
</script>
