<script lang="ts">
    import { currentTrafficManagerTab } from '../stores/store.js';
</script>


<!-- Sub-Navbar -->
<nav class="bg-gray-200 px-10 py-4 text-gray-800">
    <ul class="flex space-x-16">
        <li>
            <a 
                href="/traffic-manager/roads" 
                on:click={() => $currentTrafficManagerTab = 'Roads'} 
                class="{$currentTrafficManagerTab === 'Roads' ? 'font-bold bg-gray-800 text-gray-200 rounded-md px-4 py-2' : ''}">Routes
            </a>
        </li>
        <li>
            <a 
                href="/traffic-manager/lots" 
                on:click={() => $currentTrafficManagerTab = 'Lots'} 
                class="{$currentTrafficManagerTab === 'Lots' ? 'font-bold bg-gray-800 text-gray-200 rounded-md px-4 py-2' : ''}">Lots
            </a>
        </li>
        <li>
            <a 
                href="/traffic-manager/tractors" 
                on:click={() => $currentTrafficManagerTab = 'Tractors'} 
                class="{$currentTrafficManagerTab === 'Tractors' ? 'font-bold bg-gray-800 text-gray-200 rounded-md px-4 py-2' : ''}">Tracteurs
            </a>
        </li>
    </ul>
</nav>
