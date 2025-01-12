export interface Chapter {
  id: number;
  createdAt: string;
  updatedAt: string;
  deletedAt?: string | null;

  chapter_no: number;
  novel_id: number;
  title: string;
  chapter_url: string;
  body: string;
}
