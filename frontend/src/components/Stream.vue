<script setup lang="ts">
import { onBeforeUnmount, ref } from "vue";
import { ConnectError, Code } from "@bufbuild/connect";
import { MetricsService } from "../../gen/api/v1/metrics_connect";
import { getClient } from "../lib/api";

const client = getClient(MetricsService);

const abortCtl = ref(new AbortController());

const percent = ref<number | null>(null);

onBeforeUnmount(() => {
  stopStream();
});

async function startStream() {
  stopStream();

  try {
    for await (const res of client.cpuUsageStream(
      {},
      { signal: abortCtl.value.signal },
    )) {
      percent.value = res.percent;
    }
  } catch (err: any) {
    if (err instanceof ConnectError && err.code != Code.Canceled) {
      console.error(err);
    }
  }

  percent.value = null;
}

function stopStream() {
  if (!abortCtl.value.signal.aborted) {
    abortCtl.value.abort();
  }
  abortCtl.value = new AbortController();

  percent.value = null;
}
</script>

<template>
  <div class="card">
    <p>ServerSideStreaming Test</p>
    <button type="button" @click="startStream">Start</button>
    <button type="button" @click="stopStream">Stop</button>
    <p>CPU: {{ percent?.toFixed(2) ?? "-" }} %</p>
  </div>
</template>

<style scoped></style>
