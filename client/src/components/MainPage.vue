<template>
  <b-container>
    <b-row class="justify-content-md-center">
      <h1>Tech Debt Clock</h1>
    </b-row>
    <b-row class="justify-content-md-center">
      <h2>{{$store.getters['techDebt']}}</h2>
    </b-row>
    <b-row>
      <b-col><b-button variant="danger" size="lg" @click="addDebt(1)">+</b-button></b-col>
      <b-col><b-button variant="danger" size="lg" @click="addDebt(10)">+</b-button></b-col>
      <b-col><b-button variant="danger" size="lg" @click="addDebt(100)">+</b-button></b-col>
      <b-col><b-button variant="danger" size="lg" @click="addDebt(1000)">+</b-button></b-col>
    </b-row>
    <b-row>
      <b-col><h5>Little</h5></b-col>
      <b-col><h5>Bit</h5></b-col>
      <b-col><h5>Lot</h5></b-col>
      <b-col><h5>Fuck-Ton</h5></b-col>
    </b-row>
    <b-row>
      <b-col><b-button variant="success" size="lg" @click="addDebt(-1)">-</b-button></b-col>
      <b-col><b-button variant="success" size="lg" @click="addDebt(-10)">-</b-button></b-col>
      <b-col><b-button variant="success" size="lg" @click="addDebt(-100)">-</b-button></b-col>
      <b-col><b-button variant="success" size="lg" @click="addDebt(-1000)">-</b-button></b-col>
    </b-row>
    <b-row>
    <b-col>
      <b-btn v-b-modal.changeUsernameModal>Launch demo modal</b-btn>
    </b-col>
    </b-row>

    <!-- Modal Component -->
    <b-modal id="changeUsernameModal"
    centered
    title="Set Alias"
    @ok="setUsername"
    @hidden="resetUsername">
      <b-container fluid>
        <b-row>
          <b-col>
      <b-form-group id="usernameInputGroup"
                    label="Alias:"
                    label-for="usernameInput"
                    horizontal>
        <b-form-input id="usernameInput"
                      type="text"
                      v-model="tempUsername"
                      required
                      placeholder="Enter email">
        </b-form-input>
      </b-form-group>
          </b-col>
        </b-row>
        <b-row>
          <b-col>
      <b-form-group label="Color"
      horizontal>
        <b-form-input type="color"
        placeholder="#ff0000"
        v-model="tempColor">
        </b-form-input>
      </b-form-group>
          </b-col>
        </b-row>
      </b-container>
    </b-modal>
  </b-container>
</template>

<script>
export default {
  name: 'HelloWorld',
  data () {
    return {
      logs: '',
      tempUsername: 'Zero Cool',
      tempColor: '#ff0000',
      username: 'Zero Cool',
      color: '#ff0000'
    }
  },
  methods: {
    addDebt: function (count) {
      console.log('Adding Debt: ', count)
      this.$store.dispatch('addDebt', count)
    },
    setUsername: function () {
      console.log('Setting Username')
      this.username = this.tempUsername
      this.color = this.tempColor
      this.$store.dispatch('setUsername', this.username)
    },
    resetUsername: function () {
      this.tempUsername = this.username
      this.tempColor = this.color
    }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
</style>
