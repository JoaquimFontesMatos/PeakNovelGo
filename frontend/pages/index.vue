<template>
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
</template>

<script setup lang="ts">
    import type { LayoutItem } from '~/schemas/LayoutItem';

    useHead({
        title: '📖 Home',
    });

    onMounted(() => {
        loadLayout();
    });

    const mainItems: Ref<LayoutItem[]> = ref([
        {
            id: '1',
            type: 'RecentlyUpdatedNovels',
            props: {},
        },
        {
            id: '2',
            type: 'FeaturedNovels',
            props: {},
        },
    ]);

    const sideItems: Ref<LayoutItem[]> = ref([
        {
            id: '3',
            type: 'TopRatedNovels',
        },
        {
            id: '4',
            type: 'PopularTags',
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
                return item;
            });

        mainItems.value = inject(saved.main);
        sideItems.value = inject(saved.side);
    }
</script>
