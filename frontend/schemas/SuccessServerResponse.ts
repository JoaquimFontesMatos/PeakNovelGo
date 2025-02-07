import * as yup from 'yup';
import type { InferType } from 'yup';

const SuccessServerResponseSchema = yup.object({
    message: yup.string().required(),
});

type SuccessServerResponse = InferType<typeof SuccessServerResponseSchema>;
export { type SuccessServerResponse, SuccessServerResponseSchema };