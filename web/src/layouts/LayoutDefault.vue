<template>
  <div>
    <v-app-bar
      app
      color="primary"
      dark >
      <v-toolbar-title>Coding Together</v-toolbar-title>
    </v-app-bar>
    <v-main>
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
