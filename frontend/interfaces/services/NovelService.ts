import type { ImportedNovel } from '~/schemas/ImportedNovel';
import type { Novel, NovelSchema } from '~/schemas/Novel';
import type { PaginatedServerResponse } from '~/schemas/PaginatedServerResponse';

export interface NovelService {
  fetchNovel(novelUpdatesId: string): Promise<Novel>;
  fetchNovels(page: number, limit: number): Promise<PaginatedServerResponse<typeof NovelSchema>>;
  fetchNovelsByTag(tag: string, page: number, limit: number): Promise<PaginatedServerResponse<typeof NovelSchema>>;
  fetchNovelsByAuthor(author: string, page: number, limit: number): Promise<PaginatedServerResponse<typeof NovelSchema>>;
  fetchNovelsByGenre(genre: string, page: number, limit: number): Promise<PaginatedServerResponse<typeof NovelSchema>>;
  importByNovelUpdatesId(novelUpdatesId: string): Promise<Novel>;
  importNovel(importedNovel: ImportedNovel): Promise<Novel>;
}
