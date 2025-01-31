import type { User } from "./User";

export interface LoginResponse {
  accessToken: string;
  user: User;
}
