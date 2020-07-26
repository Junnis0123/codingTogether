<template>
  <div class="login-view">
    <v-toolbar
      color="primary"
      dark
      flat
    >
      <v-toolbar-title>로그인</v-toolbar-title>
    </v-toolbar>
    <div class="login-view__contents">
      <v-card outlined width="400" class="elevation-14">
        <v-card-text>
          <v-form ref="formRefs">
            <v-text-field
              ref="idRefs"
              required
              clearable autofocus
              v-model="session.id"
              @keyup.enter="login"
              :rules="[() => !!session.id || '아이디를 입력하세요!']"
              label="아이디"
              name="id"
              prepend-inner-icon="mdi-account"
              type="text"
            ></v-text-field>
            <v-text-field
              ref="passwordRefs"
              required
              clearable
              v-model="session.password"
              :rules="[() => !!session.password || '비밀번호를 입력하세요!']"
              @keyup.enter="login"
              id="password"
              label="비밀번호"
              name="password"
              prepend-inner-icon="mdi-lock"
              type="password"
            ></v-text-field>
          </v-form>
        </v-card-text>
        <v-card-actions>
          <router-link to="/joinMember">
            <v-btn>회원가입</v-btn>
          </router-link>
          <v-spacer></v-spacer>
          <v-btn @click="login" color="primary">로그인</v-btn>
        </v-card-actions>
      </v-card>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from '@vue/composition-api';
import loginManager from '@/services/login/login';

export default defineComponent({
  name: 'Login',
  setup() {
    return {
      ...loginManager(),
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
