import { reactive } from '@vue/composition-api';

interface UserInfo {
  [index: string]: string,
  accessToken: string;
  nickname: string;
}

const userInfo = reactive<UserInfo>({
  accessToken: '',
  nickname: '',
});

export default function useSession() {
  return {
    setToken: (token: string) => {
      userInfo.accessToken = token;

      localStorage.removeItem('token');
      localStorage.setItem('token', token);
    },
    setNickname: (nick: string) => {
      userInfo.nickname = nick;
    },
    getValue: (key: string) => userInfo[key] || localStorage.getItem(key),
  };
}
