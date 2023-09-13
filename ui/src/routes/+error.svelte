<script>
  import { page } from "$app/stores";
</script>

{#if $page.status === 404}
  <p>Sorry, the requested page was not found.</p>

  <p class="mt-6">
    <a class="anchor" href="/">Back to home</a>
  </p>
{:else}
  <aside class="alert variant-filled-error">
    <div class="alert-message">
      <p>Sorry, an error occurred.</p>
      <p>
        <code class="font-mono">{$page.status}: {$page.error.message}</code>
      </p>
      <p><b>Please retry your previous action.</b></p>
    </div>
  </aside>

  <a class="btn variant-filled w-full mt-4" href={$page.url.href}>
    Refresh page and try again
  </a>
{/if}
