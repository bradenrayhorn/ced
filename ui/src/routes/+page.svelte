<script lang="ts">
  import { page, navigating } from "$app/stores";
  import { env } from "$env/dynamic/public";
  import { ProgressRadial } from "@skeletonlabs/skeleton";
  import type { PageData } from "./$types";

  export let data: PageData;

  $: search = $page.url.searchParams.get("search") ?? "";
  $: group = data.foundGroup;
  $: isNavigating = !!$navigating;
</script>

<svelte:head>
  <title>RSVP for {env.PUBLIC_EVENT_TITLE}</title>
</svelte:head>

<h2 class="h2">RSVP for</h2>

<h1 class="h1 mt-9 mb-9">
  {env.PUBLIC_EVENT_TITLE}
</h1>

{#if search && !!group}
  <p>Is this you?</p>

  <p class="my-3"><b>{group.name}</b></p>

  <div class="flex gap-4">
    {#if !isNavigating}
      <a class="btn-sm variant-filled-primary" href={`/modify/${group.id}`}>
        Yes
      </a>
      <a class="btn-sm variant-ghost-primary" href="/">No</a>
    {:else}
      <ProgressRadial width="w-6" />
    {/if}
  </div>
{:else}
  <form action="?" autocomplete="off">
    <p class="mb-2">Please enter your full name to search for your RSVP.</p>

    <div class="flex gap-2">
      <input
        class="input"
        type="text"
        title="Search"
        name="search"
        value={search}
        disabled={isNavigating}
      />
      <button
        type="submit"
        class="btn variant-filled shrink-0"
        disabled={isNavigating}
      >
        Search
      </button>
    </div>
  </form>

  {#if search && !group}
    <aside class="alert variant-filled-error mt-2">
      <div class="alert-message">Please refine your search and try again.</div>
    </aside>
  {/if}
{/if}
