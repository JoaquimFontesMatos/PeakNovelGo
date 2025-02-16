import type { SuccessServerResponse } from '~/schemas/SuccessServerResponse';
import { AuthError } from '~/errors/AuthError';
import type { ErrorHandler } from '~/interfaces/ErrorHandler';
import type { HttpClient } from '~/interfaces/HttpClient';
import type { ResponseParser } from '~/interfaces/ResponseParser';
import { BaseUserService } from '~/services/UserService';
import type { UserService } from '~/interfaces/services/UserService';

export const useUserStore = defineStore('User', () => {
  const authStore = useAuthStore();
  const runtimeConfig = useRuntimeConfig();
  const url: string = runtimeConfig.public.apiUrl;

  const { user } = storeToRefs(authStore);

  // Initialize user service
  const httpClient: HttpClient = new FetchHttpClient();
  const responseParser: ResponseParser = new ZodResponseParser();
  const userService: UserService = new BaseUserService(url, httpClient, responseParser);

  // Initialize error handler
  const errorHandler: ErrorHandler = new BaseErrorHandler();

  // Handling update user Variables
  const updatingUser: Ref<boolean> = ref<boolean>(false);
  const updateMessage: Ref<string | null> = ref<string | null>(null);

  // Handling user prefs Variables
  const isReaderMode: Ref<boolean> = ref<boolean>(true);

  // Handling update user Variables
  const deletingUser: Ref<boolean> = ref<boolean>(false);
  const deleteMessage: Ref<string | null> = ref<string | null>(null);

  const saveUserLocalStorage = (): void => {
    localStorage.setItem('user', JSON.stringify(user.value));
  };

  const updateUserFields = async (fields: {}): Promise<void> => {
    updatingUser.value = true;
    updateMessage.value = null;

    try {
      if (!authStore.isUserLoggedIn()) {
        throw new AuthError({
          type: 'UNAUTHORIZED_ERROR',
          message: "You're not logged in!",
          cause: 'User tried to update user fields without being logged in.',
        });
      }

      const serverResponse: SuccessServerResponse = await userService.updateUserFields(fields, user.value!!.ID);

      updateMessage.value = serverResponse.message;
    } catch (error) {
      errorHandler.handleError(error, { user: user, fields: fields, location: 'user.ts -> updateUserFields' });
      throw error;
    } finally {
      updatingUser.value = false;
    }
  };

  const deleteUser = async (): Promise<void> => {
    deletingUser.value = true;
    deleteMessage.value = null;

    try {
      if (!authStore.isUserLoggedIn()) {
        throw new AuthError({
          type: 'UNAUTHORIZED_ERROR',
          message: "You're not logged in!",
          cause: 'User tried to update user fields without being logged in.',
        });
      }

      const serverResponse: SuccessServerResponse = await userService.deleteUser(user.value!!.ID);

      deleteMessage.value = serverResponse.message;
    } catch (error) {
      errorHandler.handleError(error, { user: user, location: 'user.ts -> deleteUser' });
      authStore.clearSession();

      throw error;
    } finally {
      deletingUser.value = false;
    }
  };

  return {
    user,
    isReaderMode,
    updatingUser,
    updateMessage,
    deletingUser,
    deleteMessage,
    deleteUser,
    updateUserFields,
    saveUserLocalStorage,
  };
});

if (import.meta.hot) {
  import.meta.hot.accept(acceptHMRUpdate(useUserStore, import.meta.hot));
}
