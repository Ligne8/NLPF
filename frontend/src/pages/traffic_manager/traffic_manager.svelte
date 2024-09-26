<script lang="ts">
    import { onMount } from 'svelte';

    onMount(() => {
        window.history.replaceState({}, '', '/traffic-manager/routes');
        window.dispatchEvent(new Event('popstate'));
    });
</script>