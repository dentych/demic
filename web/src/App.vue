<template>
  <h1>WebSocket example</h1>
  <button type="button" @click="newRoom"> new room</button>
  <input type="number" :roomnumber = placeholder="enter roomnumber">
  <input type="text" placeholder="enter playername">
  <button type="button" @click="joinRoom"> join room</button>
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
  },
  methods:{
    newRoom(){
      this.socket = new WebSocket("ws://localhost:8081/ws")

      this.socket.addEventListener("open", (event) => {
        this.socket.send(JSON.stringify({action: "new-room"}))
      })

      this.socket.addEventListener("message", (event) => {
        console.log("Received message: ", event.data)
      })
    },
    joinRoom(){
      this.socket = new WebSocket("ws://localhost:8081/ws")

      this.socket.addEventListener("open", (event) => {
        this.socket.send(JSON.stringify({action: "join-room",id: "",}))
      })

      this.socket.addEventListener("message", (event) => {
        console.log("Received message: ", event.data)
      })
    }
  }

}
</script>
