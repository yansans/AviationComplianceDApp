<template>
    <div>
      <h2>Create Asset</h2>
      <form @submit.prevent="submitAsset">
        <div>
          <label for="id">Asset ID:</label>
          <input type="text" id="id" v-model="formData.id" required />
        </div>
        <div>
          <label for="aircraftId">Aircraft ID:</label>
          <input type="text" id="aircraftId" v-model="formData.aircraftId" required />
        </div>
        <div>
          <label for="reportDate">Report Date:</label>
          <input type="date" id="reportDate" v-model="formData.reportDate" required />
        </div>
        <div>
          <label for="inspector">Inspector:</label>
          <input type="text" id="inspector" v-model="formData.inspector" required />
        </div>
        <div>
          <label for="description">Description:</label>
          <input type="text" id="description" v-model="formData.description" required />
        </div>
        <div>
          <label for="compliance">Compliance:</label>
          <select id="compliance" v-model="formData.compliance" required>
            <option value="true">Compliant</option>
            <option value="false">Non-Compliant</option>
          </select>
        </div>
        <button type="submit">Create New Asset</button>
      </form>
      <p v-if="responseMessage">{{ responseMessage }}</p>
    </div>
  </template>
  
  <script>
  import api from "../api/api";
  
  export default {
    data() {
      return {
        formData: {
          id: "",
          aircraft_id: "",
          report_date: "",
          inspector: "",
          description: "",
          compliance: true,
        },
        responseMessage: null,
      };
    },
    methods: {
      async submitAsset() {
        try {
          const response = await api.post("/create_asset", this.formData);
          this.responseMessage = `Asset created successfully: ${response.data.message}`;
        } catch (error) {
          console.error("Error creating asset:", error);
          this.responseMessage = "Failed to create asset.";
        }
      },
    },
  };
  </script>
  