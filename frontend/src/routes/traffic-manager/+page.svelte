<script lang="ts">
    import { onMount } from 'svelte';
    import { goto } from '$app/navigation';

    onMount(() => {
        goto('/traffic-manager/roads', { replaceState: true });
    });
</script>
