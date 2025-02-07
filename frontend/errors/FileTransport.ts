// FileTransport.ts
import type { LogEntry } from '~/models/LogEntry';

export class FileTransport {
  private url: string;

  constructor(url: string) {
    this.url = url;
  }

  // Example: Client-side error logging function
  public async log(entry: LogEntry): Promise<void> {
    try {
      if (this.url === '') return;

      await $fetch(this.url + '/log/', {
        method: 'POST',
        body: entry,
      });
    } catch (error) {
      console.error('Failed to send log:', error);
    }
  }
}
