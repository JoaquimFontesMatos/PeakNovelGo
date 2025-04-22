import type { BookmarkedNovel, BookmarkedNovelSchema } from '~/schemas/BookmarkedNovel';
import type { PaginatedServerResponse } from '~/schemas/PaginatedServerResponse';
import type { SuccessServerResponse } from '~/schemas/SuccessServerResponse';
import type { NovelSchema } from '~/schemas/Novel';

export interface BookmarkService {
    bookmarkNovel(novelId: number, userId: number): Promise<BookmarkedNovel>;
    updateBookmark(updatedBookmark: BookmarkedNovel): Promise<BookmarkedNovel>;
    unbookmarkNovel(novelId: number, userId: number): Promise<SuccessServerResponse>;
    fetchBookmarkedNovelByUser(novelId: string, userId: number): Promise<BookmarkedNovel>;
    fetchBookmarkedNovelsByUser(userId: number, page: number, limit: number): Promise<PaginatedServerResponse<typeof NovelSchema>>;
}
