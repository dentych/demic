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
                v-else
                @click="continueGame"
                :disabled="attackMode || gameEnd">
          <span v-if="rounds < 15">Turn next card</span>
          <span v-else>Everyone call out their cards</span>
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
                @click="attack"
                :disabled="!haveCardButtonEnabled">
          I have the card!
        </button>
      </div>
    </div>

    <!-- Player cards -->
    <div class="flex" v-if="cards">
      <div class="flex-auto" v-for="(card, index) in cards">
        <img :src="'/playing-cards/' + card.card + '.png'" alt="card1" v-if="card.show">
        <img :src="'/playing-cards/purple_back.png'" alt="card1" v-else @click="turnCard(index)">
      </div>
    </div>

    <!-- Overlay for choosing a player -->
    <div class="overlay bg-white" v-if="showDialogPlayerPicker">
      <div class="flex flex-col space-y-8 p-4">
        <button class="focus:outline-none bg-yellow-700 text-yellow-100 text-lg py-6 butte"
                v-for="p in players"
                v-if="p !== name && !attackedPlayers.includes(p)"
                @click="dialogChoosePlayer(p)">
          {{ p }}
        </button>
        <button class="focus:outline-none bg-yellow-600 text-yellow-100 text-lg py-6 butte"
                @click="dialogChoosePlayer(null)">
          BACK
        </button>
      </div>
    </div>

    <!-- ATTACK OVERLAY -->
    <div class="overlay bg-white" v-for="(attacker, index) in attacks">
      <div class="flex flex-col text-center p-4">
        <p class="text-xl mb-8">You have been attacked by {{ attacker }}</p>
        <button class="focus:outline-none bg-yellow-700 text-yellow-100 text-lg py-6 butte mb-8"
                @click="drink(index)">
          Drink
        </button>
        <button class="focus:outline-none bg-yellow-700 text-yellow-100 text-lg py-6 butte"
                @click="demandShowCard(index)">
          Demand show card
        </button>
      </div>
    </div>

    <!-- REJECTION OVERLAY -->
    <div class="overlay bg-white" v-for="(rejection, rejectionIndex) in rejections">
      <div class="flex flex-col text-center p-4">
        <p class="flex-grow text-xl mb-8">{{ rejection }} doesn't believe you and wants you to show that you have the
          card!<br><br>Please pick the card below</p>
        <div class="flex">
          <div class="flex-auto" v-for="(card, cardIndex) in cards">
            <img :src="'/playing-cards/purple_back.png'" alt="card1" @click="showCard(rejectionIndex, cardIndex)">
          </div>
        </div>
      </div>
    </div>

    <div v-if="popupText" class="absolute left-0 bottom-0 right-0">
      <div class="flex flex-col justify-end">
        <div class="bg-green-300 text-center mb-20 mx-10 py-4 text-xl">
          {{ popupText }}
        </div>
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

.overlay {
  @apply absolute top-0 left-0 h-screen w-screen
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
      cards: [],
      attackMode: false,
      showDialogPlayerPicker: false,
      showDialogAttacked: false,
      choosePlayerFunc: null,
      popupText: null,
      attacks: [],
      rejections: [],
      attackedPlayers: [],
      haveCardButtonEnabled: false,
      cardsVisible: true,
      gameEnd: false,
      rounds: 0
    }
  },
  methods: {
    messageHandler(evt) {
      let data = JSON.parse(evt.data)
      console.log(data.payload)
      switch (data.payload.action_type) {
        case "player-join":
          this.players.push(data.payload.target)
          break
        case "player-quit":
          let i = this.players.indexOf(data.payload.target)
          this.players.splice(i, 1)
          break
        case "host":
          this.host = data.payload.target === this.name
          break
        case "start-game":
          this.cards.forEach(x => x.show = false)
          break
        case "player-deal-hand":
          if (this.cards.length === 0) {
            for (let i = 0; i < 4; i++) {
              this.cards.push({card: "purple_back", show: true})
            }
          }
          data.payload.target.split(",").forEach((card, index) => {
            this.cards[index].card = card
          })
          this.started = true
          break
        case "attack-state":
          this.attackMode = data.payload.target === "true"
          break
        case "player-attack":
          this.attacks.push(data.payload.origin)
          break
        case "new-round":
          this.newRound()
          this.haveCardButtonEnabled = true
          this.cards.forEach(card => card.show = false)
          this.rounds++
          break
        case "player-reject-attack":
          this.rejections.push(data.payload.origin)
          break
        case "game-end":
          this.gameEnd = true
          this.haveCardButtonEnabled = false
          break
      }
    },
    newRound() {
      this.attackedPlayers = []
    },
    startGame() {
      this.ws.send(JSON.stringify({action_type: "start-game", payload: {action_type: "start-game", origin: this.name}}))
    },
    continueGame() {
      this.ws.send(JSON.stringify({action_type: "continue", payload: {action_type: "continue", origin: this.name}}))
    },
    dialogChoosePlayer(name) {
      this.showDialogPlayerPicker = false
      if (this.choosePlayerFunc !== null) {
        this.choosePlayerFunc(name)
      }
    },
    attack() {
      this.choosePlayerFunc = name => {
        if (name) {
          this.ws.send(JSON.stringify({
            action_type: "player-attack",
            payload: {action_type: "player-attack", origin: this.name, target: name}
          }))
          this.attackedPlayers.push(name)
          this.popupText = "You attacked " + name
          setTimeout(() => {
            this.popupText = null
          }, 3500)
        }
      }
      this.showDialogPlayerPicker = true
    },
    drink(index) {
      this.ws.send(JSON.stringify({
        action_type: "player-accept-attack",
        payload: {action_type: "player-accept-attack", origin: this.name, target: this.attacks[index]}
      }))
      this.attacks.splice(index, 1)
    },
    demandShowCard(index) {
      this.showDialogAttacked = false
      this.ws.send(JSON.stringify({
        action_type: "player-reject-attack",
        payload: {action_type: "player-reject-attack", origin: this.name, target: this.attacks[index]}
      }))
      this.attacks.splice(index, 1)
    },
    showCard(rejectionIndex, cardIndex) {
      this.rejections.splice(rejectionIndex, 1)
      this.cards[cardIndex].show = true
      this.ws.send(JSON.stringify({
        action_type: "player-pick-card",
        payload: {action_type: "player-pick-card", origin: this.name, target: cardIndex.toString()}
      }))
    },
    turnCard(index) {
      if (!this.gameEnd) return
      this.cards[index].show = true
      this.ws.send(JSON.stringify({
        action_type: "show-card",
        payload: {action_type: "show-card", origin: this.name, target: index.toString()}
      }))
    }
  },
  mounted() {
    this.ws = new WebSocket("ws://" + location.hostname + ":8080/ws")
    this.ws.onopen = () => {
      this.ws.send(JSON.stringify({action_type: "player-join", payload: {action_type: "player-join", origin: this.name, target: this.code}}))
    }

    this.ws.onmessage = this.messageHandler
  }
}
</script>
