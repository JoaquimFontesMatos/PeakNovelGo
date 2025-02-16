import { handleSpecificError } from '~/config/errorConfig';
import { ErrorBase } from '~/errors/ErrorBase';
import type { ErrorHandler } from '~/interfaces/ErrorHandler';

export class BaseErrorHandler implements ErrorHandler {
  handleError(error: unknown, context?: Record<string, unknown>): void {
    if (error instanceof ErrorBase) {
      handleSpecificError(error, context);
    } else {
      console.error('Non ErrorBase error encountered:', error);
    }
  }
}
