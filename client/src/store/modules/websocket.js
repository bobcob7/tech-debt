// initial state
import Vue from 'vue'
import VueNativeSock from 'vue-native-websocket'
import messages from '../../../src/messages/messages_pb'
Vue.use(VueNativeSock, 'ws://localhost:3000/ws', {})

const state = {
  username: 'test',
  techDebt: 0,
  socket: {
    isConnected: false,
    message: '',
    reconnectError: false
  }
}

// getters
const getters = {
  techDebt (state, getters) {
    return state.techDebt
  }
}

// actions
const actions = {
  setUsername (context, username) {
    context.commit('set_username', username)
  },
  getDebt () {
    Vue.prototype.$socket.send(0)
  },
  addDebt (context, count) {
    context.state.techDebt += count
    var message = new messages.TechDebt()
    message.setUsername(context.state.username)
    message.setTechdebt(count)
    let rawMessage = message.serializeBinary()
    console.log(rawMessage)
    Vue.prototype.$socket.send(rawMessage)
  }
}

// mutations
const mutations = {
  set_username (state, username) {
    console.log('Username set to ', username)
    state.username = username
    document.cookie = 'username=' + username + ';samesite'
  },
  SOCKET_ONOPEN (state, event) {
    console.log('WS Opened')
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
    var buffer = new messages.TechDebt().deserializeBinary(message)
    console.log('Got message', buffer)
    state.techDebt = buffer.techDebt
  },
  // mutations for reconnect methods
  SOCKET_RECONNECT (state, count) {
    console.info(state, count)
  },
  SOCKET_RECONNECT_ERROR (state) {
    state.socket.reconnectError = true
  }
}

export default {
  namespaced: true,
  state,
  getters,
  actions,
  mutations
}
