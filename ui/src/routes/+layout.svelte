<script>
  import "../app.postcss";
  import { autoModeWatcher } from "@skeletonlabs/skeleton";
  import { env } from "$env/dynamic/public";
</script>

<svelte:head>
  <!-- The \u{73} syntax is needed to avoid a parse error: https://github.com/sveltejs/svelte-eslint-parser/issues/346. -->
  <!-- eslint-disable-next-line svelte/no-at-html-tags -->
  {@html `<\u{73}cript>${autoModeWatcher.toString()} autoModeWatcher(); </script>`}
</svelte:head>

<div
  class="max-w-3xl w-full mx-auto pt-8 px-4 flex flex-col justify-between min-h-screen gap-1"
>
  <div>
    <slot />
  </div>

  <footer class="text-primary-300 text-sm w-full text-center pb-2">
    <a class="anchor" href={env.PUBLIC_EVENT_URL}>{env.PUBLIC_EVENT_TITLE}</a>
  </footer>
</div>
