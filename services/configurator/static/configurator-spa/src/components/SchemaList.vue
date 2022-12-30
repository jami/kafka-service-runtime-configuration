<script setup>
import SchemaListItem from "./SchemaListItem.vue";
</script>

<template>
  <div>
    <div>
      <SchemaListItem v-for="s in schemas" :key="s.id" :schema="s"></SchemaListItem>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      schemas: [],
      isLoading: false,
    };
  },
  methods: {
    getSchemas() {
      this.isLoading = true;

      fetch("/api/schema/list")
        .then((response) => response.json())
        .then((data) => {
          this.schemas = data.schemas;
          this.isLoading = false;
        });
    },
  },
  mounted() {
    this.getSchemas();
  },
};
</script>
