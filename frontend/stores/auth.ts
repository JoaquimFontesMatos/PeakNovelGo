import type { LoginForm, SignUpForm } from '~/schemas/Forms';
import { AuthService } from '~/services/AuthService';
import { AuthError } from '~/errors/AuthError';
import { ProjectError } from '~/errors/ProjectError';
import { UserError } from '~/errors/UserError';
import type { SuccessServerResponse } from '~/schemas/SuccessServerResponse';
import type { AuthSession } from '~/schemas/AuthSession';
import type { User } from '~/schemas/User';
import { logger } from '~/config';

export const useAuthStore = defineStore('Auth', () => {
  // fetch the base url from the nuxt config
  const runtimeConfig = useRuntimeConfig();
  const url: string = runtimeConfig.public.apiUrl;

  // Initialize auth service
  const authService: AuthService = new AuthService(url);

  // Auth Session Info
  const user: Ref<User | null> = ref<User | null>(null);
  const accessToken: Ref<string | null> = ref<string | null>(null);

  // Handling Login Variables
  const loginError: Ref<string | null> = ref<string | null>(null);
  const loadingLogin: Ref<boolean> = ref<boolean>(false);

  // Handling Sign Up Variables
  const signUpError: Ref<string | null> = ref<string | null>(null);
  const loadingSignUp: Ref<boolean> = ref<boolean>(false);
  const signUpMessage: Ref<string | null> = ref<string | null>(null);

  // Handling Logout Variables
  const logoutError: Ref<string | null> = ref<string | null>(null);
  const loadingLogout: Ref<boolean> = ref<boolean>(false);
  const logoutMessage: Ref<string | null> = ref<string | null>(null);

  // Handling Refresh Token Variables
  const refreshTokenError: Ref<string | null> = ref<string | null>(null);

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
    return user.value === null;
  };

  // Login function
  const login = async (form: LoginForm): Promise<void> => {
    loadingLogin.value = true;
    loginError.value = null;

    try {
      const loginResponse: AuthSession = await authService.login(form);

      setSession(loginResponse);
    } catch (error) {
      let errorType = 'Unknown';

      if (error instanceof AuthError) {
        loginError.value = error.message;
        errorType = 'AuthError';
      } else if (error instanceof ProjectError) {
        loginError.value = error.message;
        errorType = 'ProjectError';
      } else {
        loginError.value = 'An unknown error occurred.';
      }

      if (error instanceof Error) {
        logger.error('Login Error', { userEmail: form.email, location: 'auth.ts -> login', type: errorType }, error);
      }

      clearSession();
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
    refreshTokenError.value = null;

    try {
      const loginResponse: AuthSession = await authService.refreshAccessToken();

      setSession(loginResponse);

      scheduleTokenRefresh(15 * 60 * 1000);
    } catch (error) {
      if (error instanceof AuthError) {
        loginError.value = error.message;
        console.error('AuthError', error);
      } else if (error instanceof ProjectError) {
        loginError.value = error.message;
        console.error('ProjectError:', error);
      } else {
        loginError.value = 'An unknown error occurred.';
        console.error('Unexpected error:', error);
      }
      clearSession();
    }
  };

  // Attach the access token to every request
  const authorizedFetch = async (input: RequestInfo, init: RequestInit = {}) => {
    if (!accessToken.value) {
      await refreshAccessToken();
    }

    let response = await fetch(input, {
      ...init,
      headers: {
        ...init.headers,
        Authorization: `Bearer ${accessToken.value}`,
      },
    });

    // If unauthorized, try refreshing the token and retrying
    if (response.status === 401) {
      await refreshAccessToken();

      if (accessToken.value) {
        response = await fetch(input, {
          ...init,
          headers: {
            ...init.headers,
            Authorization: `Bearer ${accessToken.value}`,
          },
        });
      }
    }

    return response;
  };

  const signUp = async (form: SignUpForm): Promise<void> => {
    loadingSignUp.value = true;
    signUpError.value = null;
    signUpMessage.value = null;

    try {
      const successServerResponse: SuccessServerResponse = await authService.signUp(form);
      signUpMessage.value = successServerResponse.message;
    } catch (error) {
      if (error instanceof UserError) {
        signUpError.value = error.message;
        console.error('UserError', error);
      } else if (error instanceof ProjectError) {
        signUpError.value = error.message;
        console.error('ProjectError:', error);
      } else {
        signUpError.value = 'An unknown error occurred.';
        console.error('Unexpected error:', error);
      }
    } finally {
      loadingSignUp.value = false;
    }
  };

  // Logout function to clear session and notify backend
  const logout = async (): Promise<void> => {
    loadingLogout.value = true;
    logoutError.value = null;

    try {
      const successServerResponse: SuccessServerResponse = await authService.logout();
      logoutMessage.value = successServerResponse.message;

      clearSession();
    } catch (error) {
      if (error instanceof ProjectError) {
        logoutError.value = error.message;
        console.error('ProjectError:', error);
      } else {
        logoutError.value = 'An unknown error occurred.';
        console.error('Unexpected error:', error);
      }
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
    loginError,
    loadingLogin,
    signUpError,
    loadingSignUp,
    signUpMessage,
    logoutError,
    loadingLogout,
    initSession,
    login,
    signUp,
    logout,
    refreshAccessToken,
    authorizedFetch,
    clearSession,
    isUserLoggedIn,
  };
});

if (import.meta.hot) {
  import.meta.hot.accept(acceptHMRUpdate(useAuthStore, import.meta.hot));
}
