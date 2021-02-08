<template>
  <div class="flex bg-orange-light h-screen p-4">
    <div class="flex flex-col flex-grow border-r-4 border-yellow-600 border-dashed" v-if="ws">
      <h1 class="text-3xl border-b-2 border-dashed border-yellow-600 pb-4 mr-4 text-yellow-800 truncate">
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
    </div>
    <div class="grid grid-rows-5 gap-2">

      <div class="flex justify-center space-x-16">
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
        case"player-quit":
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