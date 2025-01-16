export interface Chapter {
  ID: number;
  CreatedAt: string;
  UpdatedAt: string;
  DeletedAt?: string | null;

  chapterNo: number;
  novelId: number;
  title: string;
  chapterUrl: string;
  body: string;
}
