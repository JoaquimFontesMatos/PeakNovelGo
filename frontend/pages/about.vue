<template>
    <Container>
        <div class="about-us-content">
            <h1 class="mb-4 text-3xl font-bold">About PeakNovelGo</h1>

            <p class="mb-4">
                PeakNovelGo is built with a passion for light novels and web development. We aim to create a seamless and enjoyable reading experience for light
                novel enthusiasts. Our goal is to provide a platform where readers can discover new worlds, immerse themselves in captivating stories, and
                connect with fellow fans.
            </p>

            <h2 class="mb-2 text-2xl font-bold">Our Mission</h2>

            <p class="mb-4">
                We strive to offer a comprehensive library of light novels, constantly updated with the latest releases and hidden gems. We believe everyone
                deserves access to high-quality literature, regardless of their location or background. Our platform is designed to be user-friendly, intuitive,
                and accessible across various devices.
            </p>

            <h2 class="mb-2 text-2xl font-bold">The Team</h2>
            <p class="mb-4">
                We are a small team of dedicated developers and light novel lovers working tirelessly to bring you the best possible reading experience.
            </p>

            <h2 class="mb-2 text-2xl font-bold">Contact Us</h2>

            <p class="mb-4">
                Have any questions, suggestions, or feedback? We'd love to hear from you! You can reach us at:
                <a href="mailto:peaknovelgo@gmail.com" class="text-accent-gold hover:underline">peaknovelgo@gmail.com</a>
            </p>

            <p class="mb-4">Follow us on social media:</p>
            <ul class="flex space-x-4">
                <!-- Replace with actual social media links -->
                <li>
                    <a href="#" class="text-accent-gold hover:underline">Twitter</a>
                </li>
                <li>
                    <a href="#" class="text-accent-gold hover:underline">Facebook</a>
                </li>
                <li>
                    <a href="#" class="text-accent-gold hover:underline">Discord</a>
                </li>
            </ul>
        </div>

        <VerticalSpacer />

        <NuxtLink to="/">
            <MainButton>Home</MainButton>
        </NuxtLink>
    </Container>
</template>

<style scoped>
</style>
<script setup lang="ts"></script>
