import type { User } from '~/models/User';

export interface AuthSession {
    accessToken: string;
    user: User;
}