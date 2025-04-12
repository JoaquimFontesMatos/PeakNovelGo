<script setup lang="ts">
    import Sortable from 'sortablejs';
    import type { LayoutItem } from '~/schemas/LayoutItem';

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
                    return [...props.mainColumn, ...props.sideColumn].find(item => item.id === id)!;
                })
                .filter(Boolean);
        };

        emit('update:mainColumn', getItemsFromContainer(mainRef.value));
        emit('update:sideColumn', getItemsFromContainer(sideRef.value));
    }

    function resolveComponent(type: string) {
        const registry: Record<string, Component> = {
            PaginatedNovelGallery: defineAsyncComponent(() => import('~/components/common/gallery/paginated/PaginatedNovelGallery.vue')),
            TextBlock: defineAsyncComponent(() => import('~/components/container/TextBlock.vue')),
            ImageBlock: defineAsyncComponent(() => import('~/components/container/ImageBlock.vue')),
        };
        return registry[type] || 'div';
    }
</script>

<template>
    <div class="flex gap-6">
        <!-- Main Column -->
        <div ref="mainRef" class="min-h-[100px] flex-1 space-y-4 rounded-sm p-4">
            <div v-for="item in mainColumn" :key="item.id" class="p-4" :data-id="item.id">
                <div class="drag-handle mb-2 cursor-move text-sm text-gray-500">:: drag</div>
                <component :is="resolveComponent(item.type)" v-bind="item.props" />
            </div>
        </div>

        <!-- Side Column -->
        <div ref="sideRef" class="min-h-[100px] w-64 space-y-4 rounded-sm bg-secondary p-4 shadow-sm">
            <div v-for="item in sideColumn" :key="item.id" class="rounded-sm bg-primary p-4 shadow-xs" :data-id="item.id">
                <div class="drag-handle mb-2 cursor-move text-sm text-gray-500">:: drag</div>
                <component :is="resolveComponent(item.type)" v-bind="item.props" />
            </div>
        </div>
    </div>
</template>
