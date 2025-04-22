<script setup lang="ts">
    const currentPage = ref(1);
    const currentLimit = ref(10);

    const bookmarkStore = useBookmarkStore();

    const { fetchingBookmarkedNovels, paginatedBookmarkedNovels } = storeToRefs(bookmarkStore);

    const onPageChange = async (newPage: number, limit: number) => {
        if (newPage === currentPage.value && limit === currentLimit.value && paginatedBookmarkedNovels.value && paginatedBookmarkedNovels.value.page != 0)
            return;

        try {
            await bookmarkStore.fetchBookmarkedNovelsByUser(newPage, limit);
        } catch {}
        currentPage.value = newPage;
        currentLimit.value = limit;
    };

    onMounted(async () => {
        await onPageChange(currentPage.value, currentLimit.value);
        // recentlyUpdated.value = await fetchRecentlyUpdatedNovels();
    });
</script>

<template>
    <section>
        <h2 class="mb-6 text-3xl font-bold text-primary-content">Following</h2>
        <PaginatedNovelGallery
            v-show="!fetchingBookmarkedNovels"
            :errorMessage="paginatedBookmarkedNovels === null ? 'No Novels Found' : null"
            :paginatedData="paginatedBookmarkedNovels"
            @page-change="onPageChange"
        />
    </section>
</template>
