import { reactive, ref } from '@vue/composition-api';
import defaultAxios from '@/services/axios/defaultAxios';
import useSession from '@/services/login/session';
import router from '@/router';
export default function loginManager() {
    const formHasErrors = ref(false);
    const session = reactive({
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
                console.log(useSession().getToken());
                router.push('/');
            }
            catch (e) {
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
//# sourceMappingURL=login.js.map