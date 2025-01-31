export interface BookmarkedNovel {
    ID: number;
    CreatedAt: string;
    UpdatedAt: string;
    DeletedAt?: string | null;

    novelId: number;
    userId: number;
    status: string;
    score: number;
    currentChapter: number;
}
