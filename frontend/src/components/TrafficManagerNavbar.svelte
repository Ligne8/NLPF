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
        if (path.includes('traffic-manager/routes'))
            currentTrafficManagerTab.set('Routes');
        else if (path.includes('traffic-manager/lots'))
            currentTrafficManagerTab.set('Lots');
        else if (path.includes('traffic-manager/tractors'))
            currentTrafficManagerTab.set('Tractors');
        else
            currentTrafficManagerTab.set(''); // Reset if on an unknown path
    });
</script>


<!-- Sub-Navbar -->
<nav class="bg-gray-200 px-10 py-4 text-gray-800">
    <ul class="flex space-x-16">
        <li>
            <a
                    href="/traffic-manager/routes"
                    on:click={() => switchTab('Routes')}
                    class="{$currentTrafficManagerTab === 'Routes' ? 'font-bold bg-gray-800 text-gray-200 rounded-md px-4 py-2' : ''}">Routes
            </a>
        </li>
        <li>
            <a
                    href="/traffic-manager/lots"
                    on:click={() => switchTab('Lots')}
                    class="{$currentTrafficManagerTab === 'Lots' ? 'font-bold bg-gray-800 text-gray-200 rounded-md px-4 py-2' : ''}">Lots
            </a>
        </li>
        <li>
            <a
                    href="/traffic-manager/tractors"
                    on:click={() => switchTab('Tractors')}
                    class="{$currentTrafficManagerTab === 'Tractors' ? 'font-bold bg-gray-800 text-gray-200 rounded-md px-4 py-2' : ''}">Tractors
            </a>
        </li>
    </ul>
</nav>
