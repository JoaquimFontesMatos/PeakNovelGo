<script setup lang="ts">
    const userStore = useUserStore();
    const { user } = storeToRefs(userStore);
    const colorMode = useColorMode();

    const handleChangeReadingPreferences = async () => {
        if (!user.value) return;

        const fields = {
            readingPreferences: JSON.stringify(user.value.readingPreferences),
        };

        if (user.value.readingPreferences.theme !== undefined && user.value.readingPreferences.theme) {
            colorMode.preference = user.value.readingPreferences.theme;
        }

        try {
            await userStore.updateUserFields(fields);
        } catch {}

        userStore.saveUserLocalStorage();
    };
</script>

<template>
    <DisplaySettings />

    <SmallVerticalSpacer />

    <AudioSettings />
</template>
