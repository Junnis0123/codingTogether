import { reactive, ref } from '@vue/composition-api';
import defaultAxios from '@/services/axios/defaultAxios';
export default function joinMemberManager() {
    const session = reactive({
        id: '',
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
            }
            catch (e) {
                console.log(e);
            }
        }
    };
    return {
        session,
        formRefs,
        joinMember,
    };
}
//# sourceMappingURL=joinMember.js.map