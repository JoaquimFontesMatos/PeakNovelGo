import { UserService } from '~/services/UserService';
import type { SuccessServerResponse } from '~/schemas/SuccessServerResponse';
import { AuthError } from '~/errors/AuthError';
import { ProjectError } from '~/errors/ProjectError';
import { UserError } from '~/errors/UserError';

export const useUserStore = defineStore('User', () => {
    const authStore = useAuthStore();
    const runtimeConfig = useRuntimeConfig();
    const url: string = runtimeConfig.public.apiUrl;

    const { user } = storeToRefs(authStore);

    // Initialize user service
    const userService: UserService = new UserService(url);

    // Handling update user Variables
    const updatingUser: Ref<boolean> = ref<boolean>(false);
    const updateUserError: Ref<string | null> = ref<string | null>(null);
    const updateMessage: Ref<string | null> = ref<string | null>(null);

    // Handling user prefs Variables
    const isReaderMode: Ref<boolean> = ref<boolean>(true);

    // Handling update user Variables
    const deletingUser: Ref<boolean> = ref<boolean>(false);
    const deleteUserError: Ref<string | null> = ref<string | null>(null);
    const deleteMessage: Ref<string | null> = ref<string | null>(null);

    const saveUserLocalStorage = (): void => {
        localStorage.setItem('user', JSON.stringify(user.value));
    };

    const updateUserFields = async (fields: {}): Promise<void> => {
        updatingUser.value = true;
        updateUserError.value = null;
        updateMessage.value = null;

        if (authStore.isUserLoggedIn()) {
            updatingUser.value = false;
            updateUserError.value = 'You\'re not logged in!';
            return;
        }

        try {
            const serverResponse: SuccessServerResponse = await userService.updateUserFields(fields, user.value!!.ID);

            updateMessage.value = serverResponse.message;
        } catch (error) {
            if (error instanceof AuthError) {
                updateUserError.value = error.message;
                console.error('AuthError', error);
            } else if (error instanceof ProjectError) {
                updateUserError.value = error.message;
                console.error('ProjectError:', error);
            } else if (error instanceof UserError) {
                updateUserError.value = error.message;
                console.error('ProjectError:', error);
            } else {
                updateUserError.value = 'An unknown error occurred.';
                console.error('Unexpected error:', error);
            }
        } finally {
            updatingUser.value = false;
        }
    };

    const deleteUser = async (): Promise<void> => {
        deletingUser.value = true;
        deleteUserError.value = null;
        deleteMessage.value = null;

        if (authStore.isUserLoggedIn()) {
            deletingUser.value = false;
            deleteUserError.value = 'You\'re not logged in!';
            return;
        }

        try {
            const serverResponse: SuccessServerResponse = await userService.deleteUser(user.value!!.ID);

            deleteMessage.value = serverResponse.message;
        } catch (error) {
            if (error instanceof AuthError) {
                deleteUserError.value = error.message;
                console.error('AuthError', error);
            } else if (error instanceof ProjectError) {
                deleteUserError.value = error.message;
                console.error('ProjectError:', error);
            } else if (error instanceof UserError) {
                deleteUserError.value = error.message;
                console.error('ProjectError:', error);
            } else {
                deleteUserError.value = 'An unknown error occurred.';
                console.error('Unexpected error:', error);
            }
            authStore.clearSession();
        } finally {
            deletingUser.value = false;
        }
    };

    return {
        user,
        isReaderMode,
        updatingUser,
        updateUserError,
        updateMessage,
        deletingUser,
        deleteUserError,
        deleteMessage,
        deleteUser,
        updateUserFields,
        saveUserLocalStorage,

    };
});

if (import.meta.hot) {
    import.meta.hot.accept(acceptHMRUpdate(useUserStore, import.meta.hot));
}
