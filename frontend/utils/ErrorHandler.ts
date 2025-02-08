import { ErrorBase } from '~/errors/ErrorBase';
import { handleSpecificError } from '~/errors/handling/ErrorConfig';

export const handleError = (error: unknown, context?: Record<string, unknown>) => {
  if (error instanceof ErrorBase) {
    handleSpecificError(error, context);
  } else {
    console.error('Non ErrorBase error encountered:', error);
  }
};
