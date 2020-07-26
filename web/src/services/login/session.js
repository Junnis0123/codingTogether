import { ref } from '@vue/composition-api';
const accessToken = ref();
export default function useSession() {
    return {
        setToken: (token) => {
            accessToken.value = token;
        },
        getToken: () => accessToken.value,
    };
}
//# sourceMappingURL=session.js.map