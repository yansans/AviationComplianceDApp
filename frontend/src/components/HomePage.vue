<template>
  <div>
    <h2>Welcome to BlogCain</h2>
    <p>{{ message }}</p>
    <button @click="fetchData">Fetch Data</button>
  </div>
</template>

<script>
import api from '../api/api';

export default {
  data() {
    return {
      message: "Try the button below to make a get request",
    };
  },
  methods: {
    async fetchData() {
      try {
        const response = await api.get("/info");
        this.message = response.data.message;
      } catch (error) {
        console.error("Error fetching data:", error);
        this.message = "Failed to fetch data from the API.";
      }
    },
  },
};
</script>
