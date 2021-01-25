<template>
  <div class="flex justify-center m-4">
    <div class="border rounded flex-auto sm:flex-none sm:w-1/4 bg-blue-100">
      <div class="flex flex-col p-5">
        <h1 class="place-self-center text-3xl">Pyramid</h1>
        <p class="font-medium text-lg mb-2">Room code:</p>
        <input class="mb-2 py-1 px-2 border rounded focus:outline-none uppercase" type="text" v-model="code"
               placeholder="Enter a valid room code">
        <p class="font-medium text-lg mb-2">Name:</p>
        <input class="mb-2 py-1 px-2 border rounded focus:outline-none uppercase" type="text" v-model="name"
               placeholder="Enter a valid room code">
        <button class="mb-2 py-2 text-white bg-blue-600 mx-2" @click="join">Join room</button>
      </div>
      <hr>
      <div class="flex flex-col p-5">
        <p class="font-medium text-lg mb-2">Create a new room:</p>
        <button class="mb-6 py-2 text-white bg-blue-600 mx-2" @click="create">Create new room</button>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      code: "",
      name: ""
    }
  },
  methods: {
    create() {
      console.log("Opening websocket")
      let ws = new WebSocket("ws://localhost:8080/ws")
      ws.addEventListener("open", function (event) {
        console.log("WebSocket opened")
        ws.send(JSON.stringify({action: "start-game"}))
        console.log(event)
      })

      ws.addEventListener("message", function (event) {
        if (event.data === "hello") {
          console.log("HELLO received")
        }
        console.log("MESSAGE:", event)
      })

      this.$store.commit("set", ws)
    },
    join() {
    }
  }
}
</script>