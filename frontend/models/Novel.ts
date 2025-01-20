export interface Novel {
  ID: number;
  CreatedAt: string;
  UpdatedAt: string;
  DeletedAt?: string | null;

  title: string;
  synopsis: string;
  coverUrl: string;
  language: string;
  status: string;
  novelUpdatesUrl: string;
  tags: Tag[];
  authors: Author[];
  genres: Genre[];
  year: string;
  releaseFrequency: string;
  novelUpdatesId: string;
}

export interface Tag {
  id: number;
  name: string;
  description: string;
}

export interface Author {
  id: number;
  name: string;
}

export interface Genre {
  id: number;
  name: string;
  description: string;
}
