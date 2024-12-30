<template>
    <div>
      <h2>Read Asset</h2>
      <form @submit.prevent="fetchAsset">
        <div>
          <label for="id">Asset ID:</label>
          <input type="text" id="id" v-model="id" required />
        </div>
        <button type="submit">Fetch Asset</button>
      </form>
      <p v-if="asset">{{ asset }}</p>
      <p v-if="responseMessage">{{ responseMessage }}</p>
    </div>
  </template>
  
  <script>
  import api from "../api/api";
  
  export default {
    data() {
      return {
        id: "",
        asset: null,
        responseMessage: null,
      };
    },
    methods: {
      async fetchAsset() {
        try {
          const response = await api.get(`/readAsset/${this.id}`);
          this.asset = JSON.stringify(response.data, null, 2);
        } catch (error) {
          console.error("Error fetching asset:", error);
          this.responseMessage = "Failed to fetch asset.";
        }
      },
    },
  };
  </script>
  