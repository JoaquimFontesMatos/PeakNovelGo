import type { ZodSchema } from 'zod';
import { ProjectError } from '~/errors/ProjectError';
import type { ResponseParser } from '~/interfaces/ResponseParser';

export class ZodResponseParser implements ResponseParser {
  async parseJSON<T>(response: Response): Promise<T> {
    try {
      return await response.json();
    } catch (error) {
      throw new ProjectError({
        type: 'INTERNAL_SERVER_ERROR',
        message: 'An unexpected error occurred while parsing JSON',
        cause: error,
      });
    }
  }
  validateSchema<T>(schema: ZodSchema<T>, data: unknown): T {
    return schema.parse(data);
  }
}
