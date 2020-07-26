import axios from 'axios';
const baseURL = process.env.VUE_APP_API_BASE;
axios.defaults.baseURL = baseURL;
// 통신 후 인터셉트
axios.interceptors.response.use((response) => response.data, async (err) => {
    const error = err;
    alert(error.response.data.Message);
    return Promise.reject(error.response.data);
});
export default axios;
//# sourceMappingURL=defaultAxios.js.map