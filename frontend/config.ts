import { BaseFileTransport } from './composables/fileTransport';
import { BaseLogger } from './composables/logger';

// config.ts

let _apiUrl = 'http://localhost:8081';
let _runMode = 'development';

if (import.meta.client) {
  const runtimeConfig = useRuntimeConfig();
  _apiUrl = runtimeConfig.public.apiUrl;
  _runMode = runtimeConfig.public.runMode;
}

const logger = new BaseLogger(_runMode === 'production' ? 'info' : 'debug');

const fileTransport = new BaseFileTransport(_apiUrl);

logger.addTransport(fileTransport.log.bind(fileTransport));

export { logger };
