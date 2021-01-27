<template>
  <div>
    <div class="flex justify-center m-4" v-if="!ws">
      <div class="border rounded flex-auto sm:flex-none sm:w-1/4 bg-blue-100">
        <div class="flex flex-col p-5">
          <h1 class="place-self-center text-3xl">Pyramid</h1>
          <p class="font-medium text-lg mb-2">Room code:</p>
          <input class="mb-2 py-1 px-2 border rounded focus:outline-none uppercase" type="text" v-model="code"
                 placeholder="Enter a valid room code" @keyup="updateCode">
          <p class="font-medium text-lg mb-2">Name:</p>
          <input class="mb-2 py-1 px-2 border rounded focus:outline-none uppercase" type="text" v-model="name"
                 placeholder="Enter a name" @keyup="updateName">
          <button class="mb-2 py-2 text-white bg-blue-600 mx-2" @click="join">Join room</button>
        </div>
        <hr>
        <div class="flex flex-col p-5">
          <p class="font-medium text-lg mb-2">Create a new room:</p>
          <button class="mb-6 py-2 text-white bg-blue-600 mx-2" @click="create">Create new room</button>
        </div>
      </div>
    </div>
    <div class="flex justify-center" v-if="ws">
      <div class="block">
        <p v-if="code">Room code: {{ code }}</p>
        <p class="font-bold">Players</p>
        <ul class="font-bold">
          <li v-for="player in players">{{ player }}</li>
        </ul>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      code: "",
      name: "",
      ws: null,
      players: []
    }
  },
  methods: {
    create() {
      this.ws = new WebSocket("ws://" + location.hostname + ":8080/ws")
      this.ws.onopen = () => {
        this.ws.send(JSON.stringify({action: {action_type: "create-game"}}))
      }

      this.ws.onmessage = this.messageHandler

      this.ws.onerror = err => {
        console.log(err)
      }
    },
    join() {
      this.ws = new WebSocket("ws://" + location.hostname + ":8080/ws")
      this.ws.onopen = () => {
        this.ws.send(JSON.stringify({room_id: this.code, action: {action_type: "player-join", origin: this.name}}))
      }

      this.ws.onmessage = this.messageHandler
    },
    updateCode() {
      localStorage.setItem("code", this.code)
    },
    updateName() {
      localStorage.setItem("name", this.name)
    },
    messageHandler(evt) {
      let data = JSON.parse(evt.data)
      switch (data.action.action_type) {
        case "game-created":
          console.log("pyramid created with ID", data.action.target)
          this.code = data.action.target
          this.updateName()
          this.updateCode()
          break
        case "player-joined":
          this.players.push(data.action.target)
          break
        case "player-left":
          let i = this.players.indexOf(data.action.target)
          this.players.splice(i, 1)
          break
      }
    }
  },
  mounted() {
    this.code = localStorage.getItem("code")
    this.name = localStorage.getItem("name")
  }
}
</script>

<style scoped>

</style>