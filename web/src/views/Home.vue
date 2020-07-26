<template>
  <div class="home-view">
    {{ nick }}
    {{ token }}
  </div>
</template>

<script lang="ts">
import { defineComponent, computed } from '@vue/composition-api';
import homeManager from '@/services/home/home';
import useSession from '../services/login/session';

export default defineComponent({
  name: 'Home',
  setup() {
    const manager = homeManager();
    manager.getNickname();

    return {
      token: useSession().getValue('token'),
      nick: computed(() => useSession().getValue('nickname')),
    };
  },
});
</script>

<style lang="scss" scoped>
  .login-view {
    &__contents {
      display: flex;
      justify-content: center;
      padding-top: 24px;
    }
  }
</style>
