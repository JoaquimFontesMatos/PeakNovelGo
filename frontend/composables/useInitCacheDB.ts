const DB_NAME = 'novelCache';
const BOOKMARK_STORE = 'bookmarks';
const CHAPTER_STORE = 'chapters';
const RECENTLY_VISITED_NOVEL_STORE = 'recentlyVisitedNovels';

export const useIndexedDB = () => {
    let db: IDBDatabase | null = null;

    const initDB = async (): Promise<IDBDatabase> => {
        return new Promise((resolve, reject) => {
            if (db) {
                resolve(db);
                return;
            }

            const request = indexedDB.open(DB_NAME, 4);

            request.onupgradeneeded = event => {
                db = (event.target as IDBOpenDBRequest).result;
                if (!db.objectStoreNames.contains(BOOKMARK_STORE)) {
                    db.createObjectStore(BOOKMARK_STORE, { keyPath: 'cacheKey' });
                }
                if (!db.objectStoreNames.contains(CHAPTER_STORE)) {
                    db.createObjectStore(CHAPTER_STORE, { keyPath: 'cacheKey' });
                }
                if (!db.objectStoreNames.contains(RECENTLY_VISITED_NOVEL_STORE)) {
                    db.createObjectStore(RECENTLY_VISITED_NOVEL_STORE, { keyPath: 'cacheKey' });
                }
            };

            request.onsuccess = event => {
                db = (event.target as IDBOpenDBRequest).result;
                resolve(db);
            };

            request.onerror = event => {
                reject((event.target as IDBOpenDBRequest).error);
            };
        });
    };

    return {
        initDB,
    };
};
