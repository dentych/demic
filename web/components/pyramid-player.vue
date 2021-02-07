<template>
  <div class="flex flex-col h-screen">
    <div class="bg-orange-dark text-center mb-6">
      <span class="text-sm">ROOM CODE: {{ code }}</span>
    </div>

    <!-- Host buttons -->
    <div v-if="host" class="mb-8">
      <div class="text-center mb-4">
        GAME MASTER ACTIONS
      </div>
      <div class="flex space-x-4 px-2">
        <button class="flex-auto focus:outline-none bg-yellow-700 text-yellow-100 text-lg py-4"
                v-if="!started"
                @click="startGame">
          Start game
        </button>
        <button class="flex-auto focus:outline-none bg-yellow-700 text-yellow-100 text-lg py-4"
                v-if="started"
                @click="continueGame">
          Turn next card
        </button>
      </div>
    </div>

    <div>
      <div class="text-center mb-4">
        PLAYER ACTIONS
      </div>
      <div class="flex space-x-4 px-2">
        <button class="flex-auto focus:outline-none bg-yellow-700 text-yellow-100 text-lg">Hejsa</button>
        <button class="flex-auto focus:outline-none bg-yellow-700 text-yellow-100 text-lg">Hejsa</button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.bg-orange-dark {
  background-color: #ffc58c;
}
</style>

<script>
export default {
  props: ["name", "code"],
  data() {
    return {
      ws: null,
      host: false,
      players: [],
      started: false
    }
  },
  methods: {
    messageHandler(evt) {
      let data = JSON.parse(evt.data)
      console.log(data.action)
      switch (data.action.action_type) {
        case "player-join":
          this.players.push(data.action.target)
          break
        case "player-left":
          let i = this.players.indexOf(data.action.target)
          this.players.splice(i, 1)
          break
        case "host":
          this.host = data.action.target === this.name
          break
        case "start-game":
          this.started = true
          break
      }
    },
    startGame() {
      this.ws.send(JSON.stringify({room_id: this.code, action: {action_type: "start-game", origin: this.name}}))
    },
    continueGame() {
      this.ws.send(JSON.stringify({room_id: this.code, action: {action_type: "continue", origin: this.name}}))
    }
  },
  mounted() {
    this.ws = new WebSocket("ws://" + location.hostname + ":8080/ws")
    this.ws.onopen = () => {
      this.ws.send(JSON.stringify({room_id: this.code, action: {action_type: "player-join", origin: this.name}}))
    }

    this.ws.onmessage = this.messageHandler
  }
}
</script>
