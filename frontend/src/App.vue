<template>
  <div id="app">
    <header>
      <h1>Aviation Compliance DApp</h1>
      <nav>
        <a><button class="btn-feature" @click="viewCreateAsset">Create Asset</button></a>
        <a><button class="btn-feature" @click="viewReadAsset">Read Asset</button></a>
        <a><button class="btn-feature" @click="viewUpdateCompliance">Update Compliance</button></a>
        <a><button class="btn-feature" @click="viewAssetHistory">View Asset History</button></a>
        <a><button class="btn-feature" @click="viewWalletSetting">Wallet Setting</button></a>
      </nav>
    </header>
    <component :is="currentView" />
  </div>
</template>

<script>
import CreateAsset from "./components/CreateAsset.vue";
import ReadAsset from "./components/ReadAsset.vue";
import UpdateCompliance from "./components/UpdateCompliance.vue";
import AssetHistory from "./components/AssetHistory.vue";
import WalletSetting from "./components/WalletSetting.vue";

export default {
  data() {
    return {
      message: "Select an action to interact with the ledger",
      currentView: null,
    };
  },

  components: {
    CreateAsset,
    ReadAsset,
    UpdateCompliance,
    AssetHistory,
    WalletSetting,
  },

  methods: {
    viewCreateAsset() {
      this.currentView = "CreateAsset";
    },
    viewReadAsset() {
      this.currentView = "ReadAsset";
    },
    viewUpdateCompliance() {
      this.currentView = "UpdateCompliance";
    },
    viewAssetHistory() {
      this.currentView = "AssetHistory";
    },
    viewWalletSetting() {
      this.currentView = "WalletSetting";
    },
  },
};
</script>

<style>
  @import './assets/styles.css';
</style>