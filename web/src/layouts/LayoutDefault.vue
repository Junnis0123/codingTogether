<template>
  <div>
    <v-app-bar
      app
      color="primary"
      dark >
      <v-toolbar-title>모각코</v-toolbar-title>
    </v-app-bar>
    <v-main>
      <div class="d-flex align-center">
        <router-link to="/About">about</router-link>
        <router-link to="/">login</router-link>
      </div>
      <router-view />
    </v-main>
  </div>
</template>

<script lang="ts">
import { defineComponent } from '@vue/composition-api';

export default defineComponent({
  name: 'LayoutDefault',
});
</script>

<style lang="scss" scoped>
  .default-layout{
    width: 1200px;
    margin: 0 auto;
    min-height: calc(100vh - 194px - 337px);
  }
</style>
