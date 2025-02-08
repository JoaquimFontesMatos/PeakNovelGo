import { logger } from '~/config';
import { AuthError } from '~/errors/AuthError';
import { BookmarkError } from '~/errors/BookmarkError';
import { ChapterError } from '~/errors/ChapterError';
import { NovelError } from '~/errors/NovelError';
import { ProjectError } from '~/errors/ProjectError';
import { TtsError } from '~/errors/TtsError';
import { UserError } from '~/errors/UserError';
import type { ToastIcon } from '~/schemas/Toast';

export const handleError = (error: unknown, context?: Record<string, unknown>) => {
  const toastStore = useToastStore();

  let errorType = 'Unknown Error';
  let errorMessage = 'An unexpected error occurred';
  let icon: ToastIcon = 'none';

  if (error instanceof AuthError) {
    errorType = 'AuthError';
    errorMessage = `Authentication Error: ${error.message}`;
    icon = 'auth';
  } else if (error instanceof UserError) {
    errorType = 'UserError';
    errorMessage = `User Error: ${error.message}`;
    icon = 'user';
  } else if (error instanceof ProjectError) {
    errorType = 'ProjectError';
    errorMessage = `Project Error: ${error.message}`;
    icon = 'project';
  } else if (error instanceof NovelError) {
    errorType = 'NovelError';
    errorMessage = `Novel Error: ${error.message}`;
    icon = 'novel';
  } else if (error instanceof ChapterError) {
    errorType = 'ChapterError';
    errorMessage = `Chapter Error: ${error.message}`;
    icon = 'chapter';
  } else if (error instanceof BookmarkError) {
    errorType = 'BookmarkError';
    errorMessage = `Bookmark Error: ${error.message}`;
    icon = 'bookmark';
  } else if (error instanceof TtsError) {
    errorType = 'TtsError';
    errorMessage = `TTS Error: ${error.message}`;
    icon = 'tts';
  }

  toastStore.addToast(errorMessage, 'error', icon);

  // Log with structured metadata
  if (error instanceof Error) {
    logger.error(
      `Error Type: ${errorType}`,
      {
        context,
        message: errorMessage,
        stack: error.stack,
      },
      error
    );
  }
};
