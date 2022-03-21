<template>
  <div>
    <div class="json-editor"></div>
    <p class="text-red-7 q-ma-none">{{ errorMessage }}</p>
  </div>
</template>

<script>
import JSONEditor from "jsoneditor";
import { ref } from "vue";

export default {
  props: {
    modelValue: {
      type: [String, Number, Object, Array],
      default: () => ({}),
    },
    mode: {
      type: String,
      default: "code",
    },
    name: {
      type: String,
      default: "jsonEditor",
    },
    height: {
      type: Number,
      default: 300,
    },
    readonly: {
      type: Boolean,
      default: false,
    },
  },
  emits: ["update:modelValue", "error", "validationError"],
  data() {
    return {
      errorMessage: ref(""),
      jsonEditor: null,
      internalChange: false,
    };
  },
  computed: {
    minHeight() {
      return this.height + "px";
    },
  },
  watch: {
    modelValue: {
      handler(newValue) {
        if (!this.internalChange) {
          this.setValue(newValue);
        }
      },
    },
  },
  mounted() {
    this.init();
  },
  unmounted() {
    this.jsonEditor?.destroy();
    this.jsonEditor = null;
  },
  methods: {
    init() {
      this.jsonEditor = new JSONEditor(
        this.$el.querySelector(".json-editor"),
        {
          name: this.name,
          mode: this.mode,
          indentation: 2,
          mainMenuBar: false,
          statusBar: false,
          onError: (err) => {
            this.$emit("error", err);
          },
          onValidationError: (err) => {
            if (err && err.length > 0) {
              this.errorMessage = "Invalid JSON format.";
            } else {
              this.errorMessage = "";
            }
            this.$emit("validationError", err);
          },
          onEditable: (node) => {
            if (this.readonly) {
              return false;
            }
            return true;
          },
          onChange: () => {
            try {
              this.internalChange = true;
              this.$emit("update:modelValue", this.jsonEditor.get());
              // prevent infinite loop
              this.$nextTick(() => {
                this.internalChange = false;
              });
            } catch (error) {}
          },
        },
        this.modelValue
      );
    },
    setValue(json) {
      if (!this.jsonEditor) {
        return false;
      }
      if (!json) {
        return false;
      }
      if (typeof json == "object") {
        this.jsonEditor.set(json);
      } else {
        this.jsonEditor.setText(json);
      }
    },
  },
};
</script>

<style lang="scss">
@import "jsoneditor/dist/jsoneditor.min.css";

.json-editor {
  height: 100%;
}
.jsoneditor {
  border: none;
  .ace_content {
    background: rgba(0, 0, 0, 0.05);
  }
}
.jsoneditor-preview {
  padding: 12px !important;
  background: rgba(0, 0, 0, 0.05);
  line-height: 1.5;
  min-height: v-bind(minHeight);
}
.ace-jsoneditor,
textarea.jsoneditor-text {
  min-height: v-bind(minHeight);
}
.ace-jsoneditor .ace_marker-layer .ace_active-line {
  background: rgba(0, 0, 0, 0.05);
}
.ace_gutter {
  display: none;
}
.ace_scroller {
  left: 0 !important;
}
</style>
