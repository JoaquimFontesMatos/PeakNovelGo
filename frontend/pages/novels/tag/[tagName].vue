<script setup lang="ts">
    const { tagName } = useRoute().params as { tagName: string };
    const currentPage = ref(1);
    const currentLimit = ref(10);

    onMounted(async () => {
        try {
            await onPageChange(currentPage.value, currentLimit.value);
        } catch {}
    });

    const { fetchingNovel, paginatedNovelsDataByTag } = storeToRefs(useNovelStore());

    const onPageChange = async (newPage: number, limit: number) => {
        if (newPage === currentPage.value && limit === currentLimit.value && paginatedNovelsDataByTag.value && paginatedNovelsDataByTag.value.page != 0) return;
        try {
            await useNovelStore().fetchNovelsByTag(tagName as string, newPage, limit);
        } catch {}
        currentPage.value = newPage;
        currentLimit.value = limit;
    };
</script>

<template>
    <Container>
        <RouteTree
            :routes="[
                { path: '/', name: 'Home' },
                { path: '/novels', name: 'Novels' },
                {
                    path: '/novels/tag/' + tagName,
                    name: tagName as string,
                },
            ]"
        />

        <VerticalSpacer />

        <LoadingBar v-show="fetchingNovel" />

        <PaginatedNovelGallery
            v-show="!fetchingNovel"
            :errorMessage="paginatedNovelsDataByTag === null ? 'No Novels Found' : null"
            :paginatedData="paginatedNovelsDataByTag"
            @page-change="onPageChange"
        />
    </Container>
</template>
