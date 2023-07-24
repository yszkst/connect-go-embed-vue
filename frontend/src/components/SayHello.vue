<script setup lang="ts">
import { ref } from "vue";

import { SayHelloService } from "../../gen/api/v1/sayhello_connect.ts";

import { getClient } from "../lib/api";

const client = getClient(SayHelloService);

const name = ref("");
const rep = ref("");

async function sayhello() {
  const res = await client.sayHello({ name: name.value });
  rep.value = res.reply;
}
</script>

<template>
  <div class="card">
    <p>SayHello Test</p>
    <label>Your name <input v-model="name" /></label>
    <button type="button" @click="sayhello()">SayHello</button>
    <p>{{ rep || "..." }}</p>
  </div>
</template>

<style scoped></style>
