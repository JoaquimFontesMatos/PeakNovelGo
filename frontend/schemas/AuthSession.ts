import * as yup from 'yup';
import { UserSchema } from '~/schemas/User';

export const AuthSessionSchema = yup.object({
    accessToken: yup.string().required(),
    user: UserSchema.required(),
});