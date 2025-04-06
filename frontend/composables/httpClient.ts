import { storeToRefs } from 'pinia';
import type { HttpClient } from '~/interfaces/HttpClient';

export class FetchHttpClient implements HttpClient {
  private authStore: ReturnType<typeof useAuthStore>;

  constructor(authStore: ReturnType<typeof useAuthStore>) {
    this.authStore = authStore;
  }

  async request(input: RequestInfo, init: RequestInit = {}): Promise<Response> {
    return fetch(input, init);
  }

  async authorizedRequest(input: RequestInfo, init: RequestInit = {}): Promise<Response> {
    const { accessToken } = storeToRefs(this.authStore);

    if (!accessToken.value) {
      await this.authStore.refreshAccessToken();
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
      await this.authStore.refreshAccessToken();

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

  debounce<T extends (...args: any[]) => void>(fn: T, delay: number): T {
    let timeoutId: ReturnType<typeof setTimeout>;
    return function (this: any, ...args: any[]) {
      clearTimeout(timeoutId);
      timeoutId = setTimeout(() => fn.apply(this, args), delay);
    } as T;
  }

  throttle<T extends (...args: any[]) => void>(fn: T, delay: number): T {
    let lastCall = 0;
    return function (this: any, ...args: any[]) {
      const now = Date.now();
      if (now - lastCall >= delay) {
        lastCall = now;
        fn.apply(this, args);
      }
    } as T;
  }

  throttleWithFlush<T extends (...args: any[]) => void>(fn: T, delay: number): T & { flush: () => void } {
    let lastCall = 0;
    let lastArgs: any[] | null = null;
    let lastThis: any;
    let timeoutId: ReturnType<typeof setTimeout> | null = null;

    const throttledFunction = function (this: any, ...args: any[]) {
      const now = Date.now();
      lastArgs = args;
      lastThis = this;

      if (!lastCall) {
        // First call
        lastCall = now;
        fn.apply(this, args);
      } else if (now - lastCall >= delay) {
        // Time to execute
        lastCall = now;
        fn.apply(this, args);
      } else if (!timeoutId) {
        // Schedule execution
        timeoutId = setTimeout(() => {
          if (lastArgs) {
            lastCall = Date.now();
            fn.apply(lastThis, lastArgs);
            lastArgs = null;
            timeoutId = null;
          }
        }, delay - (now - lastCall));
      }
    } as T & { flush: () => void };

    throttledFunction.flush = () => {
      if (timeoutId && lastArgs) {
        clearTimeout(timeoutId);
        lastCall = Date.now();
        fn.apply(lastThis, lastArgs);
        lastArgs = null;
        timeoutId = null;
      }
    };

    return throttledFunction;
  }
}
