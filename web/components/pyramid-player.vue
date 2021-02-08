<template>
  <div class="flex flex-col h-screen">
    <div class="bg-orange-dark text-center mb-6">
      <span class="text-sm">ROOM CODE: {{ code }}</span>
    </div>

    <!-- Host actions -->
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
        <button class="flex-auto focus:outline-none bg-yellow-700 text-yellow-100 text-lg py-4 butte"
                v-if="started"
                @click="continueGame"
                :disabled="!first && !attackMode">
          <span v-if="first">Turn first card</span>
          <span v-else>Turn next card</span>
        </button>
      </div>
    </div>

    <!-- Player actions -->
    <div class="flex-grow">
      <div class="text-center mb-4">
        PLAYER ACTIONS
      </div>
      <div class="flex space-x-4 px-2">
        <button class="flex-auto focus:outline-none bg-yellow-700 text-yellow-100 text-lg p-4 butte"
                :disabled="!attackMode">
          I have the card!
        </button>
      </div>
    </div>

    <!-- Player cards -->
    <div class="flex" v-if="cards">
      <div class="flex-auto" v-for="card in cards">
        <img :src="'/playing-cards/' + card + '.png'" alt="card1">
      </div>
    </div>
  </div>
</template>

<style scoped>
.bg-orange-dark {
  background-color: #ffc58c;
}
.butte:disabled {
  opacity: 50%
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
      started: false,
      cards: null,
      attackMode: false,
      first: true
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
        case "player-quit":
          let i = this.players.indexOf(data.action.target)
          this.players.splice(i, 1)
          break
        case "host":
          this.host = data.action.target === this.name
          break
        case "start-game":
          this.cards = ["purple_back", "purple_back", "purple_back", "purple_back"]
          break
        case "player-deal-hand":
          this.cards = data.action.target.split(",")
          this.started = true
          break
        case "attack-state":
          this.attackMode = data.action.target === "true"
          break
      }
    },
    startGame() {
      this.ws.send(JSON.stringify({room_id: this.code, action: {action_type: "start-game", origin: this.name}}))
    },
    continueGame() {
      if (this.first) {
        this.first = false
      }
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
