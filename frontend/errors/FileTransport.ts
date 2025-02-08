import type { LogEntry } from '~/schemas/LogEntry';

export class FileTransport {
  private url: string;

  constructor(url: string) {
    this.url = url;
  }

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
