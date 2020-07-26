import { reactive, ref } from '@vue/composition-api';
import defaultAxios from '@/services/axios/defaultAxios';
import useAxios from '@/services/axios/apiFactory';

interface JoinValue {
  id: string,
  idConfirm: boolean,
  password: string,
  passwordConfirm: string,
  nickname: string,
}

const axios = useAxios(0);

export default function joinMemberManager() {
  const session = reactive<JoinValue>({
    id: '',
    idConfirm: false,
    password: '',
    passwordConfirm: '',
    nickname: '',
  });
  const formRefs = ref();

  const joinMember = async () => {
    if (formRefs.value.validate()) {
      try {
        const formdata = new FormData();
        formdata.append('user_id', session.id);
        formdata.append('user_pw', btoa(session.password));
        formdata.append('user_nickname', session.nickname);
        const result = await defaultAxios.post('/users/', formdata);
        console.log(result);
      } catch (e) {
        console.log(e);
      }
    }
  };
  const checkUserId = async () => {
    const result = await axios.get(`/auth/duplication/${session.id}`);
    session.idConfirm = result.success;
  };

  return {
    session,
    formRefs,
    joinMember,
    checkUserId,
  };
}
