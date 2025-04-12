import { z } from 'zod';

// LoginForm schema
export const loginFormSchema = z.object({
    email: z.email().refine(val => val.length <= 255, {
        error: 'Email must be at most 255 characters long',
    }),
    password: z
        .string()
        .nonempty({ error: 'Password is required' })
        .min(8, { error: 'Password must be at least 8 characters long' })
        .max(72, { error: 'Password must be at most 72 characters long' }),
});

export type LoginForm = z.infer<typeof loginFormSchema>;

// SignUpForm schema
export const signUpFormSchema = z.object({
    email: z.email().refine(val => val.length <= 255, {
        error: 'Email must be at most 255 characters long',
    }),
    password: z
        .string()
        .nonempty({ error: 'Password is required' })
        .min(8, { error: 'Password must be at least 8 characters long' })
        .max(72, { error: 'Password must be at most 72 characters long' }),
    username: z
        .string()
        .nonempty({ error: 'Username is required' })
        .min(3, { error: 'Username must be at least 3 characters long' })
        .max(255, { error: 'Username must be at most 255 characters long' }),
    dateOfBirth: z
        .string()
        .nonempty({ error: 'Date of Birth is required' })
        .regex(/^\d{4}-\d{2}-\d{2}$/, {
            error: 'Date of Birth must be in the format YYYY-MM-DD',
        }),
});

export type SignUpForm = z.infer<typeof signUpFormSchema>;
