export interface HttpClient {
  request(input: RequestInfo, init: RequestInit): Promise<Response>;

  authorizedRequest(input: RequestInfo, init: RequestInit): Promise<Response>;
}
