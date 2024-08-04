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
    "One attendee",
    "Two attendees",
    "Three attendees",
    "Four attendees",
    "Five attendees",
    "Six attendees",
    "Seven attendees",
    "Eight attendees",
    "Nine attendees",
    "Ten attendees",
  ].slice(0, group.max_attendees + 1);

  $: isLoading = isSaving || !!$navigating;
</script>

<svelte:head>
  <title>RSVP - {group.name}</title>
</svelte:head>

<h2 class="mt-8 mb-5 h2">{group.name}</h2>

<p class="mb-5">Please confirm the number of attendees in your party.</p>

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
