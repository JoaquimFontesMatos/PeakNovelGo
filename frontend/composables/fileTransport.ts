import type { FileTransport } from '~/interfaces/FileTransport';
import type { LogEntry } from '~/schemas/LogEntry';

export class BaseFileTransport implements FileTransport {
  private url: string;

  constructor(url: string) {
    this.url = url;
  }

  public async log(entry: LogEntry): Promise<void> {
    try {
      if (this.url === '') return;

      const accessToken = localStorage.getItem('accessToken') || '';

      await $fetch(this.url + '/log/', {
        method: 'POST',
        headers: {
          Authorization: `Bearer ${accessToken}`,
        }
      });
    } catch (error) {
      console.error('Failed to send log:', error);
    }
  }
}
