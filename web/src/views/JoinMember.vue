<template>
  <div class="login-view">
    <v-toolbar
      color="primary"
      dark
      flat
    >
      <v-toolbar-title>회원가입</v-toolbar-title>
    </v-toolbar>
    <div class="login-view__contents">
      <v-card outlined width="400" class="elevation-14">
        <v-card-text>
          <v-form ref="formRefs">
            <v-text-field
              required
              clearable autofocus
              v-model="session.id"
              @keyup.enter="joinMember"
              :rules="[() => !!session.id || '아이디를 입력하세요!',
              () => !!session.idConfirm || '아이디 중복체크를 진행해주세요!']"
              label="아이디"
              name="id"
              prepend-inner-icon="mdi-account"
              append-icon="mdi-send"
              @click:append="checkUserId"
              type="text"
            ></v-text-field>
            <v-text-field
              required
              clearable
              v-model="session.password"
              :rules="[() => !!session.password || '비밀번호를 입력하세요!']"
              @keyup.enter="joinMember"
              id="password"
              label="비밀번호"
              name="password"
              prepend-inner-icon="mdi-lock"
              type="password"
            ></v-text-field>
            <v-text-field
              required
              clearable
              v-model="session.passwordConfirm"
              :rules="[() => session.password === session.passwordConfirm || '비밀번호가 일치하지 않습니다.']"
              @keyup.enter="login"
              id="passwordConfirm"
              label="비밀번호확인"
              name="passwordConfirm"
              prepend-inner-icon="mdi-lock"
              type="password"
            ></v-text-field>
            <v-text-field
              required
              clearable
              v-model="session.nickname"
              :rules="[() => !!session.nickname || '닉네임을 입력하세요!']"
              @keyup.enter="joinMember"
              id="nickname"
              label="닉네임"
              name="nickname"
              prepend-inner-icon="mdi-lock"
              type="text"
            ></v-text-field>
          </v-form>
        </v-card-text>
        <v-card-actions>
          <router-link to="/Login">
            <v-btn>취소</v-btn>
          </router-link>
          <v-spacer></v-spacer>
          <v-btn @click="joinMember" color="primary">회원가입</v-btn>
        </v-card-actions>
      </v-card>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from '@vue/composition-api';
import joinMemberManager from '@/services/joinMember/joinMember';

export default defineComponent({
  name: 'Login',
  setup() {
    return {
      ...joinMemberManager(),
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
