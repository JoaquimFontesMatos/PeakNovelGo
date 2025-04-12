<template>
    <div class="relative">
        <img id="clouds1" class="cloud absolute top-0 left-0 h-12" src="/img/cloud_01.svg" alt="clouds" />
        <img id="clouds2" class="cloud absolute top-5 left-5 h-12" src="/img/cloud_02.svg" alt="clouds" />
        <img id="clouds3" class="cloud absolute top-10 left-10 h-12" src="/img/cloud_03.svg" alt="clouds" />
        <img id="clouds4" class="cloud absolute top-2 left-2 h-12" src="/img/cloud_04.svg" alt="clouds" />
    </div>
</template>

<style scoped>
    #clouds1 {
        animation-delay: 0s, 0s;
    }

    #clouds2 {
        animation-delay: 0.2s, 0.2s;
    }

    #clouds3 {
        animation-delay: 0.4s, 0.4s;
    }

    #clouds4 {
        animation-delay: 0.6s, 0.6s;
    }

    @keyframes moveClouds {
        0% {
            transform: translateX(0px);
        }
        100% {
            transform: translateX(200px);
        }
    }

    @keyframes bounce {
        0%,
        100% {
            transform: translateY(-25%);
            animation-timing-function: cubic-bezier(0.8, 0, 1, 1);
        }
        50% {
            transform: translateY(0px);
            animation-timing-function: cubic-bezier(0, 0, 0.2, 1);
        }
    }

    .cloud {
        animation:
            moveClouds 10s ease-in-out alternate infinite,
            bounce 1s ease-in-out infinite;
    }
</style>
