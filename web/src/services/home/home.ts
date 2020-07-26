import apiAxios from '@/services/axios/apiAxios';
import useSession from '@/services/login/session';

export default function homeManager() {
  const getNickname = async () => {
    const result = await apiAxios.get('/users/me');
    useSession().setNickname(result.Data);
  };

  return {
    getNickname,
  };
}
