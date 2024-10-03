<script lang="ts">
    import { onMount } from 'svelte';

    onMount(() => {
        window.history.replaceState({}, '', '/stock-exchange/lots');
        window.dispatchEvent(new Event('popstate'));
    });
</script>