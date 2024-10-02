<script lang="ts">
    import { onMount } from 'svelte';

    onMount(() => {
        window.history.replaceState({}, '', '/trader/lots');
        window.dispatchEvent(new Event('popstate'));
    });
</script>