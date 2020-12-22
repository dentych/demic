<template>
  <h1>WebSocket example</h1>
  <p>{{ greeting }}</p>
</template>

<script>
export default {
  name: "App",
  data() {
    return {
      greeting: "Hello, world!",
      socket: null
    }
  },
  mounted() {
    this.socket = new WebSocket("ws://localhost:8081/ws")

    this.socket.addEventListener("open", (event) => {
      this.socket.send("hello backend")
    })

    this.socket.addEventListener("message", (event) => {
      console.log("Received message: ", event.data)
    })
  }
}
</script>

<style scoped>

</style>