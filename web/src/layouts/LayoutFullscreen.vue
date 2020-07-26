<template>
  <div>
    <slot />
  </div>
</template>

<script lang="ts">
import { defineComponent } from '@vue/composition-api';

export default defineComponent({
  name: 'LayoutFullscreen',
});
</script>

<style lang="scss" scoped>

</style>
