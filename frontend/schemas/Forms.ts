import * as yup from "yup";

export const loginFormSchema = yup.object({
  email: yup
    .string()
    .required("Email is required")
    .email("Please enter a valid email address")
    .max(255, "Email must be at most 255 characters long"),
  password: yup
    .string()
    .required("Password is required")
    .min(8, "Password must be at least 8 characters long")
    .max(72, "Password must be at most 72 characters long"),
});

export type LoginForm = yup.InferType<typeof loginFormSchema>;

export const signUpFormSchema = yup.object({
  email: yup
    .string()
    .required("Email is required")
    .email("Please enter a valid email address")
    .max(255, "Email must be at most 255 characters long"),
  password: yup
    .string()
    .required("Password is required")
    .min(8, "Password must be at least 8 characters long")
    .max(72, "Password must be at most 72 characters long"),
  username: yup
    .string()
    .required("Username is required")
    .min(3, "Username must be at least 3 characters long")
    .max(255, "Username must be at most 255 characters long"),
  dateOfBirth: yup
    .string()
    .required("Date of Birth is required")
    .matches(
      /^\d{4}-\d{2}-\d{2}$/,
      "Date of Birth must be in the format YYYY-MM-DD"
    ),
});

export type SignUpForm = yup.InferType<typeof signUpFormSchema>;
