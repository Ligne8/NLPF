<script lang="ts">
    import { onMount } from 'svelte';
    import { currentTrafficManagerTab } from '../stores/store.js';

    // Function to switch tab
    function switchTab(tab: string) {
        currentTrafficManagerTab.set(tab);
    }

    // Set tab based on URL on mount
    onMount(() => {
        const path = window.location.pathname;
        if (path.includes('trader/lots'))
            currentTrafficManagerTab.set('Lots');
        else if (path.includes('trader/tractors'))
            currentTrafficManagerTab.set('Tractors');
        else
            currentTrafficManagerTab.set(''); // Reset if on an unknown path
    });
</script>


<!-- Sub-Navbar -->
<nav class="bg-gray-200 px-10 py-4 text-gray-800 fixed top-24 mt-3 left-0 w-full z-40">
    <ul class="flex space-x-16">
        <li>
            <a
                    href="/trader/lots"
                    on:click={() => switchTab('Lots')}
                    class="{$currentTrafficManagerTab === 'Lots' ? 'font-bold bg-gray-800 text-gray-200 rounded-md px-4 py-2' : ''}">Lots
            </a>
        </li>
        <li>
            <a
                    href="/trader/tractors"
                    on:click={() => switchTab('Tractors')}
                    class="{$currentTrafficManagerTab === 'Tractors' ? 'font-bold bg-gray-800 text-gray-200 rounded-md px-4 py-2' : ''}">Tractors
            </a>
        </li>
    </ul>
</nav>
