<template>
  <div>
    <PyramidPlayer v-if="name" :name="name" :code="code"></PyramidPlayer>
    <PyramidMenu v-else @create="create" @join="join"></PyramidMenu>
  </div>
</template>

<script>
import PyramidMenu from "~/components/pyramid-menu";
import PyramidPlayer from "~/components/pyramid-player";

export default {
  components: {
    PyramidPlayer,
    PyramidMenu
  },
  data() {
    return {
      name: null,
      code: null
    }
  },
  methods: {
    create() {
      this.$router.push("/pyramid/host")
    },
    join(evt) {
      this.name = evt.name
      this.code = evt.code
    }
  }
}
</script>

<style scoped>

</style>