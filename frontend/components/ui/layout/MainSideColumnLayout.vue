<script setup lang="ts">
    import Sortable from 'sortablejs';
    import type { LayoutItem } from '~/schemas/LayoutItem';

    const editMode = ref(false);

    const props = defineProps<{
        mainColumn: LayoutItem[];
        sideColumn: LayoutItem[];
    }>();

    const emit = defineEmits<{
        (e: 'update:mainColumn', value: LayoutItem[]): void;
        (e: 'update:sideColumn', value: LayoutItem[]): void;
    }>();

    const mainRef = ref<HTMLElement>() as Ref<HTMLElement>;
    const sideRef = ref<HTMLElement>() as Ref<HTMLElement>;

    onMounted(() => {
        nextTick(() => {
            const group = { name: 'shared', pull: true, put: true };

            if (mainRef.value) {
                Sortable.create(mainRef.value as HTMLElement, {
                    animation: 200,
                    handle: '.drag-handle',
                    group,
                    onEnd: () => syncItems(),
                });
            }

            if (sideRef.value) {
                Sortable.create(sideRef.value as HTMLElement, {
                    animation: 200,
                    handle: '.drag-handle',
                    group,
                    onEnd: () => syncItems(),
                });
            }
        });
    });

    function syncItems() {
        // Grabs items from the DOM and syncs arrays accordingly
        const getItemsFromContainer = (container: HTMLElement | null): LayoutItem[] => {
            if (!container) return [];
            return Array.from(container.children)
                .map((el: any) => {
                    const id = el.getAttribute('data-id');
                    return [...props.mainColumn, ...props.sideColumn].find(item => item.id === id);
                })
                .filter((item): item is LayoutItem => Boolean(item));
        };

        emit('update:mainColumn', getItemsFromContainer(mainRef.value));
        emit('update:sideColumn', getItemsFromContainer(sideRef.value));
    }

    function resolveComponent(type: string) {
        const registry: Record<string, Component> = {
            PaginatedNovelGallery: defineAsyncComponent(() => import('~/components/ui/gallery/paginated/PaginatedNovelGallery.vue')),
            TextBlock: defineAsyncComponent(() => import('~/components/container/TextBlock.vue')),
            ImageBlock: defineAsyncComponent(() => import('~/components/container/ImageBlock.vue')),
            FeaturedNovels: defineAsyncComponent(() => import('~/components/ui/sections/FeaturedNovels.vue')),
            RecentlyUpdatedNovels: defineAsyncComponent(() => import('~/components/ui/sections/RecentlyUpdatedNovels.vue')),
            TopRatedNovels: defineAsyncComponent(() => import('~/components/ui/sections/TopRatedNovels.vue')),
            PopularTags: defineAsyncComponent(() => import('~/components/ui/sections/PopularTags.vue')),
        };
        return registry[type] || 'div';
    }

    function toggleEditMode() {
        editMode.value = !editMode.value;
    }
</script>

<template>
    <CircularButton
        :icon-size="12"
        :no-background="true"
        :icon-name="editMode ? 'fluent:edit-20-filled' : 'fluent:edit-20-regular'"
        :padding="3"
        @click="toggleEditMode()"
    />
    <div class="flex flex-col gap-6 md:flex-row">
        <!-- Main Column -->
        <div ref="mainRef" class="min-h-[100px] flex-1 space-y-4 rounded-sm p-4">
            <div v-for="item in mainColumn" :key="item.id" class="p-4" :data-id="item.id">
                <div v-show="editMode" class="drag-handle cursor-move text-sm text-gray-500">:: drag</div>
                <component :is="resolveComponent(item.type)" v-bind="item.props" />
            </div>
        </div>

        <!-- Side Column -->
        <div ref="sideRef" class="min-h-[100px] w-full space-y-4 rounded-sm bg-primary p-4 md:w-64 md:bg-secondary md:shadow-sm">
            <div v-for="item in sideColumn" :key="item.id" class="rounded-sm bg-primary p-4 shadow-xs" :data-id="item.id">
                <div v-show="editMode" class="drag-handle cursor-move text-sm text-gray-500">:: drag</div>
                <component :is="resolveComponent(item.type)" v-bind="item.props" />
            </div>
        </div>
    </div>
</template>
