<script setup lang="ts">
    import { hasPermission } from '~/config/permissionsConfig';

    useHead({
        title: '📖 Novels',
    });

    const currentPage = ref(1);
    const currentLimit = ref(10);

    const userStore = useUserStore();

    onMounted(async () => {
        try {
            await onPageChange(currentPage.value, currentLimit.value);
        } catch {}
    });

    const open = ref(false);

    const { fetchingNovel, paginatedNovelsData } = storeToRefs(useNovelStore());

    const onPageChange = async (newPage: number, limit: number) => {
        if (newPage === currentPage.value && limit === currentLimit.value && paginatedNovelsData.value && paginatedNovelsData.value.page != 0) return;
        try {
            await useNovelStore().fetchNovels(newPage, limit);
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
            ]"
        />

        <VerticalSpacer />

        <LoadingBar v-show="fetchingNovel" />

        <div v-show="!fetchingNovel && userStore.user && hasPermission(userStore.user, 'novels', 'update')">
            <NuxtLink :to="'/novels/update'">
                <CircularButton :icon-name="'fluent-mdl2:refresh'" :icon-size="24" :padding="4" @click="navigateTo('novels/update')" />
            </NuxtLink>
            <VerticalSpacer />
        </div>

        <PaginatedNovelGallery
            v-show="!fetchingNovel"
            :errorMessage="paginatedNovelsData === null ? 'No Novels Found' : null"
            :paginatedData="paginatedNovelsData"
            @page-change="onPageChange"
        />
    </Container>
</template>
