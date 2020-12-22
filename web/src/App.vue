<template>
  <h1>WebSocket example</h1>
  <button type="button" @click="newRoom"> new room</button><br><br>
  <input type="number" v-model="roomNumber" placeholder="enter roomnumber">
  <input type="text" v-model="playerName" placeholder="enter playername">
  <button type="button" @click="joinRoom"> join room</button>
</template>

<script>
export default {
  name: "App",
  data() {
    return {
      greeting: "Hello, world!",
      roomNumber: null,
      playerName: null,
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
      alert("Joining room " + this.roomNumber + " as " + this.playerName)
      this.socket = new WebSocket("ws://localhost:8081/ws")

      this.socket.addEventListener("open", (event) => {
        this.socket.send(JSON.stringify({action: "join-room", id: this.roomNumber, name: this.playerName}))
      })

      this.socket.addEventListener("message", (event) => {
        console.log("Received message: ", event.data)
      })
    }
  }

}
</script>
