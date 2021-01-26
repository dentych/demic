<template>
  <p> Creating pyramid game </p>
</template>

<script>
export default {
  mounted() {
    let ws = new WebSocket("ws://localhost:8080/ws")

    ws.onopen = () => {
      ws.send(JSON.stringify({action: {action_type: "create-game"}}))
      ws.onopen = null
    }

    ws.onmessage = evt => {
      let data = JSON.parse(evt.data)
      ws.onmessage = null
      this.$store.commit("set", ws)
      if (data.action.action_type === "game-created") {
        console.log("going to pyramid lala")
        this.$router.replace("/pyramid/" + data.action.target)
      }
    }
  }
}
</script>