<template>
  <div>
    <div class="flex justify-center" v-if="roomCode">Room code: {{ roomCode }}</div>
    <div>Players:</div>
    <ul>
      <li v-for="player in players">{{ player }}</li>
    </ul>
  </div>
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
    this.roomCode = this.$route.params.id
    if (this.$store.state.websocket === null) {
      if (this.$store.state.playerName !== null) {
        let ws = new WebSocket("ws://" + location.hostname + ":8080/ws")
        ws.onopen = evt => {
          ws.send(JSON.stringify({
            room_id: this.roomCode,
            action: {
              action_type: "player-join",
              origin: this.$store.state.playerName
            }
          }))
        }
        this.$store.commit("setWebsocket", ws)
      } else {
        this.$router.push("/")
        return
      }
    } else {
      console.log("Websocket is")
      console.log(this.$store.state.websocket)
    }
    this.$store.state.websocket.onmessage = evt => {
      console.log(evt)
      let data = JSON.parse(evt.data)
      if (data.action.action_type === "player-joined") {
        this.players.push(data.action.target)
      }
    }
  }
}
</script>