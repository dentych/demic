<template>
  <h1>WebSocket example</h1>
  <p>{{ greeting }}</p>
</template>

<script>
export default {
  name: "App",
  data() {
    return {
      greeting: "Hello, world!"
    }
  },
  mounted() {
    this.socket = new WebSocket("ws://localhost:8081")

    this.socket.addEventListener("open", function (event) {
      this.socket.send("hello backend")
    })

    this.socket.addEventListener("message", function (event) {
      console.log("Received message: ", event.data)
    })
  }
}
</script>

<style scoped>

</style>