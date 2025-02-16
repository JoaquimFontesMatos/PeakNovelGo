import type { SuccessServerResponse } from "~/schemas/SuccessServerResponse";

export interface UserService {
  updateUserFields(fields: {}, userId: number):  Promise<SuccessServerResponse>;
  deleteUser(userId: number):  Promise<SuccessServerResponse>;
}
