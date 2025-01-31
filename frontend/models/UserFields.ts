import type { ReadingPreferences } from './ReadingPreferences';

export interface UserFields {
    username: string,
    bio: string,
    profilePicture: string,
    preferredLanguage: string,
    readingPreferences: ReadingPreferences,
    dateOfBirth: string,
    roles: string
}