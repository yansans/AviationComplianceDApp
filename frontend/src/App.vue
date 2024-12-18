<template>
  <div id="app">
    <header>
      <h1>Blogcain Jaya</h1>
      <nav>
        <a href="#" @click="switchPage('home')">Home</a>
        <a href="#" @click="switchPage('about')">About</a>
      </nav>
    </header>

    <main>
      <HomePage v-if="currentPage === 'home'" />
      <InfoPage v-else />
    </main>
  </div>
</template>

<script>
import HomePage from './components/HomePage.vue';
import InfoPage from './components/InfoPage.vue';

export default {
  data() {
    return {
      currentPage: 'home',
    };
  },
  components: {
    HomePage,
    InfoPage,
  },
  methods: {
    switchPage(page) {
      this.currentPage = page;
    },
  },
};
</script>

<style>
/* Basic styles */
body {
  font-family: Arial, sans-serif;
  margin: 0;
  padding: 0;
}

header {
  background: #282c34;
  color: white;
  padding: 1rem;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

nav a {
  color: white;
  margin: 0 0.5rem;
  text-decoration: none;
  cursor: pointer;
}

main {
  padding: 2rem;
}
</style>
