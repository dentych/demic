<template>
  <div class="flex justify-center" v-if="roomCode">Room code: {{ roomCode }}</div>
</template>

<script>
export default {
  data() {
    return {
      roomCode: null,
      players: []
    }
  },
  mounted() {
    if (this.$store.state.websocket === null) {
      this.$router.push("/")
    } else {
      this.roomCode = this.$route.params.id
      console.log("Websocket is")
      console.log(this.$store.state.websocket)
      this.$store.state.websocket.onmessage = evt => {
        console.log(evt)
        let data = JSON.parse(evt.data)
        if (data.action.action_type === "player-join") {
          this.players.push(data.action.target)
        }
      }
    }
  }
}
</script>