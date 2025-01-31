<script setup lang="ts">
const toastStore = useToastStore()
const {toasts} = storeToRefs(toastStore)
</script>

<template>
    <div class="fixed top-4 right-4 space-y-2 z-50">
        <div
            v-for="toast in toasts"
            :key="toast.id"
            :class="[
                'w-96 p-4 rounded shadow transition-opacity duration-300',
                toast.type === 'error' ? 'bg-red-500 text-white' : 'bg-green-500 text-white',
            ]"
        >
            <button
                class="float-right"
                type="button"
                @click="toastStore.removeToast(toast.id)"
            >
                <Icon name="fluent:dismiss-24-filled" class="text-primary " :size="'1.5em'"/> <!-- Show Icon -->
            </button>
            <p>{{ toast.message }}</p>
        </div>
    </div>
</template>