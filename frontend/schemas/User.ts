import * as yup from 'yup';
import { ReadingPreferencesSchema } from '~/schemas/ReadingPreferences';

export const UserSchema = yup.object({
    ID: yup.number().required('ID is required'),
    CreatedAt: yup.string().required('Creation date is required'),
    UpdatedAt: yup.string().required('Updated date is required'),
    // DeletedAt is optional: you can use .nullable() and .notRequired() to allow missing values
    DeletedAt: yup.string().nullable().notRequired(),
    username: yup.string().required('Username is required'),
    email: yup.string().email('Invalid email format').required('Email is required'),
    emailVerified: yup.boolean().required('Email verified flag is required'),
    profilePicture: yup.string().required('Profile picture URL is required'),
    bio: yup.string().required('Bio is required'),
    roles: yup.string().required('Roles are required'),
    lastLogin: yup.string().required('Last login date is required'),
    dateOfBirth: yup.string().required('Date of birth is required'),
    preferredLanguage: yup.string().nullable().notRequired(),
    readingPreferences: ReadingPreferencesSchema.required('Reading preferences are required'),
});