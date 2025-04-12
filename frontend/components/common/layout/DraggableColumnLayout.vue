<script setup lang="ts">
    import Sortable from 'sortablejs';
    import type { LayoutItem } from '~/schemas/LayoutItem';

    const props = defineProps<{
        modelValue: LayoutItem[];
    }>();

    const emit = defineEmits<{
        (e: 'update:modelValue', value: LayoutItem[]): void;
    }>();

    const containerRef = ref<HTMLElement | null>(null);

    onMounted(() => {
        nextTick(() => {
            Sortable.create(containerRef.value as HTMLElement, {
                animation: 200,
                handle: '.drag-handle',
                onEnd: evt => {
                    const oldIndex = evt.oldIndex ?? -1;
                    const newIndex = evt.newIndex ?? -1;

                    if (oldIndex !== -1 && newIndex !== -1) {
                        const newItems = [...props.modelValue];
                        const moved = newItems.splice(oldIndex, 1)[0];
                        newItems.splice(newIndex, 0, moved);
                        emit('update:modelValue', newItems);
                    }
                },
            });
        });
    });

    /**
     * Example dynamic mapping — customize this to match your app
     */
    function resolveComponent(type: string) {
        const registry: Record<string, any> = {
            'text-block': 'TextBlock',
            'image-block': 'ImageBlock',
            'paginated-novel-gallery': 'PaginatedNovelGallery',
        };
        return registry[type] || 'div';
    }
</script>

<template>
    <div ref="containerRef" class="space-y-4">
        <div v-for="item in modelValue" :key="item.id" class="p-4">
            <div class="drag-handle mb-2 cursor-move text-sm text-gray-500">:: drag</div>

            <!-- Dynamically resolve and bind props -->
            <component :is="resolveComponent(item.type)" v-bind="item.props" />
        </div>
    </div>
</template>
