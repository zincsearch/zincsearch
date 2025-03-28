<template>
  <span v-for="item in list" :key="item.key" v-bind="item">
    <span v-if="item.isKeyWord" class="highlight">{{ item.text }}</span>
    <span v-else>{{ item.text }}</span>
  </span>
</template>

<script>
import { defineComponent, ref } from "vue";

export default defineComponent({
  name: "HighLight",
  props: {
    content: {
      type: String,
      required: true,
    },
    queryString: {
      type: String,
      default: "",
    },
  },
  data() {
    return {
      list: ref([]),
      keywords: ref([]),
    };
  },
  watch: {
    content: {
      handler() {
        this.init();
      },
    },
    queryString: {
      handler() {
        this.keywords = this.getKeywords(this.queryString);
        this.init();
      },
    },
  },
  mounted() {
    this.keywords = this.getKeywords(this.queryString);
    this.init();
  },
  methods: {
    init() {
      this.list = this.splitToList(this.content, this.keywords);
    },
    splitToList(content, keywords) {
      let arr = [
        {
          isKeyWord: false,
          text: content,
        },
      ];
      for (let i = 0; i < keywords.length; i++) {
        const keyword = keywords[i];
        let j = 0;
        while (j < arr.length) {
          let rec = arr[j];
          let record = rec.text.split(keyword);
          if (record.length > 1) {
            // delete j replace by new
            arr.splice(j, 1);
            let recKeyword = {
              isKeyWord: true,
              text: keyword,
            };
            for (let k = 0; k < record.length; k++) {
              let r = {
                isKeyWord: false,
                text: record[k],
              };
              if (k == record.length - 1) {
                arr.splice(j + k * 2, 0, r);
              } else {
                arr.splice(j + k * 2, 0, r, recKeyword);
              }
            }
          }
          j = j + record.length;
        }
      }
      return arr;
    },
    getKeywords(queryString) {
      if (!queryString || queryString.trim().length == 0) {
        return [];
      }

      let arr = [];
      // queryString + " " is for special split regular
      // split by space, but ignore double quotation marks
      const groups = (queryString + " ").split(/ s*(?![^"]*"\ )/);
      for (let i = 0; i < groups.length - 1; i++) {
        const group = groups[i];
        if (!group || group.trim().length == 0) {
          continue;
        }
        // group + ":" is for special split regular
        // split by :, but ignore "
        const fieldWordArr = (group + ":").split(/:s*(?![^"]*"\:)/);
        let keyword = group;
        if (fieldWordArr.length > 2) {
          keyword = fieldWordArr[1];
        }
        // delete start and end of * and "
        keyword = keyword
          .replace(/(^\**)|(\**$)/g, "")
          .replace(/(^"*)|("*$)/g, "");
        if (keyword.trim().length > 0) {
          // make sure key not empty or not space
          arr.push(keyword);
        }
      }

      let new_arr = [];
      for (let i = 0; i < arr.length; i++) {
        let start = 0;
        let in_word = false;
        for (let k = 0; k < arr[i].length; k++) {
          let c = arr[i].charCodeAt(k);
          // 0001 ~ 007e is ascii
          // ff60 ~ ff9f is half width katakana
          if ((c >= 0x0001 && c <= 0x007e) || (c >= 0xff60 && c <= 0xff9f)) {
            in_word = true;
          } else {
            if (in_word) {
              new_arr.push(arr[i].substring(start, k));
            }
            new_arr.push(arr[i].substring(k, k + 1));
            start = k + 1;
            in_word = false;
          }
        }
        if (in_word) {
          new_arr.push(arr[i].substring(start, arr[i].length));
        }
      }

      return new_arr;
    },
  },
});
</script>
<style lang="scss">
.highlight {
  background-color: rgb(255, 213, 0);
}
</style>
