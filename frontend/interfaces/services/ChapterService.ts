import type { Chapter, ChapterSchema } from '~/schemas/Chapter';
import type { PaginatedServerResponse } from '~/schemas/PaginatedServerResponse';

export interface ChapterService {
  fetchChapter(novelUpdatesId: string, chaptNo: number): Promise<Chapter>;
  fetchChapters(novelUpdatesId: string, page: number, limit: number): Promise<PaginatedServerResponse<typeof ChapterSchema>>;
}
