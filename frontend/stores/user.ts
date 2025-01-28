import type { ReadingPreferences } from '~/models/ReadingPreferences';

export const useUserStore = defineStore('User', () => {
    const authStore = useAuthStore();
    const runtimeConfig = useRuntimeConfig();
    const url: string = runtimeConfig.public.apiUrl;

    const { user } = storeToRefs(authStore);

    const updatingUser = ref<boolean>(false);
    const updateUserError = ref<string | null>(null);

    const isReaderMode = ref<boolean>(true);

    const saveUserLocalStorage = () => {
        localStorage.setItem('user', JSON.stringify(user.value));
    };

    const updateUserFields = async (fields: {}) => {
        updatingUser.value = true;
        updateUserError.value = null;

        try {

            if (user.value === null) {
                updatingUser.value = false;
                updateUserError.value = 'You\'re not logged in!';
                return;
            }

            const response: Response | undefined = await authStore.authorizedFetch(url + '/user/' + user.value.ID + '/fields',
                {
                    method: 'PUT',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify(fields),
                });

            if (response === null) {
                updatingUser.value = false;
                updateUserError.value = 'Unexpected Error';
                return;
            }

            if (response.status === 200) {
                const message = await response.json();
                console.log(message);
            } else {
                const errorResponse = await response.json();

                updateUserError.value = errorResponse.error;
            }
        } catch (error) {
            console.log('Unexpected Error:', error);
            updateUserError.value = 'Unexpected error';
        } finally {
            updatingUser.value = false;
        }
    };

    return { user, isReaderMode, updatingUser, updateUserError, updateUserFields, saveUserLocalStorage };
});

if (import.meta.hot) {
    import.meta.hot.accept(acceptHMRUpdate(useUserStore, import.meta.hot));
}
