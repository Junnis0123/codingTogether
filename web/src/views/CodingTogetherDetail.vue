<template>
  <div class="coding-together-detail">
    <v-row v-for="(list, idx) in lists"
           :key="idx">
      <v-col cols="12">
        <v-card
          class="mx-auto"
          tile
        >
          <v-list subheader header two-line>
            <v-subheader>My CodingTogether</v-subheader>
            <v-item-group v-model="list.tab" color="primary">
              <v-list-item v-for="(item, index) in list.list" :key="index" @click="() => {}">
                <v-list-item-content>
                  <v-list-item-title v-text="item.codingTogetherName"></v-list-item-title>
                  <v-list-item-subtitle class="text--primary"
                                        v-text="item.codingTogetherOrgnizerName">
                  </v-list-item-subtitle>
                  <v-list-item-subtitle v-text="item.codingTogetherCreateTime">
                  </v-list-item-subtitle>
                </v-list-item-content>
              </v-list-item>
            </v-item-group>
          </v-list>
        </v-card>
      </v-col>
    </v-row>

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
    manager.getList();
    manager.getOwnList();

    return {
      lists: manager.lists,
      token: useSession().getValue('token'),
      nick: computed(() => useSession().getValue('nickname')),
    };
  },
});
</script>

<style lang="scss" scoped>
  .home-view {
    padding: 12px;
    &__my-info{
      font-weight: bold;
    }
    &__contents {
      display: flex;
      justify-content: center;
      padding-top: 24px;
    }
  }
</style>
