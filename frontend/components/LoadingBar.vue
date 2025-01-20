<template>
  <div class="relative h-40 overflow-hidden">
    <img
      id="clouds1"
      class="top-0 left-0 h-12 cloud"
      src="@img/cloud_01.svg"
      alt="clouds"
    />
    <img
      id="clouds2"
      class="top-5 left-5 h-12 cloud"
      src="@img/cloud_02.svg"
      alt="clouds"
    />
    <img
      id="clouds3"
      class="top-10 left-10 h-12 cloud"
      src="@img/cloud_03.svg"
      alt="clouds"
    />
    <img
      id="clouds4"
      class="top-2 left-2 h-12 cloud"
      src="@img/cloud_04.svg"
      alt="clouds"
    />
    <img
      id="clouds5"
      class="top-0 left-3 h-12 cloud hidden md:block"
      src="@img/cloud_01.svg"
      alt="clouds"
    />
    <img
      id="clouds6"
      class="top-4 left-5 h-12 cloud hidden md:block"
      src="@img/cloud_02.svg"
      alt="clouds"
    />
    <img
      id="clouds7"
      class="top-3 left-9 h-12 cloud hidden md:block"
      src="@img/cloud_03.svg"
      alt="clouds"
    />
    <img
      id="clouds8 "
      class="top-3 left-2 h-12 cloud hidden md:block"
      src="@img/cloud_04.svg"
      alt="clouds"
    />
    <img
      id="clouds9"
      class="top-2 left-2 h-12 cloud"
      src="@img/cloud_04.svg"
      alt="clouds"
    />
    <img
      id="clouds10"
      class="top-0 left-3 h-12 cloud hidden md:block"
      src="@img/cloud_01.svg"
      alt="clouds"
    />
    <img
      id="clouds11"
      class="top-4 left-5 h-12 cloud hidden md:block"
      src="@img/cloud_02.svg"
      alt="clouds"
    />
    <img
      id="clouds12"
      class="top-3 left-9 h-12 cloud hidden md:block"
      src="@img/cloud_03.svg"
      alt="clouds"
    />
    <img
      id="clouds13 "
      class="top-3 left-2 h-12 cloud hidden md:block"
      src="@img/cloud_04.svg"
      alt="clouds"
    />
    <img
      id="clouds14"
      class="top-2 left-2 h-12 cloud"
      src="@img/cloud_04.svg"
      alt="clouds"
    />
    <img
      id="clouds15"
      class="top-0 left-3 h-12 cloud hidden md:block"
      src="@img/cloud_01.svg"
      alt="clouds"
    />
    <img
      id="clouds16"
      class="top-4 left-5 h-12 cloud hidden md:block"
      src="@img/cloud_02.svg"
      alt="clouds"
    />
    <img
      id="clouds17"
      class="top-3 left-9 h-12 cloud hidden md:block"
      src="@img/cloud_03.svg"
      alt="clouds"
    />
    <img
      id="clouds18 "
      class="top-3 left-2 h-12 cloud hidden md:block"
      src="@img/cloud_04.svg"
      alt="clouds"
    />
  </div>
</template>

<style scoped>
/* Individual delays for a staggered effect */
#clouds1 {
  animation-delay: 0.25s;
}
#clouds2 {
  animation-delay: 1.23s;
}
#clouds3 {
  animation-delay: 2.44s;
}
#clouds4 {
  animation-delay: 3.31s;
}
#clouds5 {
  animation-delay: 4.34s;
}
#clouds6 {
  animation-delay: 5.64s;
}
#clouds7 {
  animation-delay: 6.11s;
}
#clouds8 {
  animation-delay: 7.24s;
}
#clouds9 {
  animation-delay: 8.34s;
}
#clouds10 {
  animation-delay: 9.31s;
}
#clouds11 {
  animation-delay: 10.23s;
}
#clouds12 {
  animation-delay: 11.44s;
}
#clouds13 {
  animation-delay: 12.31s;
}
#clouds14 {
  animation-delay: 13.34s;
}
#clouds15 {
  animation-delay: 14.64s;
}
#clouds16 {
  animation-delay: 15.11s;
}
#clouds17 {
  animation-delay: 16.24s;
}
#clouds18 {
  animation-delay: 17.34s;
}

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
  position: absolute;
  animation: moveAndBounce 10s linear infinite;
}
</style>
