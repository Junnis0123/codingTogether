import { AxiosRequestConfig } from 'axios';
import Vue from 'vue';
import session from '@/services/login/session';
import axios from './defaultAxios';

const baseURL = process.env.VUE_APP_API_BASE;
const apiAxios = axios.create();

// 통신 전 인터셉트
apiAxios.interceptors.request.use((config) => {
  const con: AxiosRequestConfig = config;
  // const token: string = Vue.$cookies.get('access-token');

  con.baseURL = baseURL;

  con.headers = {
    Authorization: `Bearer ${session().getValue('token')}`,
  };
  return con;
}, async (err) => {
  const error = err;
  return Promise.reject(error);
});

// 통신 후 인터셉트
apiAxios.interceptors.response.use((response) => response,
  async (err) => {
    const error = err;
    return Promise.reject(error.response);
  });

export default apiAxios;
