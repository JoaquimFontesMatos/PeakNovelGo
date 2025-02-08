import { z } from 'zod';

// Log level enum (explicit type)
const LogLevelSchema = z.enum(['debug', 'info', 'warn', 'error']);

// Error schema (optional)
const ErrorSchema = z
  .object({
    name: z.string(),
    message: z.string(),
    stack: z.string().optional(),
  })
  .optional();

// LogEntry schema
const LogEntrySchema = z.object({
  level: LogLevelSchema,
  message: z.string(),
  timestamp: z.string(),
  context: z.record(z.unknown()).optional(),
  error: ErrorSchema,
});

type LogEntry = z.infer<typeof LogEntrySchema>;

export { type LogEntry, LogEntrySchema };
