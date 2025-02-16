export interface ErrorHandler {
  handleError(error: unknown, context?: Record<string, unknown>): void;
}
