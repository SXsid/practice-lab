<script setup>
    const name ="Harsh"
import HelloWorld from './components/HelloWorld.vue'
const AlertValue=(count)=>{


    alert("current value is "+count)

}
</script>

<template>
  <HelloWorld :name="name" v-on:countClicked="AlertValue" />
</template>
