import Vue from 'vue'
import Vuex from 'vuex'

import messages from '../../src/messages/messages_pb'

Vue.use(Vuex)

export default new Vuex.Store({
  state: {
    username: 'Zero Cool',
    color: '#ff0000',
    techDebt: 0,
    socket: {
      isConnected: false,
      message: '',
      reconnectError: false
    }
  },
  getters: {
    techDebt (state, getters) {
      console.log('Getting techDebt: ', state.techDebt)
      return state.techDebt
    }
  },
  mutations: {
    set_user_info (state, {username, color}) {
      console.log('Username set to ', username + ' and ' + color)
      state.username = username
      state.color = color
      document.cookie = 'username=' + username + ';color=' + color + ';samesite'
    },
    SOCKET_ONOPEN (state, event) {
      Vue.prototype.$socket = event.currentTarget
      state.socket.isConnected = true
    },
    SOCKET_ONCLOSE (state, event) {
      state.socket.isConnected = false
    },
    SOCKET_ONERROR (state, event) {
      console.error(state, event)
    },
    // default handler called for all methods
    SOCKET_ONMESSAGE (state, message) {
      state.socket.message = message
      console.log('Got response')
      var fileReader = new FileReader()
      fileReader.onload = function (event) {
        var arrayBuffer = event.target.result
        console.log(arrayBuffer)
        var buffer = messages.TechDebt.deserializeBinary(arrayBuffer)
        var pbObject = buffer.toObject()
        console.log(pbObject)
        state.techDebt = pbObject.techdebt
      }
      fileReader.readAsArrayBuffer(message.data)
    },
    // mutations for reconnect methods
    SOCKET_RECONNECT (state, count) {
      console.info(state, count)
    },
    SOCKET_RECONNECT_ERROR (state) {
      state.socket.reconnectError = true
    }
  },
  actions: {
    sendMessage (context, message) {
      Vue.prototype.$socket.send(message)
    },
    setUserInfo (context, {username, color}) {
      context.commit('set_user_info', {username, color})
    },
    getDebt () {
      Vue.prototype.$socket.send(0)
    },
    addDebt (context, count) {
      context.state.techDebt += count
      var message = new messages.TechDebt()
      message.setUsername(context.state.username)
      message.setColor(context.state.color)
      message.setTechdebt(count)
      console.log('Username: ' + context.state.username)
      console.log('Color: ' + context.state.color)
      let rawMessage = message.serializeBinary()
      console.log(rawMessage)
      Vue.prototype.$socket.send(rawMessage)
    }
  }
})
