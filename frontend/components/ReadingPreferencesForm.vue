<script setup lang="ts">
const userStore = useUserStore()

const {user} = storeToRefs(userStore)

const handleChangeReadingPreferences = async() => {
    if (user.value) {
        const fields = {
            readingPreferences: JSON.stringify(user.value.readingPreferences),
        }

        await userStore.updateUserFields(fields)

        userStore.saveUserLocalStorage()
    }
}
</script>

<template>
    <div v-if="user" class="form-container space-y-6 p-4 bg-secondary rounded-lg shadow-lg">
        <fieldset class="border-t border-accent-gold-dark pt-4">
            <legend class="text-xl font-semibold text-primary-content ml-3.5 px-3.5">Reading Preferences</legend>
            <!-- Atomic Reading -->
            <div class="form-group flex items-center space-x-2">
                <input
                    id="atomicReading"
                    name="atomicReading"
                    type="checkbox"
                    v-model="user.readingPreferences.atomicReading"
                    @change="handleChangeReadingPreferences()"
                />
                <label for="atomicReading" class="text-sm font-medium text-secondary-content">Enable Atomic Reading</label>
            </div>

            <!-- Font Selector -->
            <div class="form-group">
                <label for="fontPreview" class="block text-sm font-medium text-secondary-content">Font</label>
                <select
                    id="fontPreview"
                    name="font"
                    v-model="user.readingPreferences.font"
                    @change="handleChangeReadingPreferences()"
                >
                    <option value="font-serif">Serif</option>
                    <option value="font-sans">Sans Serif</option>
                    <option value="font-mono">Monospace</option>
                    <option value="font-montserrat">Montserrat</option>
                    <option value="font-noto">Noto Sans</option>
                    <option value="font-raleway">Raleway</option>
                </select>
                <p class="mt-2" :class="user.readingPreferences.font">
                    Preview: The quick brown fox jumps over the lazy dog.
                </p>
            </div>
        </fieldset>

        <!-- Theme Selector -->
        <div class="form-group">
            <label for="theme" class="block text-sm font-medium text-secondary-content">Theme</label>
            <select
                id="theme"
                name="theme"
                v-model="user.readingPreferences.theme"
                @change="handleChangeReadingPreferences()"
            >
                <option value="Unfollow">Unfollow</option>
                <option value="Reading">Reading</option>
                <option value="Completed">Completed</option>
                <option value="On-Hold">On-Hold</option>
                <option value="Dropped">Dropped</option>
                <option value="Plan to Read">Plan to Read</option>
            </select>
        </div>

        <!-- Text-to-Speech Section -->
        <fieldset class="border-t border-accent-gold-dark pt-4">
            <legend class="text-lg font-semibold text-primary-content ml-3.5 px-3.5">Text-to-Speech</legend>
            <div class="mt-2 space-y-4">
                <!-- Autoplay -->
                <div class="form-group flex items-center space-x-2">
                    <input
                        id="autoplay"
                        name="autoplay"
                        type="checkbox"
                        v-model="user.readingPreferences.tts.autoplay"
                        @change="handleChangeReadingPreferences()"
                    />
                    <label for="autoplay" class="text-sm font-medium text-secondary-content">Autoplay</label>
                </div>

                <!-- Voice Selector -->
                <div class="form-group">
                    <label for="voice" class="block text-sm font-medium text-secondary-content">Voice</label>
                    <select
                        id="voice"
                        name="voice"
                        v-model="user.readingPreferences.tts.voice"
                        @change="handleChangeReadingPreferences()"
                    >
                        <option value="en-US-AriaNeural">"en-US AriaNeural"</option>
                        <option value="Reading">Reading</option>
                        <option value="Completed">Completed</option>
                        <option value="On-Hold">On-Hold</option>
                        <option value="Dropped">Dropped</option>
                        <option value="Plan to Read">Plan to Read</option>
                    </select>
                </div>

                <!-- Speed Input -->
                <div class="form-group">
                    <label for="speed" class="block text-sm font-medium text-secondary-content">Speed</label>
                    <input
                        id="speed"
                        name="speed"
                        type="number"
                        step="0.1"
                        min="0.1"
                        max="3"
                        v-model="user.readingPreferences.tts.speed"
                        @change="handleChangeReadingPreferences()"
                    />
                </div>
            </div>
        </fieldset>
    </div>
</template>
