import type { LoginForm, SignInForm } from "~/models/Forms";
import type { LoginResponse } from "~/models/Response";
import type { User } from "~/models/User";

export const useAuthStore = defineStore("Auth", () => {
  const user = ref<User | null>(null);
  const accessToken = ref<string | null>(null);
  const runtimeConfig = useRuntimeConfig();
  const url = runtimeConfig.public.apiUrl;

  // Function to set the access token and user info
  const setSession = (loginResponse: LoginResponse) => {
    accessToken.value = loginResponse.accessToken;
    user.value = loginResponse.user;

    if (import.meta.client) {
      localStorage.setItem("accessToken", loginResponse.accessToken);
      localStorage.setItem("user", JSON.stringify(loginResponse.user));
    }
  };

  const clearSession = () => {
    accessToken.value = null;
    user.value = null;

    if (import.meta.client) {
      localStorage.removeItem("accessToken");
      localStorage.removeItem("user");
    }
  };

  const initSession = async () => {
    if (import.meta.client) {
      const storedAccessToken = localStorage.getItem("accessToken");
      const storedUser = localStorage.getItem("user");

      if (storedAccessToken) {
        accessToken.value = storedAccessToken;
      }

      if (storedUser) {
        user.value = JSON.parse(storedUser);
      }
    }

    await checkAuthStatus();
  };

  // Login function
  const login = async (form: LoginForm) => {
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
        setSession(loginResponse);
      } else {
        console.log(response.status);
      }
    } catch (error) {
      console.error("Login error:", error);
    }
  };

  // Automatically refresh the access token when needed
  const refreshAccessToken = async () => {
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
      await refreshAccessToken();
    }

    const headers = {
      ...init.headers,
      Authorization: `Bearer ${accessToken.value}`,
    };

    return fetch(input, { ...init, headers });
  };

  const signIn = async (form: SignInForm) => {
    try {
      const response = await fetch(`${url}/auth/signin`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(form),
      });

      if (response.status === 200) {
        const data = await response.json();
        user.value = data.user;
      } else {
        console.log(response.status);
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
    initSession,
    login,
    signIn,
    logout,
    refreshAccessToken,
    authorizedFetch,
    deleteUser,
  };
});

if (import.meta.hot) {
  import.meta.hot.accept(acceptHMRUpdate(useAuthStore, import.meta.hot));
}
