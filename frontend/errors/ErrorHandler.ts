import type { Logger } from './Logger';

export function setupErrorHandling(logger: Logger): void {
  process.on('uncaughtException', error => {
    logger.error('Uncaught exception', {}, error);
    process.exit(1);
  });

  // Handle unhandled promise rejections (e.g., uncaught async errors)
  process.on('unhandledRejection', reason => {
    logger.error('Unhandled promise rejection', {}, reason instanceof Error ? reason : new Error(String(reason)));
  });
}
