<script setup>
    const auth = useAuthStore();
    const loginDialog = ref(null);
    const { user } = storeToRefs(auth);

    const handleClick = (event, callback) => {
        if (user.value === null) {
            event.preventDefault();
            event.stopPropagation(); // Fully stop event propagation
            loginDialog.value.showDialog = true;
        } else if (callback) {
            callback(); // Execute the action only if logged in
        }
    };
</script>

<template>
    <div>
        <slot :handleClick="handleClick" />
        <LoginDialog ref="loginDialog" />
    </div>
</template>
