import { z } from 'zod';
import { ReadingPreferencesSchema } from '~/schemas/ReadingPreferences';

const UserSchema = z.object({
  ID: z.number(),
  CreatedAt: z.string(),
  UpdatedAt: z.string(),
  DeletedAt: z.string().nullable().optional(),
  username: z.string(),
  email: z.string().email(),
  emailVerified: z.boolean(),
  profilePicture: z.string(),
  bio: z.string(),
  roles: z.string(),
  lastLogin: z.string(),
  dateOfBirth: z.string(),
  preferredLanguage: z.string().nullable().optional(),
  readingPreferences: ReadingPreferencesSchema,
});

type User = z.infer<typeof UserSchema>;
export { type User, UserSchema };
