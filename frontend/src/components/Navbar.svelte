<script lang="ts">
    import { userRole, currentTab } from '../stores/store.js';
    import type { UserRole } from '../stores/store.js';
  
    // Role permissions
    const rolePermissions: Record<UserRole, string[]> = {
        admin: ['Lots', 'Tractors', 'TrafficManager', 'Trader', 'StockExchange'],
        trafficManager: ['TrafficManager'],
        trader: ['Trader', 'StockExchange'],
        client: ['Lots', 'Tractors', 'StockExchange']
    };
  
    // Function to check user access
    function hasAccess(tab: string): boolean {
        const role = $userRole;
        return rolePermissions[role].includes(tab);
    }
  
    // Function to switch tab
    function switchTab(tab: string) {
        currentTab.set(tab);
    }
</script>


<!-- Navbar -->
<nav class="bg-gray-800 p-4 text-white">
    <ul class="flex space-x-8">
        {#if hasAccess('Lots')}
            <li><a href="/lots" on:click={() => switchTab('Lots')} class="{$currentTab === 'Lots' ? 'font-bold' : ''}">Lots</a></li>
        {/if}
        {#if hasAccess('Tractors')}
            <li><a href="/tractors" on:click={() => switchTab('Tractors')} class="{$currentTab === 'Tractors' ? 'font-bold' : ''}">Tracteurs</a></li>
        {/if}
        {#if hasAccess('TrafficManager')}
            <li><a href="/traffic-manager" on:click={() => switchTab('TrafficManager')} class="{$currentTab === 'TrafficManager' ? 'font-bold' : ''}">Traffic manager</a></li>
        {/if}
        {#if hasAccess('Trader')}
            <li><a href="/trader" on:click={() => switchTab('Trader')} class="{$currentTab === 'Trader' ? 'font-bold' : ''}">Trader</a></li>
        {/if}
        {#if hasAccess('StockExchange')}
            <li><a href="/stock-exchange" on:click={() => switchTab('StockExchange')} class="{$currentTab === 'StockExchange' ? 'font-bold' : ''}">Bourse</a></li>
        {/if}
    </ul>
</nav>  