<template>
  <p> Creating pyramid game </p>
</template>

<script>
export default {
  mounted() {
    console.log("Creating new game")
    let ws = new WebSocket("ws://" + location.hostname + ":8080/ws")

    ws.onopen = () => {
      ws.send(JSON.stringify({action: {action_type: "create-game"}}))
      ws.onopen = null
    }

    ws.onmessage = evt => {
      let data = JSON.parse(evt.data)
      ws.onmessage = null
      this.$store.commit("setWebsocket", ws)
      if (data.action.action_type === "game-created") {
        console.log("going to pyramid lala")
        this.$router.replace("/pyramid/" + data.action.target)
      }
    }
  }
}
</script>