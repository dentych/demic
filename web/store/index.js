export const state = () => ({
    websocket: null
})

export const mutations = {
    set(state, websocket) {
        state.websocket = websocket;
    }
}