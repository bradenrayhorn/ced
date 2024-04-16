<script lang="ts">
  import type { PageData } from "./$types";

  export let data: PageData;
</script>

<svelte:head>
  <title>About ced</title>
</svelte:head>

<h2 class="text-3xl font-bold mb-8">About ced</h2>

<div class="flex flex-col gap-4">
  <p>ced is a self-hosted RSVP service.</p>

  <p>
    Source code is available on{" "}<a
      class="anchor"
      href="https://github.com/bradenrayhorn/ced">GitHub</a
    >. Licensed under
    <a
      class="anchor"
      href="https://github.com/bradenrayhorn/ced/blob/main/LICENSE">AGPLv3</a
    >.
  </p>

  <p>
    Icons courtesy of Material Design Icons, licensed under <a
      class="anchor"
      href="https://www.apache.org/licenses/LICENSE-2.0.txt">Apache-2.0</a
    >. Source available on
    <a class="anchor" href="https://github.com/google/material-design-icons"
      >GitHub</a
    >.
  </p>

  <p>
    Servers running using <a class="anchor" href="https://go.dev">Go</a> and
    <a class="anchor" href="https://nodejs.org/">Node.js</a>.
  </p>
</div>

<h2 class="text-3xl font-bold my-8">Bundled Licenses</h2>

<div class="overflow-hidden">
  <h3 class="text-2xl font-bold my-8">Client</h3>

  <div class="w-full whitespace-pre-wrap">
    {data.notices.ui}
  </div>

  <h3 class="text-2xl font-bold my-8">Server</h3>

  <div class="w-full whitespace-pre-wrap">
    {data.notices.uiServer}
  </div>

  <div class="w-full whitespace-pre-wrap">
    {data.notices.server}
  </div>
</div>
