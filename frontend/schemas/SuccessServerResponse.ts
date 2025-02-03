import * as yup from 'yup';

export const SuccessServerResponseSchema = yup.object({
    message: yup.string().required(),
});