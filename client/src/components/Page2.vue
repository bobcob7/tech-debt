<template>
  <div class="hello">
    <h1>Alter Websockets Page</h1>
    <router-link to="/">Go Back</router-link>
    <div>
      <input v-model="nextMessage">
      <button @click="sendMessage(nextMessage)">Send</button>
    </div>
    <div>
      <button @click="logs=''">Clear</button>
    </div>
    <h3>Messages</h3>
    <p v-html="logs"></p>
  </div>
</template>

<script>
export default {
  name: 'HelloWorld',
  data () {
    return {
      logs: '',
      nextMessage: ''
    }
  },
  methods: {
    onOpen: function (event) {
      console.log('Opened WS')
    },
    onClose: function (event) {
      console.log('Closed WS')
    },
    onMessage: function (event) {
      console.log('Message WS')
      console.log(event.data)
      this.logs += 'Received: ' + event.data + '<br>'
    },
    onError: function (event) {
      console.log('Error WS')
    },
    sendMessage: function (message) {
      this.logs += 'Sending: ' + message + '<br>'
      // this.$socket.send(message)
      this.$store.dispatch('websocket/send', message)
    }
  },
  created: function () {
    // this.$connect()
    this.$options.sockets.onmessage = this.onMessage
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
</style>
