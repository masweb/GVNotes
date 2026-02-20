<script lang="ts" setup>
import { Splitpanes, Pane } from 'splitpanes'
</script>

<template>
 <splitpanes style="height: 100vh" class="split-dark">
  <pane>
   <splitpanes class="split-dark panels-limited">
    <pane>1</pane>
    <pane>2</pane>
    <pane>3</pane>
   </splitpanes>
  </pane>
  <pane max-size="80" min-size="30">contenido principal</pane>
 </splitpanes>
</template>

<style>
.splitpanes__pane {
 /*interior de cada panel*/
}

.panels-limited > .splitpanes__pane {
 min-width: 120px;
 max-width: 420px;
}
</style>
