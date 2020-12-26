<template>
  <h1>WebSocket example</h1>
  <button type="button" @click="newRoom"> new room</button><br><br>
  <input type="text" v-model="roomNumber" placeholder="enter room number">
  <input type="text" v-model="playerName" placeholder="enter player name">
  <button type="button" @click="joinRoom"> join room</button><br><br>
  <button type="button" @click="startGame"> Start game</button>
  <button type="button" @click="endGame"> End game</button>
  <button type="button" @click="newCard"> Get new card</button>
  <button type="button" @click="printRoom"> Print room</button><br><br>
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
        this.socket.send(JSON.stringify({action: "join-room",name: this.playerName, room: this.roomNumber}))
      })

      this.socket.addEventListener("message", (event) => {
        console.log("Received message: ", event.data)
      })
    },
    startGame(){
      this.socket = new WebSocket("ws://localhost:8081/ws")

      this.socket.addEventListener("open", (event) => {
        this.socket.send(JSON.stringify({action: "start-game", room: this.roomNumber}))
      })

      this.socket.addEventListener("message", (event) => {
        console.log("Received message: ", event.data)
      })
    },
    newCard(){
      this.socket = new WebSocket("ws://localhost:8081/ws")

      this.socket.addEventListener("open", (event) => {
        this.socket.send(JSON.stringify({action: "new-card", room: this.roomNumber, name: this.playerName}))
      })

      this.socket.addEventListener("message", (event) => {
        console.log("Received message: ", event.data)
      })
    },
    endGame(){
      this.socket = new WebSocket("ws://localhost:8081/ws")

      this.socket.addEventListener("open", (event) => {
        this.socket.send(JSON.stringify({action: "end-game", room: this.roomNumber}))
      })

      this.socket.addEventListener("message", (event) => {
        console.log("Received message: ", event.data)
      })
    },
    printRoom(){
      this.socket = new WebSocket("ws://localhost:8081/ws")

      this.socket.addEventListener("open", (event) => {
        this.socket.send(JSON.stringify({action: "print-room", room: this.roomNumber}))
      })

      this.socket.addEventListener("message", (event) => {
        console.log("Received message: ", event.data)
      })
    }
  }
}
</script>
