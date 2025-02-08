import { ProjectError } from '~/errors/ProjectError';

export const parseJSONPromise = async (response: Response): Promise<any> => {
  try {
    return await response.json();
  } catch (error) {
    throw new ProjectError({
      name: 'INTERNAL_SERVER_ERROR',
      message: 'An unexpected error occurred while parsing JSON',
      cause: error,
    });
  }
};
