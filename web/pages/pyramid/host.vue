<template>
  <div class="flex bg-orange-light h-screen w-screen p-4">

    <div class="flex-none flex flex-col border-r-4 border-yellow-600 border-dashed w-2/12" v-if="ws">
      <h1 class="text-3xl border-b-2 border-dashed border-yellow-600 pb-4 mr-4 text-yellow-800">
        <span class="font-extrabold">ROOM CODE:</span> {{ code }}
      </h1>
      <div class="py-4 mr-4 border-b-2 border-dashed border-yellow-600 text-2xl">
        <span class="text-yellow-800 font-bold">PLAYERS:</span>
        <br>
        <span class="ml-2" v-if="!players.length">( No players )</span>
        <ul class="ml-2 truncate">
          <li class="my-1" v-for="player in players">{{ player }}</li>
        </ul>
      </div>
      <div class="py-4 mr-4 border-b-2 border-dashed border-yellow-600 text-xl">
        <span class="font-bold">Wall of Rules:</span>
        <br>
        The board contains 15 cards. Each round, 1 card on the board is turned. If you have a card with the same rank,
        you may give out sips to 1 player for each card. It is okay to lie about having the card. If the player you
        give sips to doesn't believe you, you have to show the correct card. If you fail to do so, YOU drink instead.
        If you show the correct card, the other player drinks double. In the end, you call out loud the rank of each
        card and then press it to show the others. For each card you get wrong, you drink.
      </div>
    </div>

    <div class="grid grid-rows-5 flex-flow-col gap-2">
      <div class="flex justify-center">
        <div class="h-full">
          <img :src="'/playing-cards/' + cards[14] + '.png'" class="h-full" v-if="cards[14]">
          <img v-else src="/playing-cards/purple_back.png" class="h-full">
        </div>
      </div>

      <div class="flex justify-center space-x-12">
        <div class="h-full">
          <img :src="'/playing-cards/' + cards[12] + '.png'" class="h-full" v-if="cards[12]">
          <img v-else src="/playing-cards/purple_back.png" class="h-full">
        </div>
        <div class="h-full transform rotate-90">
          <img :src="'/playing-cards/' + cards[13] + '.png'" class="h-full" v-if="cards[13]">
          <img v-else src="/playing-cards/purple_back.png" class="h-full">
        </div>
      </div>

      <div class="flex justify-center space-x-12">
        <div class="h-full transform rotate-90">
          <img :src="'/playing-cards/' + cards[9] + '.png'" class="h-full" v-if="cards[9]">
          <img v-else src="/playing-cards/purple_back.png" class="h-full">
        </div>
        <div class="h-full">
          <img :src="'/playing-cards/' + cards[10] + '.png'" class="h-full" v-if="cards[10]">
          <img v-else src="/playing-cards/purple_back.png" class="h-full">
        </div>
        <div class="h-full transform rotate-90">
          <img :src="'/playing-cards/' + cards[11] + '.png'" class="h-full" v-if="cards[11]">
          <img v-else src="/playing-cards/purple_back.png" class="h-full">
        </div>
      </div>

      <div class="flex justify-center space-x-16">
        <div class="h-full transform rotate-90">
          <img :src="'/playing-cards/' + cards[5] + '.png'" class="h-full" v-if="cards[5]">
          <img v-else src="/playing-cards/purple_back.png" class="h-full">
        </div>
        <div class="h-full">
          <img :src="'/playing-cards/' + cards[6] + '.png'" class="h-full" v-if="cards[6]">
          <img v-else src="/playing-cards/purple_back.png" class="h-full">
        </div>
        <div class="h-full transform rotate-90">
          <img :src="'/playing-cards/' + cards[7] + '.png'" class="h-full" v-if="cards[7]">
          <img v-else src="/playing-cards/purple_back.png" class="h-full">
        </div>
        <div class="h-full">
          <img :src="'/playing-cards/' + cards[8] + '.png'" class="h-full" v-if="cards[8]">
          <img v-else src="/playing-cards/purple_back.png" class="h-full">
        </div>
      </div>

      <div class="flex justify-center space-x-20">
        <div class="h-full">
          <img :src="'/playing-cards/' + cards[0] + '.png'" class="h-full" v-if="cards[0]">
          <img v-else src="/playing-cards/purple_back.png" class="h-full">
        </div>
        <div class="h-full transform rotate-90">
          <img :src="'/playing-cards/' + cards[1] + '.png'" class="h-full" v-if="cards[1]">
          <img v-else src="/playing-cards/purple_back.png" class="h-full">
        </div>
        <div class="h-full">
          <img :src="'/playing-cards/' + cards[2] + '.png'" class="h-full" v-if="cards[2]">
          <img v-else src="/playing-cards/purple_back.png" class="h-full">
        </div>
        <div class="h-full transform rotate-90">
          <img :src="'/playing-cards/' + cards[3] + '.png'" class="h-full" v-if="cards[3]">
          <img v-else src="/playing-cards/purple_back.png" class="h-full">
        </div>
        <div class="h-full">
          <img :src="'/playing-cards/' + cards[4] + '.png'" class="h-full" v-if="cards[4]">
          <img v-else src="/playing-cards/purple_back.png" class="h-full">
        </div>
      </div>
    </div>

    <div class="flex-none flex flex-col border-l-4 border-yellow-600 border-dashed w-2/12">
      <h1 class="text-3xl border-b-2 border-dashed border-yellow-600 pb-4 mx-4 text-yellow-800 truncate text-center">
        ACTIONS
      </h1>
      <div class="flex flex-col overflow-hidden">
        <p class="p-4" v-for="text in actionTexts">{{ text }}</p>
      </div>
    </div>
  </div>
