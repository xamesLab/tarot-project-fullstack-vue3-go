import { createApp } from "vue";
import { createPinia } from "pinia";

import App from "./App.vue";
import { router, setupRouterGuards } from "./router";
import { useUiPreferencesStore } from "./stores/ui-preferences.store";

import "./assets/styles/variables.css";
import "./assets/styles/base.css";

const app = createApp(App);
const pinia = createPinia();

setupRouterGuards(pinia);

app.use(pinia);
useUiPreferencesStore(pinia).initialize();
app.use(router);
app.mount("#app");
