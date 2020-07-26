<template>
  <v-app>
    <the-layout>
      <router-view />
    </the-layout>
  </v-app>
</template>

<script lang="ts">
import { defineComponent } from '@vue/composition-api';
import TheLayout from '@/layouts/TheLayout.vue';

export default defineComponent({
  name: 'App',
  components: {
    TheLayout,
  },
});
</script>
