import type { AuthSession } from '~/schemas/AuthSession';
import type { LoginForm, SignUpForm } from '~/schemas/Forms';
import type { SuccessServerResponse } from '~/schemas/SuccessServerResponse';

export interface AuthService {
  login(form: LoginForm): Promise<AuthSession>;
  signUp(form: SignUpForm): Promise<SuccessServerResponse>;
  refreshAccessToken(refreshToken: string): Promise<AuthSession>;
  verifyToken(token: string): Promise<SuccessServerResponse>;
  logout(refreshToken: string): Promise<SuccessServerResponse>;
}
