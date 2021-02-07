<template>
  <div class="flex justify-center m-4">
    <div class="border-2 border-yellow-600 rounded flex-auto sm:flex-none sm:w-1/4 bg-orange-light">
      <div class="flex flex-col p-5">
        <h1 class="place-self-center text-5xl mb-4 text-yellow-700">PYRAMID</h1>
        <p class="font-medium text-lg mb-2">Room code:</p>
        <input class="mb-2 py-1 px-2 text-xl border border-yellow-600 rounded focus:outline-none uppercase"
               type="text" v-model="code"
               placeholder="Enter a valid room code"
               @input="update">
        <p class="font-medium text-lg mb-2">Name:</p>
        <input class="mb-2 py-1 px-2 text-xl border border-yellow-600 rounded focus:outline-none uppercase"
               type="text" v-model="name"
               placeholder="Enter a name"
               @input="update">
        <button class="mb-2 py-2 text-white bg-yellow-600 mx-2" @click="$emit('join', {code: code, name: name})">Join room</button>
      </div>
      <hr>
      <div class="flex flex-col p-5">
        <p class="font-medium text-lg mb-2">Create a new room:</p>
        <button class="mb-6 py-2 text-white bg-yellow-600 mx-2" @click="$emit('create')">Create new room</button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.bg-orange-light {
  background-color: #ffddba;
}
</style>

<script>
export default {
  data() {
    return {
      name: "",
      code: ""
    }
  },
  methods: {
    update() {
      this.name = this.name?.toUpperCase()
      this.code = this.code?.toUpperCase()
      localStorage.setItem("name", this.name)
      localStorage.setItem("code", this.code)
    }
  },
  mounted() {
    this.name = localStorage.getItem("name")
    this.code = localStorage.getItem("code")
  }
}
</script>

<style scoped>

</style>