<script setup lang="ts">
    const currentPage = ref(1);
    const currentLimit = ref(10);

    const { fetchingNovel, paginatedNovelsData } = storeToRefs(useNovelStore());

    const onPageChange = async (newPage: number, limit: number) => {
        if (newPage === currentPage.value && limit === currentLimit.value && paginatedNovelsData.value && paginatedNovelsData.value.page != 0) return;

        try {
            await useNovelStore().fetchNovels(newPage, limit);
        } catch {}
        currentPage.value = newPage;
        currentLimit.value = limit;
    };

    onMounted(async () => {
        await onPageChange(currentPage.value, currentLimit.value);
        // topRatedNovels.value = await fetchTopRatedNovels();
    });
</script>

<template>
    <section class="">
        <h2 class="mb-6 text-3xl font-bold text-primary-content">Top Rated</h2>
        <PaginatedNovelGallery
            v-show="!fetchingNovel"
            :errorMessage="paginatedNovelsData === null ? 'No Novels Found' : null"
            :paginatedData="paginatedNovelsData"
            @page-change="onPageChange"
        />
    </section>
</template>
