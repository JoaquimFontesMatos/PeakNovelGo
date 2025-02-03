import * as yup from 'yup';

export const ErrorServerResponseSchema = yup.object({
    error: yup.string().required(),
});