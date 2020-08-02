import defaultAxios from '@/services/axios/defaultAxios';
import apiAxios from '@/services/axios/apiAxios';
import { AxiosResponse } from 'axios';

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

  const get = async<T = AxiosResponse> (url: string, params?: object) => {
    try {
      const result = await axios.get<T>(url, params);
      return result.data;
    } catch (e) {
      alert(e.message);
      return false;
    }
  };
  const post = async<T = AxiosResponse> (url: string, formData?: FormData) => {
    try {
      const result = await axios.post<T>(url, formData);
      return result.data;
    } catch (e) {
      alert(e.message);
      return false;
    }
  };

  return {
    get,
    post,
  };
}
