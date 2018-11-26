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
  </b-container>
</template>

<script>
export default {
  name: 'HelloWorld',
  data () {
    return {
      logs: ''
    }
  },
  methods: {
    addDebt: function (count) {
      console.log('Adding Debt: ', count)
      this.$store.dispatch('addDebt', count)
    }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
</style>
