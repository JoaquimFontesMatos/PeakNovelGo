export interface HttpClient {
  request(input: RequestInfo, init: RequestInit): Promise<Response>;

  authorizedRequest(input: RequestInfo, init: RequestInit): Promise<Response>;

  debounce<T extends (...args: any[]) => void>(fn: T, delay: number): T;

  throttle<T extends (...args: any[]) => void>(fn: T, delay: number): T;
}
