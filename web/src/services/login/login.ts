import { reactive, ref } from '@vue/composition-api';
import defaultAxios from '@/services/axios/defaultAxios';
import useSession from '@/services/login/session';
import router from '@/router';

interface Session {
  id: string,
  password: string,
}

export default function loginManager() {
  const formHasErrors = ref<Boolean>(false);
  const session = reactive<Session>({
    id: '',
    password: '',
  });
  const idRefs = ref();
  const formRefs = ref();
  const passwordRefs = ref();

  const login = async () => {
    if (formRefs.value.validate()) {
      try {
        const formdata = new FormData();
        formdata.append('user_id', session.id);
        formdata.append('user_pw', btoa(session.password));
        const result = await defaultAxios.post('/auth/login', formdata);
        useSession().setToken(result.Data);
        console.log(useSession().getValue('token'));
        router.push('/');
      } catch (e) {
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
