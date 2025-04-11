<script setup lang="ts">
    import type { PaginatedServerResponse } from '~/schemas/PaginatedServerResponse';
    import type { NovelSchema } from '~/schemas/Novel';
    import PaginatorNumbered from '~/components/paginator/PaginatorNumbered.vue';

    defineProps<{
        errorMessage: string | null;
        paginatedData: PaginatedServerResponse<typeof NovelSchema> | null;
        onPageChange: (newPage: number, limit: number) => Promise<void>;
    }>();

    const statusLabels: Record<string, string> = {
        Ongoing: 'blue-500',
        Completed: 'green-500',
        'On-Hold': 'yellow-500',
        Hiatus: 'yellow-500',
        Dropped: 'red-500',
    };

    const statusBorderClass = (status: string) => {
        switch (status) {
            case 'Completed':
                return 'border-green-500';
            case 'Ongoing':
                return 'border-blue-500';
            case 'Hiatus':
            case 'On-Hold':
                return 'border-yellow-500';
            case 'Cancelled':
            case 'Dropped':
                return 'border-red-500';
            default:
                return 'border-gray-500';
        }
    };

    const statusClass = (status: string): string => {
        return `text-${statusLabels[status] || 'gray-500'} ${statusBorderClass(status)}`;
    };
</script>

<template>
    <FeedbackErrorAlert v-if="errorMessage !== '' && errorMessage !== null">
        Error:
        {{ errorMessage == '' ? 'No Novels Found' : errorMessage }}
    </FeedbackErrorAlert>

    <div v-else-if="paginatedData && paginatedData.data.length > 0" class="@container">
        <ul v-auto-animate class="grid grid-cols-1 justify-center gap-2 @md:gap-5 @lg:grid-cols-2 @xl:gap-8 @3xl:grid-cols-3 @6xl:grid-cols-4">
            <li
                v-for="novel in paginatedData.data"
                :key="novel.novelUpdatesId"
                class="h-full rounded-md border-2 border-transparent bg-secondary transition-all duration-150 hover:scale-[1.01] hover:border-accent-gold hover:brightness-105 hover:drop-shadow-md"
            >
                <NuxtLink :to="'/novels/' + novel.novelUpdatesId">
                    <div class="flex h-24 w-full flex-row text-secondary-content @max-3xs:relative @sm:h-[7rem] @md:h-[9rem]">
                        <img :src="novel.coverUrl" alt="Cover Image" class="h-full w-1/4 rounded-s-md object-cover @max-3xs:w-full" />
                        <div
                            class="float-right w-3/4 p-4 @max-3xs:absolute @max-3xs:bottom-0 @max-3xs:left-0 @max-3xs:h-min @max-3xs:w-full @max-3xs:bg-primary/50 @max-3xs:p-1 @max-3xs:text-xs @max-3xs:backdrop-blur-md"
                        >
                            <p class="truncate">
                                {{ novel.title }}
                            </p>
                            <div class="flex flex-row gap-2 @max-3xs:hidden">
                                <!-- Status Label -->
                                <span class="truncate rounded-sm border p-0.5 text-xs" :class="statusClass(novel.status)">
                                    {{ novel.status }}
                                </span>
                                <!-- Total Chapters -->
                                <div class="flex items-center gap-1">
                                    <Icon name="fluent:book-open-16-regular" />
                                    <p class="text-xs">
                                        {{ novel.latestChapter }}
                                        Chapters
                                    </p>
                                </div>
                            </div>
                            <div>
                                <!-- TODO: Add Rating Label -->
                                <span></span>
                                <!-- TODO: Add General Info, like Rank, Views, Bookmarks and Comments -->
                                <div></div>
                            </div>
                        </div>
                    </div>
                </NuxtLink>
            </li>
        </ul>

        <SpacersVertical />

        <PaginatorNumbered :totalPages="paginatedData.totalPages" :total="paginatedData.total" @page-change="(page, limit) => onPageChange(page, limit)" />
    </div>
    <div v-else>
        <FeedbackErrorAlert>No novels found</FeedbackErrorAlert>
    </div>
</template>
