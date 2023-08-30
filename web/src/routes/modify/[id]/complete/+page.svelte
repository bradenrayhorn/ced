<script lang="ts">
  import IconComplete from "~icons/mdi/check-circle";
  import { page } from "$app/stores";
  import { env } from "$env/dynamic/public";
</script>

<svelte:head>
  <title>RSVP Complete</title>
</svelte:head>

<div class="w-full text-center">
  <h2 class="h2">RSVP Complete</h2>

  <IconComplete class="mx-auto my-4 text-tertiary-500 text-7xl" />

  <p>Your RSVP has been received!</p>

  <div class="mt-6 flex flex-col gap-4">
    <a class="anchor" href={env.PUBLIC_EVENT_URL}>View event details</a>
    <a class="anchor" href={`/modify/${$page.params.id}`}>Edit RSVP</a>
  </div>
</div>
