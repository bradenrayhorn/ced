<script lang="ts">
  import { enhance } from "$app/forms";
  import { page, navigating } from "$app/stores";
  import { env } from "$env/dynamic/public";
  import { ProgressRadial } from "@skeletonlabs/skeleton";
  import type { PageData } from "./$types";

  export let data: PageData;

  let isProcessing = false;

  $: search = $page.url.searchParams.get("search") ?? "";
  $: groups = data.foundGroups;
  $: isLoading = !!$navigating || isProcessing;
</script>

<svelte:head>
  <title>RSVP for {env.PUBLIC_EVENT_TITLE}</title>
</svelte:head>

<h2 class="h2">RSVP for</h2>

<h1 class="h1 mt-9 mb-9">
  {env.PUBLIC_EVENT_TITLE}
</h1>

{#if search && groups.length > 0}
  {#if groups.length === 1}
    <p>Is this you?</p>
  {:else}
    <p>Is one of these options you?</p>
  {/if}

  {#if groups.length === 1}
    <p class="my-3"><b>{groups[0].name}</b></p>

    <div class="flex gap-4">
      {#if !isLoading}
        <a
          class="btn-sm variant-filled-primary"
          href={`/modify/${groups[0].id}`}
        >
          Yes
        </a>
        <a class="btn-sm variant-ghost-primary" href="/">No</a>
      {:else}
        <ProgressRadial width="w-6" />
      {/if}
    </div>
  {:else}
    <form
      class="my-6"
      method="POST"
      action="?/toModify"
      use:enhance={() => {
        isLoading = true;

        return async ({ update }) => {
          await update({ reset: false });
          isLoading = false;
        };
      }}
    >
      {#each groups as group, index (group.id)}
        <label class="flex items-center gap-2 mb-2">
          <input
            class="radio"
            type="radio"
            name="group"
            value={group.id}
            checked={index === 0}
          />
          <p>{group.name}</p>
        </label>
      {/each}

      <label class="flex items-center gap-2 mb-2">
        <input class="radio" type="radio" name="group" value={"0"} />
        <p>None of the above</p>
      </label>

      <button
        type="submit"
        class="btn variant-filled shrink-0 mt-6"
        disabled={isLoading}
      >
        Continue
      </button>
    </form>
  {/if}
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
        disabled={isLoading}
      />
      <button
        type="submit"
        class="btn variant-filled shrink-0"
        disabled={isLoading}
      >
        Search
      </button>
    </div>
  </form>

  {#if search && groups.length === 0}
    <aside class="alert variant-filled-error mt-2">
      <div class="alert-message">Please refine your search and try again.</div>
    </aside>
  {/if}
{/if}
