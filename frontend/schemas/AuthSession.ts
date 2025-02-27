import { z } from 'zod';
import { UserSchema } from '~/schemas/User';

const AuthSessionSchema = z.object({
  accessToken: z.string().nonempty('Access token is required'),
  user: UserSchema,
});

type AuthSession = z.infer<typeof AuthSessionSchema>;

export { type AuthSession, AuthSessionSchema };
