export const state = () => ({
    websocket: null,
    roomCode: null,
    playerName: null
})

export const mutations = {
    setWebsocket(state, websocket) {
        state.websocket = websocket;
    },
    setRoomCode(state, roomCode) {
        state.roomCode = roomCode
    },
    setPlayerName(state, playerName) {
        state.playerName = playerName
    }
}