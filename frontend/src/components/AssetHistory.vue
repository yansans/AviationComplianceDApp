<template>
    <div>
      <h2>Asset History</h2>
      <form @submit.prevent="fetchHistory">
        <div>
          <label for="id">Asset ID:</label>
          <input type="text" id="id" v-model="id" required />
        </div>
        <button type="submit">Fetch Asset History</button>
      </form>
      <pre v-if="history">{{ history }}</pre>
      <p v-if="responseMessage">{{ responseMessage }}</p>
    </div>
  </template>
  
  <script>
  import api from "../api/api";
  
  export default {
    data() {
      return {
        id: "",
        history: null,
        responseMessage: null,
      };
    },
    methods: {
      async fetchHistory() {
        try {
          const response = await api.get(`/asset_history/${this.id}`);
          this.history = JSON.stringify(response.data, null, 2);
        } catch (error) {
          console.error("Error fetching history:", error);
          this.responseMessage = "Failed to fetch history.";
        }
      },
    },
  };
  </script>
  