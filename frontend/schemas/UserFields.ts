import { z } from 'zod';
import { ReadingPreferencesSchema } from './ReadingPreferences';

const UserFieldsSchema = z.object({
  username: z.string(),
  bio: z.string(),
  profilePicture: z.string(),
  preferredLanguage: z.string().nullable().optional(),
  readingPreferences: ReadingPreferencesSchema,
  dateOfBirth: z.string(),
  roles: z.string(),
});

type UserFields = z.infer<typeof UserFieldsSchema>;
export { type UserFields, UserFieldsSchema };
