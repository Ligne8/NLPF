<script lang="ts">
    import { onMount } from 'svelte';

    onMount(() => {
        window.history.replaceState({}, '', '/history/lots');
        window.dispatchEvent(new Event('popstate'));
    });
</script>