import type { LogEntry, LogLevel } from '~/schemas/LogEntry';

export interface Logger {
  setLevel(level: LogLevel): void;
  addTransport(transport: (entry: LogEntry) => void): void;
  debug(message: string, context?: Record<string, unknown>): void;
  info(message: string, context?: Record<string, unknown>): void;
  warn(message: string, context?: Record<string, unknown>, error?: Error): void;
  error(message: string, context?: Record<string, unknown>, error?: Error): void;
}
