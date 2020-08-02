import useSession from '@/services/login/session';
import useAxios from '@/services/axios/apiFactory';
import DefaultInterface, { CodingTogether } from '@/services/interface';
import { computed, reactive } from '@vue/composition-api';

const axios = useAxios(1);

interface ListResult extends DefaultInterface{
  data: string;
}

export default function homeManager() {
  const lists = reactive({
    all: {
      list: [],
      tab: 0,
    },
    own: {
      list: [],
      tab: 0,
    },
  });

  const getNickname = async () => {
    const result = await axios.get('/users/me');
    if (result) useSession().setNickname(result.data);
  };

  const getList = async () => {
    const result = await axios.get<ListResult>('/codingTogethers/');
    if (result) {
      lists.all.list = JSON.parse(result.data);
    }
  };
  const getOwnList = async () => {
    const result = await axios.get<ListResult>('/codingTogethers/me');
    if (result) {
      lists.own.list = JSON.parse(result.data);
    }
  };

  return {
    getNickname,
    getList,
    getOwnList,
    lists,
  };
}
