import type { HttpClient } from '~/interfaces/HttpClient';

export class FetchHttpClient implements HttpClient {
  async request(input: RequestInfo, init: RequestInit = {}): Promise<Response> {
    return fetch(input, init);
  }

  async authorizedRequest(input: RequestInfo, init: RequestInit = {}): Promise<Response> {
    const authStore = useAuthStore();
    const { accessToken } = storeToRefs(authStore);

    if (!accessToken.value) {
      await authStore.refreshAccessToken();
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
      await authStore.refreshAccessToken();

      if (accessToken.value) {
        response = await this.request(input, {
          ...init,
          headers: {
            ...init.headers,
            Authorization: `Bearer ${accessToken.value}`,
          },
        });
      }
    }

    return response;
  }
}
