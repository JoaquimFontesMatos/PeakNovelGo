<script setup lang="ts">
    import { computed, onMounted, onUnmounted, ref } from 'vue';

    const cloudCount = 18;
    const dragonProps = ref({
        animationDelay: 0,
        animationDuration: 8,
        scale: 1,
        zIndex: 0,
    });

    // Randomize dragon properties periodically
    const randomizeDragonProps = () => {
        dragonProps.value = {
            animationDelay: Math.random() * 5,
            // Delay between 0 and 5 seconds
            animationDuration: 8 + Math.random() * 5,
            // Duration between 8 and 13 seconds
            scale: 0.8 + Math.random() * 0.5,
            // Scale between 0.8 to 1.3
            zIndex: Math.floor(Math.random() * 3),
            // Random z-index between 0 and 2
        };
    };

    let intervalId: ReturnType<typeof setTimeout>;
    onMounted(() => {
        // Re-randomize dragon properties every 10 seconds
        intervalId = setInterval(randomizeDragonProps, dragonProps.value.animationDuration * 1000 + dragonProps.value.animationDelay * 1000);
        randomizeDragonProps(); // Initial randomization
    });

    onUnmounted(() => {
        clearInterval(intervalId);
    });

    const getCloudImage = (i: number) => {
        const images = ['/img/cloud_01.svg', '/img/cloud_02.svg', '/img/cloud_03.svg', '/img/cloud_04.svg'];
        return images[i % images.length];
    };

    // Generate cloud delays dynamically
    const cloudDelays = computed(() => Array.from({ length: cloudCount }, (_, i) => 0.45 * i + Math.random() * 0.1));
</script>

<template>
    <div class="relative flex h-16 items-center overflow-hidden">
        <!-- Clouds -->
        <img
            v-for="i in cloudCount"
            :key="`cloud-$ {i}`"
            class="cloud h-12"
            :style="{
                animationDelay: `${cloudDelays[i - 1]}s`,
                zIndex: Math.floor(i / 6),
            }"
            :src="getCloudImage(i)"
            alt="cloud"
        />

        <!-- Dragon -->
        <img
            src="/img/dragon.png"
            class="dragon"
            :style="{
                animationDelay: `${dragonProps.animationDelay}s`,
                animationDuration: `${dragonProps.animationDuration}s`,
                transform: `scale(${dragonProps.scale})`,
                zIndex: dragonProps.zIndex,
            }"
            alt=""
        />
    </div>
</template>

<style scoped>
    /* Keyframes for dynamic movement */
    @keyframes moveAndBounce {
        0% {
            opacity: 0;
            left: -10%; /* Start outside the left edge */
            transform: translateY(0px);
        }
        10% {
            opacity: 0.25;
        }
        25% {
            opacity: 0.75;
            left: 20%;
            transform: translateY(-15%);
        }
        50% {
            opacity: 1;
            left: 50%;
            transform: translateY(0px);
        }
        75% {
            opacity: 0.75;
            left: 80%;
            transform: translateY(-15%);
        }
        90% {
            opacity: 0.25;
        }
        100% {
            opacity: 0;
            left: 110%; /* End outside the right edge */
            transform: translateY(0px);
        }
    }

    .cloud {
        opacity: 0;
        animation-duration: calc(8s + 4 * var(--index, 0));
        position: absolute;
        transform: scale(0.8);
        animation: moveAndBounce 10s linear infinite;
    }

    .cloud:nth-child(odd) {
        transform: scale(1.1);
    }

    .cloud:nth-child(even) {
        transform: scale(0.9);
    }

    .cloud:hover {
        filter: brightness(1.5) contrast(1.2);
        transition: filter 0.3s ease;
    }

    @keyframes dragonFlight {
        0% {
            left: -20%; /* Start off-screen */
            transform: translateY(0) scale(1);
            opacity: 0;
        }
        25% {
            transform: translateY(-20%) scale(1.1);
        }
        50% {
            left: 50%; /* Mid-screen */
            transform: translateY(10%) scale(1.2);
            opacity: 1;
        }
        75% {
            transform: translateY(-10%) scale(1.1);
        }
        100% {
            left: 120%; /* End off-screen */
            transform: translateY(0) scale(1);
            opacity: 0.5;
        }
    }

    @keyframes fall {
        0% {
            opacity: 1;
            transform: translateY(0);
        }
        100% {
            opacity: 0; /* Fade out as it falls */
            transform: translateY(100px); /* You can add a translateY to make it fall further */
        }
    }

    .dragon-container {
        position: relative;
        width: 100%;
        height: 100%;
    }

    .dragon {
        opacity: 0;
        position: absolute;
        height: 3rem;
        animation: dragonFlight 10s linear infinite;
        transition:
            transform 0.3s ease,
            filter 0.3s ease;
    }

    .dragon:hover {
        animation-play-state: paused;
        filter: brightness(1.5) contrast(1.2);
        transition: filter 0.3s ease;
    }

    .background {
        background: linear-gradient(to bottom, #f0f8ff, #add8e6);
        animation: backgroundMove 20s linear infinite;
    }

    @keyframes backgroundMove {
        0% {
            background-position: 0% 0%;
        }
        100% {
            background-position: 100% 100%;
        }
    }
</style>
