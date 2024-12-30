<template>
    <div>
      <h2>Update Compliance</h2>
      <form @submit.prevent="updateCompliance">
        <div>
          <label for="id">Asset ID:</label>
          <input type="text" id="id" v-model="id" required />
        </div>
        <div>
          <label for="compliance">Compliance:</label>
          <select id="compliance" v-model="compliance" required>
            <option value="true">Compliant</option>
            <option value="false">Non-Compliant</option>
          </select>
        </div>
        <button type="submit">Update Compliance</button>
      </form>
      <p v-if="responseMessage">{{ responseMessage }}</p>
    </div>
  </template>
  
  <script>
  import api from "../api/api";
  
  export default {
    data() {
      return {
        id: "",
        compliance: "true",
        responseMessage: null,
      };
    },
    methods: {
      async updateCompliance() {
        try {
          const response = await api.post("/updateCompliance", { id: this.id, compliance: this.compliance });
          this.responseMessage = `Compliance updated for asset: ${response.data.id}`;
        } catch (error) {
          console.error("Error updating compliance:", error);
          this.responseMessage = "Failed to update compliance.";
        }
      },
    },
  };
  </script>
  