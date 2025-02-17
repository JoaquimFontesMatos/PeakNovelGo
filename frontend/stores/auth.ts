import type { LoginForm, SignUpForm } from '~/schemas/Forms';
import type { SuccessServerResponse } from '~/schemas/SuccessServerResponse';
import type { AuthSession } from '~/schemas/AuthSession';
import type { User } from '~/schemas/User';
import type { ErrorHandler } from '~/interfaces/ErrorHandler';
import type { AuthService } from '~/interfaces/services/AuthService';
import { BaseAuthService } from '~/services/AuthService';
import type { HttpClient } from '~/interfaces/HttpClient';
import type { ResponseParser } from '~/interfaces/ResponseParser';

export const useAuthStore = defineStore('Auth', () => {
  const runtimeConfig = useRuntimeConfig();
  const url: string = runtimeConfig.public.apiUrl;
  const httpClient: HttpClient = new FetchHttpClient(useAuthStore());
  const responseParser: ResponseParser = new ZodResponseParser();
  const $authService: AuthService = new BaseAuthService(url, httpClient, responseParser);
  const $errorHandler: ErrorHandler = new BaseErrorHandler();

  // Auth Session Info
  const user = ref<User | null>(null);
  const accessToken = ref<string | null>(null);

  // Handling Login Variables
  const loadingLogin = ref<boolean>(false);

  // Handling Sign Up Variables
  const loadingSignUp = ref<boolean>(false);
  const signUpMessage = ref<string | null>(null);

  // Handling Logout Variables
  const loadingLogout = ref<boolean>(false);
  const logoutMessage: Ref<string | null> = ref<string | null>(null);

  // Handling Verify Token Variables
  const loadingVerifyToken = ref<boolean>(false);
  const verifyTokenMessage = ref<string | null>(null);

  // Function to set the access token and user info
  const setSession = (loginResponse: AuthSession) => {
    accessToken.value = loginResponse.accessToken;
    user.value = loginResponse.user;

    if (import.meta.client) {
      localStorage.setItem('user', JSON.stringify(loginResponse.user));
      localStorage.setItem('accessToken', loginResponse.accessToken);
    }
  };

  const clearSession = (): void => {
    accessToken.value = null;
    user.value = null;

    if (import.meta.client) {
      localStorage.removeItem('user');
      localStorage.removeItem('accessToken');
    }
  };

  const isUserLoggedIn = (): boolean => {
    return user.value !== null;
  };

  // Login function
  const login = async (form: LoginForm): Promise<void> => {
    loadingLogin.value = true;

    try {
      const loginResponse: AuthSession = await $authService.login(form);

      setSession(loginResponse);
    } catch (error) {
      $errorHandler.handleError(error, { userEmail: form.email, location: 'auth.ts -> login' });
      clearSession();
      throw error;
    } finally {
      loadingLogin.value = false;
    }
  };

  const scheduleTokenRefresh = (expiresIn: number) => {
    setTimeout(async () => {
      await refreshAccessToken();
    }, expiresIn - 60 * 1000); // Refresh 1 minute before expiry
  };

  const refreshAccessToken = async (): Promise<void> => {
    try {
      const loginResponse: AuthSession = await $authService.refreshAccessToken();

      setSession(loginResponse);

      scheduleTokenRefresh(15 * 60 * 1000);
    } catch (error) {
      $errorHandler.handleError(error, { user: user, accessToken: accessToken, location: 'auth.ts -> refreshAccessToken' });
      clearSession();
      throw error;
    }
  };

  const verifyToken = async (token: string): Promise<void> => {
    loadingVerifyToken.value = true;
    verifyTokenMessage.value = null;
    try {
      const message: SuccessServerResponse = await $authService.verifyToken(token);
      verifyTokenMessage.value = message.message;
    } catch (error) {
      $errorHandler.handleError(error, { user: user, token: token, location: 'auth.ts -> verifyToken' });
      clearSession();
      throw error;
    } finally {
      loadingVerifyToken.value = false;
    }
  };

  const signUp = async (form: SignUpForm): Promise<void> => {
    loadingSignUp.value = true;
    signUpMessage.value = null;

    try {
      const successServerResponse: SuccessServerResponse = await $authService.signUp(form);
      signUpMessage.value = successServerResponse.message;
    } catch (error) {
      $errorHandler.handleError(error, { user: user, location: 'auth.ts -> signUp' });
      throw error;
    } finally {
      loadingSignUp.value = false;
    }
  };

  // Logout function to clear session and notify backend
  const logout = async (): Promise<void> => {
    loadingLogout.value = true;

    try {
      const successServerResponse: SuccessServerResponse = await $authService.logout();
      logoutMessage.value = successServerResponse.message;

      clearSession();
    } catch (error) {
      $errorHandler.handleError(error, { user: user, accessToken: accessToken, location: 'auth.ts -> logout' });
      throw error;
    } finally {
      loadingLogout.value = false;
    }
  };

  const initSession = async () => {
    if (import.meta.client) {
      const storedUser = localStorage.getItem('user');
      const storedAccessToken = localStorage.getItem('accessToken');

      if (storedUser) user.value = JSON.parse(storedUser);
      if (storedAccessToken) accessToken.value = storedAccessToken;

      if (!accessToken.value) {
        // If no access token, try refreshing the session
        await refreshAccessToken();
      }
    }
  };

  return {
    user,
    accessToken,
    loadingLogin,
    loadingSignUp,
    signUpMessage,
    loadingLogout,
    loadingVerifyToken,
    verifyTokenMessage,
    initSession,
    login,
    signUp,
    logout,
    refreshAccessToken,
    verifyToken,
    setSession,
    clearSession,
    isUserLoggedIn,
  };
});

if (import.meta.hot) {
  import.meta.hot.accept(acceptHMRUpdate(useAuthStore, import.meta.hot));
}
