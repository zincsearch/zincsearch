<template>
  <SearchComponent />
</template>

<script>
// @ is an alias to /src
// import HelloWorld from "@/components/HelloWorld.vue";
import SearchComponent from "@/components/SearchComponent.vue";
// import SearchBar from '../components/SearchBar.vue';

export default {
  name: "Search",
  components: {
    SearchComponent,
  },
  data: function () {
    return {
      columns: [
        {
          name: "field",
          required: true,
          label: "Fields",
          align: "left",
          field: (row) => row.name,
          format: (val) => `${val}`,
          sortable: true,
        },
      ],
      rows: [],
    };
  },
};
</script>

<style scoped>
.field-list {
  width: 90%;
  margin: 10px;
}

.result-list {
  width: 98%;
  margin-top: 10px;
}
</style>
