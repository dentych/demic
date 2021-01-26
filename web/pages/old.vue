<template>
  <div>
    <div class="flex flex-col h-screen">
      <div class="flex flex-col flex-grow items-center">
        <div class="flex" v-for="row in rows.length">
          <div v-for="n in row">
            <img :src="'/playing-cards/' + cards[rows[row - 1] + n - 1] + '.png'"
                 class="w-12 m-3 transform"
                 :class="{'rotate-90': !isEven(rows[row - 1] + n - 1)}"
                 @click="highlightCard(cards[rows[row - 1] + n - 1])">
          </div>
        </div>
      </div>
      <div>
        <div class="flex flex-col bg-blue-100">
          <p class="self-center font-semibold m-2">Actions</p>
          <div class="flex justify-around">
            <button class="bg-blue-600 p-2 text-white">Give sip</button>
            <button class="bg-blue-600 p-2 text-white" @click="turnCard">Turn card</button>
          </div>
        </div>
      </div>
      <div>
        <div class="flex flex-col bg-blue-100">
          <p class="self-center font-semibold m-2">Your cards</p>
          <div class="flex justify-around">
            <img src="/playing-cards/blue_back.png" alt="BLUE BACK" class="w-12 m-2">
            <img src="/playing-cards/blue_back.png" alt="BLUE BACK" class="w-12 m-2">
            <img src="/playing-cards/blue_back.png" alt="BLUE BACK" class="w-12 m-2">
            <img src="/playing-cards/blue_back.png" alt="BLUE BACK" class="w-12 m-2">
          </div>
        </div>
      </div>
    </div>
    <div v-if="highlightedCard" class="fixed top-0 left-0 w-screen h-screen"
         @click="highlightCard(null)">
      <div class="absolute top-0 left-0 w-full h-full bg-black opacity-50">HEJ</div>
      <div class="absolute top-0 left-0 w-full h-full flex justify-center items-center">
        <img :src="'/playing-cards/' + highlightedCard + '.png'" class="h-64">
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: "pyramid",
  data() {
    return {
      rows: [0, 1, 3, 6, 10],
      cards: ["10C", "KH", "blue_back", "blue_back", "blue_back", "blue_back", "blue_back", "blue_back", "blue_back", "blue_back", "blue_back", "blue_back", "blue_back", "blue_back", "blue_back"],
      highlightedCard: null
    }
  },
  methods: {
    isEven(number) {
      return number % 2 === 0
    },
    highlightCard(card) {
      this.highlightedCard = card
    },
    turnCard() {
      for (let i = 0; i < this.cards.length; i++) {
        if (this.cards[i] === "blue_back") {
          this.cards.splice(i, 1, "8C")
          break
        }
      }
      console.log(this.cards)
    }
  }
}
</script>