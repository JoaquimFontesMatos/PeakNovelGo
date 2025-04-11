<script setup lang="ts">
    import type { PaginatedServerResponse } from '~/schemas/PaginatedServerResponse';
    import type { ChapterSchema } from '~/schemas/Chapter';
    import PaginatorNumbered from '~/components/paginator/PaginatorNumbered.vue';
    const { novelTitle } = useRoute().params;

    defineProps<{
        errorMessage: string | null;
        paginatedData: PaginatedServerResponse<typeof ChapterSchema> | null;
        onPageChange: (newPage: number, limit: number) => Promise<void>;
    }>();

    const [parent] = useAutoAnimate({ duration: 150 });
</script>

<template>
    <FeedbackErrorAlert v-if="!paginatedData || !paginatedData.data || paginatedData.data.length === 0">
        Error:
        {{ errorMessage === '' || errorMessage === null ? 'No Chapters Found' : errorMessage }}
    </FeedbackErrorAlert>

    <div v-else-if="paginatedData?.data?.length > 0" class="min-h-72">
        <PaginatorNumbered :totalPages="paginatedData.totalPages" :total="paginatedData.total" @page-change="(page, limit) => onPageChange(page, limit)" />
        <SpacersVertical />
        <ul ref="parent" class="grid grid-cols-1 justify-center gap-10 gap-y-2 md:grid-cols-2 lg:grid-cols-3">
            <li v-for="chapter in paginatedData.data" :key="chapter.chapterNo">
                <NuxtLink
                    class="block h-full rounded-md border-2 border-transparent bg-secondary transition-all duration-150 hover:scale-[1.01] hover:border-accent-gold hover:brightness-105 hover:drop-shadow-md"
                    :to="'/novels/' + novelTitle + '/' + chapter.chapterNo"
                >
                    <span class="ml-2">Chapter {{ chapter.chapterNo }}</span>
                </NuxtLink>
            </li>
        </ul>
        <SpacersVertical />
        <PaginatorNumbered
            v-show="paginatedData.limit > 25"
            :totalPages="paginatedData.totalPages"
            :total="paginatedData.total"
            @page-change="(page, limit) => onPageChange(page, limit)"
        />
    </div>

    <FeedbackErrorAlert v-else>Error: Failed to Fetch Chapters</FeedbackErrorAlert>
</template>
