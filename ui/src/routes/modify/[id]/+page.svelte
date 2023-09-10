<script lang="ts">
  import { enhance } from "$app/forms";
  import { navigating } from "$app/stores";
  import { ProgressRadial } from "@skeletonlabs/skeleton";
  import type { PageData } from "./$types";

  export let data: PageData;
  let isSaving = false;

  $: group = data.data;

  $: options = [
    "Decline to attend",
    "One guest",
    "Two guests",
    "Three guests",
    "Four guests",
    "Five guests",
    "Six guests",
    "Seven guests",
    "Eight guests",
    "Nine guests",
    "Ten guests",
  ].slice(0, group.max_attendees + 1);

  $: isLoading = isSaving || !!$navigating;
</script>

<svelte:head>
  <title>Modify RSVP</title>
</svelte:head>

<h1 class="h1">Modify RSVP</h1>

<h3 class="mt-6 mb-2 h3">{group.name}</h3>

<p class="mb-2">Please confirm the number of guests in your party.</p>

<form
  method="POST"
  action="?/modify"
  autocomplete="off"
  use:enhance={() => {
    isSaving = true;

    return async ({ update }) => {
      await update({ reset: false });
      isSaving = false;
    };
  }}
>
  {#each options as option, index (option)}
    <label class="flex items-center gap-2">
      <input
        class="radio"
        type="radio"
        name="attendees"
        value={index}
        checked={group.attendees === index}
        disabled={isLoading}
      />
      <p>{option}</p>
    </label>
  {/each}

  <button
    class="btn variant-filled w-full mt-6 mb-2 h-10"
    type="submit"
    disabled={isLoading}
  >
    {#if isLoading}
      <ProgressRadial width="w-6" />
    {:else}
      Confirm
    {/if}
  </button>
</form>

<a class="anchor mt-6 text-sm" href="/">Back to search</a>
