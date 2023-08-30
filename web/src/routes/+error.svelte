<script>
  import { page } from "$app/stores";
</script>

<aside class="alert variant-filled-error">
  <div class="alert-message">
    Sorry, an error occurred.
    <p><b>Please retry your previous action.</b></p>
  </div>
</aside>

<a class="btn variant-filled w-full mt-4" href={$page.url.href}>
  Refresh page and try again
</a>

<p class="mt-6">
  <code class="code text-sm">{$page.status}: {$page.error.message}</code>
</p>
