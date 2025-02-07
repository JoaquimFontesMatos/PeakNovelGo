import type { LogEntry, LogLevel } from '~/models/LogEntry';

export class Logger {
  private transports: Array<(entry: LogEntry) => void> = [];

  private readonly levelPriority: Record<LogLevel, number> = {
    debug: 0,
    info: 1,
    warn: 2,
    error: 3,
  };
  private currentLevel: LogLevel;

  constructor(defaultLevel: LogLevel = 'info') {
    this.currentLevel = defaultLevel;
  }

  public setLevel(level: LogLevel): void {
    this.currentLevel = level;
  }

  private shouldLog(level: LogLevel): boolean {
    return this.levelPriority[level] >= this.levelPriority[this.currentLevel];
  }

  private log(level: LogLevel, message: string, context?: Record<string, unknown>, error?: Error): void {
    if (!this.shouldLog(level)) return;

    const entry: LogEntry = {
      level,
      message,
      timestamp: new Date().toISOString(),
      context,
      error,
    };

    if (level === 'debug') {
      console.log(entry);
    }
    // Send to transports (console, file, external service, etc.)
    this.sendToTransports(entry);
  }

  public addTransport(transport: (entry: LogEntry) => void): void {
    this.transports.push(transport);
  }

  private sendToTransports(entry: LogEntry): void {
    this.transports.forEach(transport => transport(entry));
  }

  // Public methods for each log level
  public debug(message: string, context?: Record<string, unknown>): void {
    this.log('debug', message, context);
  }

  public info(message: string, context?: Record<string, unknown>): void {
    this.log('info', message, context);
  }

  public warn(message: string, context?: Record<string, unknown>, error?: Error): void {
    this.log('warn', message, context, error);
  }

  public error(message: string, context?: Record<string, unknown>, error?: Error): void {
    this.log('error', message, context, error);
  }
}
