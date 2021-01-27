<template>
  <div id="demo">
    <button v-on:click="show = !show" class="bg-blue-200 hover:outline-none">
      Toggle
    </button>
    <transition name="fade">
      <div v-if="show" class="w-32 h-32 bg-black text-white transition-colors">hello</div>
    </transition>
  </div>
</template>

<style>
.fade-enter-active, .fade-leave-active {
  transition: opacity 500ms;
}
.fade-enter, .fade-leave-to {
  opacity: 0;
}
</style>

<script>
export default {
  data() {
    return {
      show: true
    }
  },
  methods: {
    toggle() {
      this.show = !this.show
    }
  }
}
</script>