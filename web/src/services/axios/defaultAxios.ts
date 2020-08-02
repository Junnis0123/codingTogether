import axios from 'axios';

const baseURL = process.env.VUE_APP_API_BASE;

axios.defaults.baseURL = baseURL;

// 통신 후 인터셉트
axios.interceptors.response.use((response) => response,
  async (err) => Promise.reject(err.response.data));

export default axios;
