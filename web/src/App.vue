<template>
  <h1>WebSocket example</h1>
  <button type="button" @click="newRoom"> New room</button><br><br>
  <input type="text" v-model="roomNumber" placeholder="Enter room number">
  <input type="text" v-model="playerName" placeholder="Enter player name">
  <button type="button" @click="joinRoom"> Join room</button>
  <button type="button" @click="removePlayer"> Remove player</button><br><br>
  <button type="button" @click="startGame"> Start game</button><br><br>
  Hand index (0-3):<input type="text" v-model="handIDX" placeholder="Enter Hand Index">
  <button type="button" @click="newCard"> Get new card</button><br><br>
  <button type="button" @click="endGame"> End game</button><br><br>
  <button type="button" @click="deleteRoom"> Delete Room</button><br><br>
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
      handIDX: null,
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
      // alert("Joining room " + this.roomNumber + " as " + this.playerName)
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
        this.socket.send(JSON.stringify({action: "new-card", room: this.roomNumber, name: this.playerName, card: this.handIDX}))
      })

      this.socket.addEventListener("message", (event) => {
        console.log("Received message: ", event.data)
      })
    },
    endGame(){
      this.socket = new WebSocket("ws://localhost:8081/ws")

      this.socket.addEventListener("open", (event) => {
        this.socket.send(JSON.stringify({action: "end-game", rgoom: this.roomNumber}))
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
    },
    deleteRoom(){
      this.socket = new WebSocket("ws://localhost:8081/ws")

      this.socket.addEventListener("open", (event) => {
        this.socket.send(JSON.stringify({action: "delete-room", room: this.roomNumber}))
      })

      this.socket.addEventListener("message", (event) => {
        console.log("Received message: ", event.data)
      })
    },
    removePlayer(){
      this.socket = new WebSocket("ws://localhost:8081/ws")

      this.socket.addEventListener("open", (event) => {
        this.socket.send(JSON.stringify({action: "remove-player",name: this.playerName, room: this.roomNumber}))
      })

      this.socket.addEventListener("message", (event) => {
        console.log("Received message: ", event.data)
      })
    }
  }
}
</script>
