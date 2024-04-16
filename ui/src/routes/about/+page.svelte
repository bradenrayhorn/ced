<script lang="ts">
  import type { PageData } from "./$types";

  export let data: PageData;
</script>

<svelte:head>
  <title>About ced</title>
</svelte:head>

<h2 class="text-3xl font-bold mb-8">About ced</h2>

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

<h2 class="text-3xl font-bold my-8">Bundled Licenses</h2>

<div class="overflow-hidden">
  <div class="w-full whitespace-pre-wrap">
    {data.notices.ui}
  </div>

  <div class="w-full whitespace-pre-wrap">
    {data.notices.server}
  </div>
</div>
