<template>
  <div class="flex bg-blue-200 h-screen p-4">
    <div class="flex flex-col w-6/12 border-r-4 border-blue-600 border-dashed" v-if="ws">
      <h1 class="text-2xl border-b-2 border-dashed border-blue-400 pb-4 mr-4">Room code: {{ code }}</h1>
      <div class="py-4 mr-4 border-b-2 border-dashed border-blue-400">
        <span class="font-bold">Players:</span><br>
        <span v-if="!players.length">( No players )</span>
        <ul>
          <li v-for="player in players">{{ player }}</li>
        </ul>
      </div>
    </div>
    <div class="flex-grow grid grid-rows-5 gap-2">

      <div class="flex justify-center space-x-4">
        <div class="h-full">
          <img src="/playing-cards/blue_back.png" class="h-full" v-if="!cards[14]">
          <img :src="'/playing-cards/' + cards[14] + '.png'" class="h-full" v-if="cards[14]">
        </div>
      </div>

      <div class="flex justify-center space-x-4">
        <div class="h-full">
          <img src="/playing-cards/blue_back.png" class="h-full" v-if="!cards[12]">
          <img :src="'/playing-cards/' + cards[12] + '.png'" class="h-full" v-if="cards[12]">
        </div>
        <div class="h-full">
          <img src="/playing-cards/blue_back.png" class="h-full" v-if="!cards[13]">
          <img :src="'/playing-cards/' + cards[13] + '.png'" class="h-full" v-if="cards[13]">
        </div>
      </div>

      <div class="flex justify-center space-x-4">
        <div class="h-full">
          <img src="/playing-cards/blue_back.png" class="h-full" v-if="!cards[9]">
          <img :src="'/playing-cards/' + cards[9] + '.png'" class="h-full" v-if="cards[9]">
        </div>
        <div class="h-full">
          <img src="/playing-cards/blue_back.png" class="h-full" v-if="!cards[10]">
          <img :src="'/playing-cards/' + cards[10] + '.png'" class="h-full" v-if="cards[10]">
        </div>
        <div class="h-full">
          <img src="/playing-cards/blue_back.png" class="h-full" v-if="!cards[11]">
          <img :src="'/playing-cards/' + cards[11] + '.png'" class="h-full" v-if="cards[11]">
        </div>
      </div>

      <div class="flex justify-center space-x-4">
        <div class="h-full">
          <img src="/playing-cards/blue_back.png" class="h-full" v-if="!cards[5]">
          <img :src="'/playing-cards/' + cards[5] + '.png'" class="h-full" v-if="cards[5]">
        </div>
        <div class="h-full">
          <img src="/playing-cards/blue_back.png" class="h-full" v-if="!cards[6]">
          <img :src="'/playing-cards/' + cards[6] + '.png'" class="h-full" v-if="cards[6]">
        </div>
        <div class="h-full">
          <img src="/playing-cards/blue_back.png" class="h-full" v-if="!cards[7]">
          <img :src="'/playing-cards/' + cards[7] + '.png'" class="h-full" v-if="cards[7]">
        </div>
        <div class="h-full">
          <img src="/playing-cards/blue_back.png" class="h-full" v-if="!cards[8]">
          <img :src="'/playing-cards/' + cards[8] + '.png'" class="h-full" v-if="cards[8]">
        </div>
      </div>

      <div class="flex justify-center space-x-4">
        <div class="h-full">
          <img src="/playing-cards/blue_back.png" class="h-full" v-if="!cards[0]">
          <img :src="'/playing-cards/' + cards[0] + '.png'" class="h-full" v-if="cards[0]">
        </div>
        <div class="h-full">
          <img src="/playing-cards/blue_back.png" class="h-full" v-if="!cards[1]">
          <img :src="'/playing-cards/' + cards[1] + '.png'" class="h-full" v-if="cards[1]">
        </div>
        <div class="h-full">
          <img src="/playing-cards/blue_back.png" class="h-full" v-if="!cards[2]">
          <img :src="'/playing-cards/' + cards[2] + '.png'" class="h-full" v-if="cards[2]">
        </div>
        <div class="h-full">
          <img src="/playing-cards/blue_back.png" class="h-full" v-if="!cards[3]">
          <img :src="'/playing-cards/' + cards[3] + '.png'" class="h-full" v-if="cards[3]">
        </div>
        <div class="h-full">
          <img src="/playing-cards/blue_back.png" class="h-full" v-if="!cards[4]">
          <img :src="'/playing-cards/' + cards[4] + '.png'" class="h-full" v-if="cards[4]">
        </div>
      </div>

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
      players: [],
      cards: []
    }
  },
  mounted() {
    this.ws = new WebSocket("ws://" + location.hostname + ":8080/ws")
    this.ws.onopen = () => {
      this.ws.send(JSON.stringify({action: {action_type: "create-game"}}))
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
      switch (data.action.action_type) {
        case"create-game":
          console.log("pyramid created with ID", data.action.target)
          this.code = data.action.target
          break
        case"player-join":
          if (data.action.target !== "HOST") {
            this.players.push(data.action.target)
            console.log(this.players)
          }
          break
        case"player-left":
          let i = this.players.indexOf(data.action.target)
          this.players.splice(i, 1)
          break
        case"gm-start-game":
          break
        case "player-deal-hand":
          this.cards.push(data.action.target)
          break
      }
    },
  }
}
</script>

<style scoped>

</style>