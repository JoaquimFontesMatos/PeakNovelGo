<script setup lang="ts">
const toastStore = useToastStore()
const {toasts} = storeToRefs(toastStore)
</script>

<template>
    <div class="fixed top-2 right-2 w-full md:top-4 md:right-4 space-y-2 z-50">
        <div
            v-for="toast in toasts"
            :key="toast.id"
            :class="[
                'float-right w-2/3 md:w-96 p-2 md:p-4 rounded shadow transition-opacity duration-300',
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