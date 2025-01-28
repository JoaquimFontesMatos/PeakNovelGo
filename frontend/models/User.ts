import type { ReadingPreferences } from '~/models/ReadingPreferences';

export interface User {
    ID: number;
    CreatedAt: string;
    UpdatedAt: string;
    DeletedAt?: string | null;

    username: string;
    email: string;
    emailVerified: boolean;
    profilePicture: string;
    bio: string;
    roles: string;
    lastLogin: string;
    dateOfBirth: string;
    preferredLanguage: string;
    readingPreferences: ReadingPreferences;
}
