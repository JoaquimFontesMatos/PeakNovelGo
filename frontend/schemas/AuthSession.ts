import * as yup from 'yup';
import { UserSchema } from '~/schemas/User';
import type { InferType } from 'yup';

const AuthSessionSchema = yup.object({
    accessToken: yup.string().required(),
    user: UserSchema.required(),
});

type AuthSession = InferType<typeof AuthSessionSchema>;

export { type AuthSession, AuthSessionSchema };