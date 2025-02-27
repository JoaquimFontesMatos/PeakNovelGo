import { ToastIconMap, type Toast, type ToastIcon, type ToastType } from '~/schemas/Toast';

export const useToastStore = defineStore('Toast', () => {
  const toasts = ref<Toast[]>([]);
  const toastId = ref(1);

  const addToast = (message: string, type: ToastType, icon?: ToastIcon) => {
    const id = ++toastId.value;

    const iconValue = icon && ToastIconMap[icon] ? ToastIconMap[icon] : ToastIconMap['none'];

    toasts.value.push({ id, message, type, icon: iconValue });

    // Remove the toast after 3 seconds
    setTimeout(() => {
      removeToast(id);
    }, 3000);
  };

  // Remove a toast by ID
  const removeToast = (id: number) => {
    const index = toasts.value.findIndex(toast => toast.id === id);
    if (index !== -1) toasts.value.splice(index, 1);
  };

  return { toasts, addToast, removeToast };
});

if (import.meta.hot) {
  import.meta.hot.accept(acceptHMRUpdate(useToastStore, import.meta.hot));
}