</template>

<style>
.bg-orange-light {
  background-color: #ffddba;
}
</style>

<script>
export default {
  data() {
    return {
      code: "",
      name: "",
      ws: null,
      players: [],
      cards: [],
      actionTexts: []
    }
  },
  mounted() {
    this.ws = new WebSocket("ws://" + location.hostname + ":8080/ws")
    this.ws.onopen = () => {
      this.ws.send(JSON.stringify({action_type: "create-game", payload: {action_type: "create-game"}}))
    }

    this.ws.onmessage = this.messageHandler

    this.ws.onerror = err => {
      console.log(err)
    }
  },
  methods: {
    messageHandler(evt) {
      console.log(evt)
      let data = JSON.parse(evt.data)
      switch (data.payload.action_type) {
        case"create-game":
          console.log("pyramid created with ID", data.payload.target)
          this.code = data.payload.target
          break
        case"player-join":
          if (data.payload.target !== "HOST") {
            this.actionTexts.unshift("Player '" + data.payload.target + "' joined the game!")
            this.players.push(data.payload.target)
            console.log(this.players)
          }
          break
        case"player-quit":
          this.actionTexts.unshift("Player '" + data.payload.target + "' quit the game!")
          let i = this.players.indexOf(data.payload.target)
          this.players.splice(i, 1)
          break
        case "player-deal-hand":
          this.actionTexts.unshift("New card turned. Players can attack!")
          this.cards.push(data.payload.target)
          break
        case "player-pick-card":
          let currentCardValue = this.cards[this.cards.length - 1].slice(0, -1)
          let target = data.payload.target.slice(0, -1)
          if (currentCardValue === target) {
            this.actionTexts.unshift(data.payload.origin + " picked the CORRECT card!")
          } else {
            this.actionTexts.unshift(data.payload.origin + " picked the WRONG card, and must drink!")
          }
          break
        case "player-attack":
          this.actionTexts.unshift(data.payload.origin + " attacks " + data.payload.target + "!")
          break
        case "player-accept-attack":
          this.actionTexts.unshift(data.payload.origin + " chose to drink!")
          break
        case "player-reject-attack":
          this.actionTexts.unshift(data.payload.origin + " rejected the attack from " + data.payload.target + ", who must now show the card!")
          break
        case "show-card":
          this.actionTexts.unshift(data.payload.origin + " shows card " + data.payload.target.slice(0, -1))
          break
        case "game-end":
          this.actionTexts.unshift("No more cards to turn. All players will take turns calling out loud their cards and press them!")
          break
      }
    },
  }
}
</script>

<style scoped>

</style>