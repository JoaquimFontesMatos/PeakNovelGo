import type { User } from '~/schemas/User';
import type { Novel } from '~/schemas/Novel';
import type { Paragraph } from '~/schemas/Paragraph';
import type { Chapter } from '~/schemas/Chapter';

type Role = 'admin' | 'moderator' | 'user';

type Permission<DataType = any> = boolean | ((user: User, data: DataType) => boolean);

type Permissions = {
    novels: {
        dataType: Novel;
        action: 'view' | 'create' | 'update';
    };
    chapters: {
        dataType: Chapter;
        action: 'view' | 'create' | 'update' | 'delete';
    };
    tts: {
        dataType: Paragraph[];
        action: 'generate';
    };
};

type RolesWithPermissions = {
    [R in Role]: {
        [Resource in keyof Permissions]: {
            [Action in Permissions[Resource]['action']]: Permission<Permissions[Resource]['dataType']>;
        };
    };
};

const ROLES: RolesWithPermissions = {
    admin: {
        novels: {
            view: true,
            create: true,
            update: true,
        },
        chapters: {
            view: true,
            create: true,
            update: true,
            delete: true,
        },
        tts: {
            generate: true,
        },
    },
    moderator: {
        novels: {
            view: true,
            create: true,
            update: true,
        },
        chapters: {
            view: true,
            create: true,
            update: true,
            delete: true,
        },
        tts: {
            generate: true,
        },
    },
    user: {
        novels: {
            view: true,
            create: false,
            update: false,
        },
        chapters: {
            view: true,
            create: false,
            update: false,
            delete: false,
        },
        tts: {
            generate: true,
        },
    },
};

export function hasPermission<Resource extends keyof Permissions>(
    user: User,
    resource: Resource,
    action: Permissions[Resource]['action'],
    data?: Permissions[Resource]['dataType'],
): boolean {
    return user.roles.split(';').some(roleStr => {
        if (!isRole(roleStr)) return false;

        const role: Role = roleStr as Role;
        const permission = ROLES[role]?.[resource]?.[action];

        if (permission == null) return false;

        if (typeof permission === 'boolean') return permission;
        return data != null && permission(user, data);
    });
}

// Type guard to check if a string is a valid Role
function isRole(role: string): role is Role {
    return role === 'admin' || role === 'moderator' || role === 'user';
}