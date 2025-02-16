import type { LogEntry } from '~/schemas/LogEntry';

export interface FileTransport {
  log(entry: LogEntry): Promise<void>;
}
