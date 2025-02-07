import { setupErrorHandling } from '~/errors/ErrorHandler';
import { logger } from '~/config';
import { FileTransport } from '~/errors/FileTransport';
import { resolve } from 'path';
// plugins/error-handler.client.ts
export default defineNuxtPlugin(nuxtApp => {
  // Access the Vue application instance from nuxtApp.vueApp
  const originalErrorHandler = nuxtApp.vueApp.config.errorHandler;
  nuxtApp.vueApp.config.errorHandler = (error, vm, info) => {
    // Log the error (make sure logger is imported or accessible)
    logger.error('Vue error', { component: vm?.$options?.name, info }, error as Error | undefined);
    // Call the original error handler if it exists
    if (originalErrorHandler) originalErrorHandler(error, vm, info);
  };
});
