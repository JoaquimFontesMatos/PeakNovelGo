import { ErrorBase, type ErrorNames } from '../ErrorBase';
import type { ToastIcon, ToastType } from '~/schemas/Toast';
import { logger } from '~/config';

// Define a structured type-safe configuration
type ErrorConfig = {
  [N in keyof ErrorNames]: {
    [T in ErrorNames[N]]: {
      showToast?: boolean;
      showDialog?: boolean;
      log?: boolean;
      icon?: ToastIcon;
      toastType?: ToastType;
      formatMessage?: (error: Error) => string;
    };
  };
};

// Define the error configurations
export const errorConfig: ErrorConfig = {
  AuthError: {
    REFRESH_TOKEN_ERROR: {
      showToast: true,
      showDialog: false,
      log: true,
      icon: 'auth',
      toastType: 'error',
      formatMessage: error => `Authentication Error: ${error.message}`,
    },
    LOGOUT_ERROR: {
      showToast: true,
      showDialog: false,
      log: true,
      icon: 'auth',
      toastType: 'error',
      formatMessage: error => `Authentication Error: ${error.message}`,
    },
    INVALID_SESSION_DATA: {
      showToast: true,
      showDialog: false,
      log: true,
      icon: 'auth',
      toastType: 'error',
      formatMessage: error => `Authentication Error: ${error.message}`,
    },
    UNAUTHORIZED_ERROR: {
      showToast: true,
      showDialog: true,
      log: true,
      icon: 'auth',
      toastType: 'warning',
      formatMessage: error => `Authentication Error: ${error.message}`,
    },
    INVALID_CREDENTIALS_ERROR: {
      showToast: true,
      showDialog: false,

      log: false,
      icon: 'auth',
      toastType: 'error',
      formatMessage: error => `Authentication Error: ${error.message}`,
    },
    SIGN_UP_ERROR: {
      showToast: true,
      showDialog: false,
      log: true,
      toastType: 'error',
      formatMessage: error => `Authentication Error: ${error.message}`,
    },
  },
  UserError: {
    USER_NOT_FOUND_ERROR: {
      showToast: true,
      showDialog: false,
      log: true,
      icon: 'user',
      toastType: 'error',
      formatMessage: error => `User Error: ${error.message}`,
    },
    USER_DEACTIVATED_ERROR: {
      showToast: true,
      showDialog: false,
      log: true,
      icon: 'user',
      toastType: 'error',
      formatMessage: error => `User Error: ${error.message}`,
    },
    USER_CONFLICT_ERROR: {
      showToast: true,
      showDialog: false,
      log: true,
      icon: 'user',
      toastType: 'error',
      formatMessage: error => `User Error: ${error.message}`,
    },
    GET_USER_ERROR: {
      showToast: true,
      showDialog: false,
      log: true,
      icon: 'user',
      toastType: 'error',
      formatMessage: error => `User Error: ${error.message}`,
    },
    CREATE_USER_ERROR: {
      showToast: true,
      showDialog: true,
      log: true,
      toastType: 'error',
      formatMessage: error => `User Error: ${error.message}`,
    },
    UPDATE_USER_ERROR: {
      showToast: true,
      showDialog: false,
      log: true,
      icon: 'user',
      toastType: 'error',
      formatMessage: error => `User Error: ${error.message}`,
    },
    INVALID_USER_DATA: {
      showToast: true,
      showDialog: false,
      log: true,
      icon: 'user',
      toastType: 'error',
      formatMessage: error => `User Error: ${error.message}`,
    },
  },
  ProjectError: {
    INTERNAL_SERVER_ERROR: {
      showToast: true,
      showDialog: false,
      log: true,
      icon: 'project',
      toastType: 'error',
      formatMessage: error => `Project Error: ${error.message}`,
    },
    CREATE_NOVEL_ERROR: {
      showToast: true,
      showDialog: true,
      log: true,
      toastType: 'error',
      formatMessage: error => `Project Error: ${error.message}`,
    },
    UPDATE_NOVEL_ERROR: {
      showToast: true,
      showDialog: false,
      log: true,
      icon: 'project',
      toastType: 'error',
      formatMessage: error => `Project Error: ${error.message}`,
    },
    NETWORK_ERROR: {
      showToast: true,
      showDialog: false,
      log: true,
      icon: 'project',
      toastType: 'error',
      formatMessage: error => `Project Error: ${error.message}`,
    },
    VALIDATION_ERROR: {
      showToast: true,
      showDialog: false,
      log: true,
      icon: 'project',
      toastType: 'error',
      formatMessage: error => `Project Error: ${error.message}`,
    },
    INVALID_RESPONSE_ERROR: {
      showToast: true,
      showDialog: false,
      log: true,
      icon: 'project',
      toastType: 'error',
      formatMessage: error => `Project Error: ${error.message}`,
    },
  },
  BookmarkError: {
    CREATE_BOOKMARK_ERROR: {
      showToast: true,
      showDialog: true,
      log: true,
      toastType: 'error',
      formatMessage: error => `Bookmark Error: ${error.message}`,
    },
    UPDATE_BOOKMARK_ERROR: {
      showToast: true,
      showDialog: false,
      log: true,
      icon: 'bookmark',
      toastType: 'error',
      formatMessage: error => `Bookmark Error: ${error.message}`,
    },
    BOOKMARK_NOT_FOUND: {
      showToast: false,
      showDialog: false,
      log: true,
      icon: 'bookmark',
      toastType: 'error',
      formatMessage: error => `Bookmark Error: ${error.message}`,
    },
  },
  NovelError: {
    GET_NOVEL_ERROR: {
      showToast: true,
      showDialog: false,
      log: true,
      icon: 'novel',
      toastType: 'error',
      formatMessage: error => `Novel Error: ${error.message}`,
    },
    CREATE_NOVEL_ERROR: {
      showToast: true,
      showDialog: true,
      log: true,
      toastType: 'error',
      formatMessage: error => `Novel Error: ${error.message}`,
    },
    UPDATE_NOVEL_ERROR: {
      showToast: true,
      showDialog: false,
      log: true,
      icon: 'novel',
      toastType: 'error',
      formatMessage: error => `Novel Error: ${error.message}`,
    },
    NOVEL_NOT_FOUND: {
      showToast: false,
      showDialog: false,
      log: true,
      icon: 'novel',
      toastType: 'error',
      formatMessage: error => `Novel Error: ${error.message}`,
    },
  },
  ChapterError: {
    GET_CHAPTER_ERROR: {
      showToast: true,
      showDialog: false,
      log: true,
      icon: 'chapter',
      toastType: 'error',
      formatMessage: error => `Chapter Error: ${error.message}`,
    },
    CREATE_CHAPTER_ERROR: {
      showToast: true,
      showDialog: true,
      log: true,
      toastType: 'error',
      formatMessage: error => `Chapter Error: ${error.message}`,
    },
    UPDATE_CHAPTER_ERROR: {
      showToast: true,
      showDialog: false,
      log: true,
      icon: 'chapter',
      toastType: 'error',
      formatMessage: error => `Chapter Error: ${error.message}`,
    },
    NO_CHAPTER_FOUND: {
      showToast: false,
      showDialog: false,
      log: true,
      icon: 'chapter',
      toastType: 'error',
      formatMessage: error => `Chapter Error: ${error.message}`,
    },
  },
  TtsError: {
    GET_TTS_ERROR: {
      showToast: true,
      showDialog: false,
      log: true,
      icon: 'tts',
      toastType: 'error',
      formatMessage: error => `TTS Error: ${error.message}`,
    },
    CREATE_TTS_ERROR: {
      showToast: true,
      showDialog: true,
      log: true,
      toastType: 'error',
      formatMessage: error => `TTS Error: ${error.message}`,
    },
    UPDATE_TTS_ERROR: {
      showToast: true,
      showDialog: false,
      log: true,
      icon: 'tts',
      toastType: 'error',
      formatMessage: error => `TTS Error: ${error.message}`,
    },
  },
};

export function handleSpecificError<N extends keyof ErrorNames, T extends ErrorNames[N]>(error: ErrorBase<N, T>, context?: Record<string, unknown>) {
  const toastStore = useToastStore();

  // Explicitly type `error.name` and `error.type`
  const errorName = error.name as N;
  const errorType = error.type as T;

  // Ensure the correct config lookup
  const config = errorConfig[errorName]?.[errorType] as ErrorConfig[N][T] | undefined;

  if (!config) {
    console.error('Unhandled error:', error);
    return;
  }

  // Handle logging
  if (config.log) {
    const generalError = { name: errorName, message: error.message, stack: error.stack };
    logger.error(
      `Error Type: ${errorType} - ${errorName}`,
      {
        context: context,
        message: error.message,
        stack: error.stack,
      },
      generalError
    );
  }

  // Handle toast notifications
  if (config.showToast) {
    // Ensure error is passed correctly with type assertion
    const message = config.formatMessage ? config.formatMessage(error as ErrorBase<N, T>) : error.message;
    toastStore.addToast(message, config.toastType || 'error', config.icon || 'none');
  }

  // Handle dialogs
  if (config.showDialog) {
    console.log(`Show dialog: ${error.message}`);
  }
}
