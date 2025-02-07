import { ProjectError } from '~/errors/ProjectError';

function assert(
    condition: boolean,
    message: string,
    context?: unknown,
): asserts condition {
    if (!condition) {
        throw new ProjectError({
            name: 'VALIDATION_ERROR',
            message: `${message}. Context: ${JSON.stringify(context)}`,
        });
    }
}

export { assert };