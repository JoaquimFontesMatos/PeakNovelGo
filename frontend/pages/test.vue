<script setup lang="ts">
    import type { LayoutItem } from '~/schemas/LayoutItem';

    const currentPage = ref(1);
    const currentLimit = ref(10);

    useHead({
        title: '📖 Home',
    });

    onMounted(async () => {
        await onPageChange(currentPage.value, currentLimit.value);
        // featuredNovels.value = await fetchFeaturedNovels();

        // Fetch recently updated novels (replace with your actual API call)
        // recentlyUpdated.value = await fetchRecentlyUpdatedNovels();

        // Fetch popular tags (replace with your actual API call)
        popularTags.value = await fetchPopularTags();

        // topRatedNovels.value = await fetchTopRatedNovels();

        loadLayout();
    });

    const { fetchingNovel, paginatedNovelsData } = storeToRefs(useNovelStore());

    const onPageChange = async (newPage: number, limit: number) => {
        if (newPage === currentPage.value && limit === currentLimit.value && paginatedNovelsData.value && paginatedNovelsData.value.page != 0) return;

        try {
            await useNovelStore().fetchNovels(newPage, limit);
        } catch {}
        currentPage.value = newPage;
        currentLimit.value = limit;
    };

    //import HeroSection from '~/components/HeroSection.vue';  // Import the hero section

    const featuredNovels = ref([]);
    const recentlyUpdated = ref([]);
    const popularTags: Ref<string[]> = ref([]);
    const topRatedNovels = ref([]);

    // Placeholder functions - replace these with your actual API calls
    async function fetchFeaturedNovels() {
        return [];
    }

    async function fetchRecentlyUpdatedNovels() {
        return [];
    }

    async function fetchTopRatedNovels() {
        return [];
    }

    async function fetchPopularTags(): Promise<string[]> {
        return [
            'Adapted to Manhua',
            'Transformation Ability',
            'Transmigration',
            'Transplanted Memories',
            'Vampires',
            'Weak to Strong',
            'Wealthy Characters',
            'Werebeasts',
            'Zombies',
        ];
    }

    const mainItems: Ref<LayoutItem[]> = ref([
        {
            id: '1',
            type: 'PaginatedNovelGallery',
            props: {
                errorMessage: paginatedNovelsData === null ? 'No Novels Found' : null,
                paginatedData: paginatedNovelsData,
                onPageChange: onPageChange,
            },
        },
        {
            id: '2',
            type: 'PaginatedNovelGallery',
            props: {
                errorMessage: paginatedNovelsData === null ? 'No Novels Found' : null,
                paginatedData: paginatedNovelsData,
                onPageChange: onPageChange,
            },
        },
        {
            id: '3',
            type: 'PaginatedNovelGallery',
            props: {
                errorMessage: paginatedNovelsData === null ? 'No Novels Found' : null,
                paginatedData: paginatedNovelsData,
                onPageChange: onPageChange,
            },
        },
    ]);

    const sideItems: Ref<LayoutItem[]> = ref([
        {
            id: '4',
            type: 'TextBlock',
        },
        {
            id: '5',
            type: 'ImageBlock',
        },
    ]);

    function saveLayout() {
        const strip = (items: LayoutItem[]) =>
            items.map(({ id, type, props }) => ({
                id,
                type,
                props: {
                    ...props,
                    onPageChange: undefined, // remove function before saving
                },
            }));

        localStorage.setItem(
            'layout',
            JSON.stringify({
                main: strip(mainItems.value),
                side: strip(sideItems.value),
            })
        );
    }

    function loadLayout() {
        const saved = JSON.parse(localStorage.getItem('layout') || 'null');
        if (!saved) return;

        const inject = (items: LayoutItem[]) =>
            items.map(item => {
                if (item.type === 'PaginatedNovelGallery') {
                    return {
                        ...item,
                        props: {
                            ...item.props,
                            onPageChange,
                        },
                    };
                }
                return item;
            });

        mainItems.value = inject(saved.main);
        sideItems.value = inject(saved.side);
    }
</script>

<template>
    <div class="home-page">
        <Container>
            <MainSideColumnLayout
                :main-column="mainItems"
                :side-column="sideItems"
                @update:mainColumn="
                    val => {
                        mainItems = val;
                        saveLayout();
                    }
                "
                @update:sideColumn="
                    val => {
                        sideItems = val;
                        saveLayout();
                    }
                "
            />
        </Container>
    </div>
</template>
