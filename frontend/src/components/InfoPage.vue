<template>
  <div>
    <h2>About Us</h2>
    <p>This is a sample button for post request</p>

    <form @submit.prevent="submitData">
      <div>
        <label for="name">Name:</label>
        <input type="text" id="name" v-model="formData.name" required />
      </div>
      <div>
        <label for="email">Email:</label>
        <input type="email" id="email" v-model="formData.email" required />
      </div>
      <button type="submit">Submit</button>
    </form>

    <p v-if="responseMessage">{{ responseMessage }}</p>
  </div>
</template>

<script>
import api from '../api/api';

export default {
  data() {
    return {
      formData: {
        name: "",
        email: "",
      },
      responseMessage: null,
    };
  },
  methods: {
    async submitData() {
      try {
        const response = await api.post("/submit", this.formData);
        this.responseMessage = `Success: ${response.data.name} (${response.data.email})`;
      } catch (error) {
        console.error("Error submitting data:", error);
        this.responseMessage = "Failed to submit data.";
      }
    },
  },
};
</script>
