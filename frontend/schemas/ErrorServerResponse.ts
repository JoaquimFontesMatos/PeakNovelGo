import { z } from 'zod';

const ErrorServerResponseSchema = z.object({
  error: z.string().nonempty('Error message is required'),
  code: z.string().nonempty('Error code is required')
});

type ErrorServerResponse = z.infer<typeof ErrorServerResponseSchema>;

export { type ErrorServerResponse, ErrorServerResponseSchema };
