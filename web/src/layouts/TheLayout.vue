<template>
    <component :is="layout">
        <slot />
    </component>
</template>

<script lang="ts">
import { defineComponent, computed } from '@vue/composition-api';
import LayoutDefault from '@/layouts/LayoutDefault.vue';
import LayoutFullscreen from '@/layouts/LayoutFullscreen.vue';
import ELayout from '@/layouts/layout.enum';

const defaultLayout = ELayout.Default;

export default defineComponent({
  name: 'TheLayout',
  components: {
    LayoutDefault,
    LayoutFullscreen,
  },
  setup(props, { root }) {
    const layout = computed(() => `layout-${root.$route.meta.layout || defaultLayout}`);
    return {
      layout,
    };
  },
});
</script>

<style lang="scss" scoped>

</style>
