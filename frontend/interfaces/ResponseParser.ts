import type { ZodSchema } from 'zod';

export interface ResponseParser {
  parseJSON<T>(response: Response): Promise<T>;
  validateSchema<T>(schema: ZodSchema<T>, data: unknown): T;
}
