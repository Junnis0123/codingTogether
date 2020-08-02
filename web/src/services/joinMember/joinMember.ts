import { reactive, ref } from '@vue/composition-api';
import useAxios from '@/services/axios/apiFactory';
import DefaultInterface from '@/services/interface';

interface JoinValue{
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
  const formRefs = ref<any>();

  const joinMember = async () => {
    if (formRefs.value.validate()) {
      try {
        const formdata = new FormData();
        formdata.append('userID', session.id);
        formdata.append('userPW', btoa(session.password));
        formdata.append('userNickname', session.nickname);
        const result = await axios.post('/users/', formdata);
        console.log(result);
      } catch (e) {
        console.log(e);
      }
    }
  };
  const checkUserId = async () => {
    const result = await axios.get<DefaultInterface>(`/auth/duplication/${session.id}`);
    if (result) session.idConfirm = result.success;
  };

  return {
    session,
    formRefs,
    joinMember,
    checkUserId,
  };
}
