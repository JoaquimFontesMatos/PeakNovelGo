<script setup lang="ts">
    import { useWindowSize } from '@vueuse/core';

    const drawerRef = ref(null);
    const drawerHeight = ref(300);
    const isMobile: Ref<boolean> = ref(false);

    const minHeight = 100;
    const maxHeight = ref(400);

    const { width, height } = useWindowSize();

    onMounted(() => {
        isMobile.value = width.value < 500;
        maxHeight.value = height.value * 0.9;
    });

    watch(width, newWidth => {
        isMobile.value = newWidth < 500;
    });

    let startY = 0;
    let startHeight = 0;

    const onPointerDown = (event: PointerEvent) => {
        if (!isMobile.value) return;

        startY = event.clientY;
        startHeight = drawerHeight.value;

        window.addEventListener('pointermove', onPointerMove);
        window.addEventListener('pointerup', onPointerUp);
    };

    const onPointerMove = (event: PointerEvent) => {
        const delta = startY - event.clientY;
        drawerHeight.value = Math.min(maxHeight.value, Math.max(minHeight, startHeight + delta));
    };

    const onPointerUp = () => {
        window.removeEventListener('pointermove', onPointerMove);
        window.removeEventListener('pointerup', onPointerUp);
    };
</script>

<template>
    <Transition :name="isMobile ? 'drawer-up' : 'drawer-right'" :key="isMobile ? 'isMobile' : 'notMobile'" appear>
        <div
            ref="drawerRef"
            :class="[
                'fixed z-50 overflow-hidden shadow-xl backdrop-blur-xs transition-all duration-0',
                isMobile ? 'right-0 bottom-0 left-0 rounded-t-2xl' : 'top-0 right-0 h-full w-[clamp(33.33%,1/3vw,50%)]',
            ]"
            :style="isMobile ? { height: `${drawerHeight}px` } : {}"
        >
            <!-- Drag handle only on mobile -->
            <div v-if="isMobile" class="flex w-full cursor-row-resize touch-none justify-center py-3" @pointerdown.passive="onPointerDown">
                <div class="h-2 w-16 rounded-full bg-accent-gold-dark opacity-50 hover:opacity-100"></div>
            </div>

            <SpacersSmallVertical v-if="!isMobile" />

            <!-- Drawer content slot -->
            <div class="h-full overflow-auto px-4">
                <slot />
            </div>
        </div>
    </Transition>
</template>

<style scoped>
    /* Custom scrollbars (optional) */
    ::-webkit-scrollbar {
        width: 6px;
    }

    ::-webkit-scrollbar-thumb {
        background-color: #d1d5db;
        border-radius: 9999px;
    }
</style>
