import { reactive, ref } from '@vue/composition-api';
import useSession from '@/services/login/session';
import router from '@/router';
import useAxios from '@/services/axios/apiFactory';

const axios = useAxios(0);

interface Session {
  id: string,
  password: string,
}

interface LoginResult {
  success: boolean;
  message: string;
  errors: string;
  accessToken: string;
  refreshToken: string;
}

export default function loginManager() {
  const session = reactive<Session>({
    id: '',
    password: '',
  });
  const idRefs = ref<any>();
  const formRefs = ref<any>();
  const passwordRefs = ref();

  const login = async () => {
    if (formRefs.value.validate()) {
      const formdata = new FormData();
      formdata.append('userID', session.id);
      formdata.append('userPW', btoa(session.password));
      const result = await axios.post<LoginResult>('/auth/login', formdata);
      if (result) {
        useSession().setToken(result.accessToken);
        console.log(useSession().getValue('token'));
        await router.push('/');
      } else {
        idRefs.value.focus();
      }
    }
  };

  return {
    session,
    idRefs,
    formRefs,
    passwordRefs,
    login,
  };
}
