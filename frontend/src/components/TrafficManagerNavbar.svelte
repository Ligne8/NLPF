<script lang="ts">
    import { currentTrafficManagerTab } from '../stores/store.js';
</script>


<!-- Sub-Navbar -->
<nav class="bg-gray-200 p-2 text-gray-800">
    <ul class="flex space-x-8">
        <li><a href="/traffic-manager/roads" on:click={() => $currentTrafficManagerTab = 'Roads'} class="{$currentTrafficManagerTab === 'Roads' ? 'font-bold' : ''}">Routes</a></li>
        <li><a href="/traffic-manager/lots" on:click={() => $currentTrafficManagerTab = 'Lots'} class="{$currentTrafficManagerTab === 'Lots' ? 'font-bold' : ''}">Lots</a></li>
        <li><a href="/traffic-manager/tractors" on:click={() => $currentTrafficManagerTab = 'Tractors'} class="{$currentTrafficManagerTab === 'Tractors' ? 'font-bold' : ''}">Tracteurs</a></li>
    </ul>
</nav>