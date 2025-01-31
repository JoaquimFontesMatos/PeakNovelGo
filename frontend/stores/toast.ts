export const useToastStore = defineStore('Toast', () => {
    const toasts = ref<{
        id: number; message: string; type: 'error' | 'success'
    }[]>([]);
    const toastId = ref(1);

    const addToast = (message: string, type: 'error' | 'success') => {
        const id = ++toastId.value;
        toasts.value.push({ id, message, type });

        // Remove the toast after 3 seconds
        setTimeout(() => {
            removeToast(id);
        }, 3000);
    };

    // Remove a toast by ID
    const removeToast = (id: number) => {
        const index = toasts.value.findIndex((toast) => toast.id === id);
        if (index !== -1) toasts.value.splice(index, 1);
    };

    return { toasts, addToast, removeToast };
});


if (import.meta.hot) {
    import.meta.hot.accept(acceptHMRUpdate(useToastStore, import.meta.hot));
}