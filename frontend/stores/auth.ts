import type {LoginForm, SignUpForm} from "~/veeSchemas/Forms";
import type {LoginResponse} from "~/models/Response";
import type {User} from "~/models/User";

export const useAuthStore = defineStore("Auth", () => {
  const user: Ref<User | null> = ref<User | null>(null);
  const accessToken: Ref<string | null> = ref<string | null>(null);
  const runtimeConfig = useRuntimeConfig();
  const url: string = runtimeConfig.public.apiUrl;
  const loginError: Ref<string | null> = ref<string | null>(null);
  const loadingLogin: Ref<boolean | null> = ref<boolean>(false);
  const signUpError: Ref<string | null> = ref<string | null>(null)

  // Function to set the access token and user info
  const setSession = (loginResponse: LoginResponse) => {
    accessToken.value = loginResponse.accessToken;
    user.value = loginResponse.user;

    if (import.meta.client) {
      localStorage.setItem("user", JSON.stringify(loginResponse.user));
    }
  };

  const clearSession = (): void => {
    accessToken.value = null;
    user.value = null;

    if (import.meta.client) {
      localStorage.removeItem("user");
    }
  };

  const initSession = async () => {
    if (import.meta.client) {
      const storedUser = localStorage.getItem("user");

      if (storedUser) {
        user.value = JSON.parse(storedUser);
      }
    }

    await checkAuthStatus();
  };

  // Login function
  const login = async (form: LoginForm): Promise<void> => {
    loadingLogin.value = true;
    try {
      const response = await fetch(`${url}/auth/login`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(form),
        credentials: "include",
      });

      if (response.status === 200) {
        const loginResponse: LoginResponse = await response.json();
        loginError.value = null;
        setSession(loginResponse);
      } else {
        // Handle errors based on response status
        const errorResponse = await response.json();
        if (response.status === 400) {
          loginError.value = errorResponse.error;
        } else if (response.status === 401) {
          loginError.value = errorResponse.error;
        } else if (response.status === 500) {
          loginError.value = errorResponse.error;
        } else {
          console.error("Unexpected error:", response.status);
        }
      }
    } catch (error) {
      console.error("Login error:", error);
    } finally {
      loadingLogin.value = false;
    }
  };

  // Automatically refresh the access token when needed
  const refreshAccessToken = async (): Promise<void> => {
    try {
      const response = await fetch(`${url}/auth/refresh-token`, {
        method: "POST",
        credentials: "include",
      });

      if (response.status === 200) {
        const loginResponse: LoginResponse = await response.json();
        setSession(loginResponse);
      } else {
        clearSession(); // Logout if refresh fails
      }
    } catch (error) {
      console.error("Failed to refresh token:", error);
      clearSession();
    }
  };

  // Attach the access token to every request
  const authorizedFetch = async (
    input: RequestInfo,
    init: RequestInit = {}
  ) => {


    if (
      !accessToken.value ||
      accessToken.value === "" ||
      accessToken.value === null
    ) {
      console.error("Accessed refresh")

      await refreshAccessToken();
    }

    const headers = {
      ...init.headers,
      Authorization: `Bearer ${accessToken.value}`,
    };

    return fetch(input, {...init, headers});
  };

  const signUp = async (form: SignUpForm) => {
    try {
      const registerData = {
        username: form.username,
        email: form.email,
        password: form.password,
        bio: "Please edit me",
        profilePicture: "Please edit me",
        dateOfBirth: form.dateOfBirth,
      };

      const response = await fetch(`${url}/auth/register`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(registerData),
      });

      if (response.status === 200) {
        const data = await response.json();
        user.value = data.user;
      } else {
        const errorResponse = await response.json();

        signUpError.value = errorResponse.error;
        console.log(response);
      }
    } catch (error) {
      console.error("Sign-in error:", error);
    }
  };

  // Logout function to clear session and notify backend
  const logout = async () => {
    try {
      await fetch(`${url}/auth/logout`, {
        method: "POST",
        credentials: "include", // Send the refresh token cookie
      });

      clearSession();
    } catch (error) {
      console.error("Logout error:", error);
    }
  };

  const deleteUser = async (id: number) => {
    accessToken.value = null;

    try {
      await authorizedFetch(`${url}/user/${id}`, {
        method: "DELETE",
        headers: {
          "Content-Type": "application/json",
        },
      });
    } catch (error) {
      console.log("Delete User", error);
    }
  };

  // Check if the user is logged in by verifying the access token
  const checkAuthStatus = async () => {
    if (accessToken.value) {
      // If there's an access token, consider the user as logged in
      return true;
    } else {
      // If no access token, try refreshing the session
      await refreshAccessToken();
      return !!accessToken.value; // Return true if access token is set
    }
  };

  return {
    user,
    accessToken,
    loginError,
    loadingLogin,
    initSession,
    login,
    signUp,
    logout,
    refreshAccessToken,
    authorizedFetch,
    deleteUser,
  };
});

if (import.meta.hot) {
  import.meta.hot.accept(acceptHMRUpdate(useAuthStore, import.meta.hot));
}
