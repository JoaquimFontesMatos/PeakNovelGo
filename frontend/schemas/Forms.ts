import { z } from 'zod';

// LoginForm schema
export const loginFormSchema = z.object({
  email: z
    .string()
    .nonempty({ message: 'Email is required' })
    .email({ message: 'Please enter a valid email address' })
    .max(255, { message: 'Email must be at most 255 characters long' }),
  password: z
    .string()
    .nonempty({ message: 'Password is required' })
    .min(8, { message: 'Password must be at least 8 characters long' })
    .max(72, { message: 'Password must be at most 72 characters long' }),
});

export type LoginForm = z.infer<typeof loginFormSchema>;

// SignUpForm schema
export const signUpFormSchema = z.object({
  email: z
    .string()
    .nonempty({ message: 'Email is required' })
    .email({ message: 'Please enter a valid email address' })
    .max(255, { message: 'Email must be at most 255 characters long' }),
  password: z
    .string()
    .nonempty({ message: 'Password is required' })
    .min(8, { message: 'Password must be at least 8 characters long' })
    .max(72, { message: 'Password must be at most 72 characters long' }),
  username: z
    .string()
    .nonempty({ message: 'Username is required' })
    .min(3, { message: 'Username must be at least 3 characters long' })
    .max(255, { message: 'Username must be at most 255 characters long' }),
  dateOfBirth: z
    .string()
    .nonempty({ message: 'Date of Birth is required' })
    .regex(/^\d{4}-\d{2}-\d{2}$/, {
      message: 'Date of Birth must be in the format YYYY-MM-DD',
    }),
});

export type SignUpForm = z.infer<typeof signUpFormSchema>;
