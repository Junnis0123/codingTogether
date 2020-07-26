import defaultAxios from '@/services/axios/defaultAxios';
import apiAxios from '@/services/axios/apiAxios';

enum AxiosType {
  Default,
  Api,
}

const createAxios = (type: AxiosType) => {
  switch (type) {
    case AxiosType.Default:
    default:
      return defaultAxios;
    case AxiosType.Api:
      return apiAxios;
  }
};

export default function useAxios(type: AxiosType) {
  const axios = createAxios(type);

  const get = async (url: string, params?: object) => {
    try {
      const result = await axios.get(url, params);
      return result;
    } catch (e) {
      alert(e.Message);
      return false;
    }
  };

  return {
    get,
  };
}
