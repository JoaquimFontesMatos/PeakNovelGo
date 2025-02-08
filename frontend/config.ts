import { FileTransport } from './errors/handling/FileTransport';
import { Logger } from './errors/handling/Logger';

// config.ts

let _apiUrl = 'http://localhost:8081';
let _runMode = 'development';

if (import.meta.client) {
  const runtimeConfig = useRuntimeConfig();
  _apiUrl = runtimeConfig.public.apiUrl;
  _runMode = runtimeConfig.public.runMode;
}

const logger = new Logger(_runMode === 'production' ? 'info' : 'debug');

const fileTransport = new FileTransport(_apiUrl);

logger.addTransport(fileTransport.log.bind(fileTransport));

export { logger };
