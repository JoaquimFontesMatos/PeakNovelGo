import * as yup from 'yup';
import type { InferType } from 'yup';

const ErrorServerResponseSchema = yup.object({
    error: yup.string().required(),
});

type ErrorServerResponse = InferType<typeof ErrorServerResponseSchema>;

export { type ErrorServerResponse, ErrorServerResponseSchema };