<script lang="ts">

    export let userRole: UserRole;
    export let currentTab: string;
    export let onTabClick: (tab: string) => void;

    // Role permissions
    const rolePermissions: Record<UserRole, string[]> = {
        admin: ['Lots', 'Tractors', 'TrafficManager', 'Trader', 'StockExchange'],
        trafficManager: ['TrafficManager'],
        trader: ['Trader', 'StockExchange'],
        client: ['Lots', 'Tractors', 'StockExchange']
    };

    // Function to check user access
    function hasAccess(tab: string): boolean {
        return rolePermissions[userRole].includes(tab);
    }

</script>


<!-- Navbar -->
<nav class="bg-gray-800 p-4 text-white">
    <ul class="flex space-x-8">
        {#if hasAccess('Lots')}
            <li><a href="#" on:click={() => currentTab = 'Lots'} class="{currentTab === 'Lots' ? 'font-bold' : ''}">Lots</a></li>
        {/if}
        {#if hasAccess('Tractors')}
            <li><a href="#" on:click={() => currentTab = 'Tractors'} class="{currentTab === 'Tractors' ? 'font-bold' : ''}">Tracteurs</a></li>
        {/if}
        {#if hasAccess('TrafficManager')}
            <li><a href="#" on:click={() => currentTab = 'TrafficManager'} class="{currentTab === 'TrafficManager' ? 'font-bold' : ''}">Traffic manager</a></li>
        {/if}
        {#if hasAccess('Trader')}
            <li><a href="#" on:click={() => currentTab = 'Trader'} class="{currentTab === 'Trader' ? 'font-bold' : ''}">Trader</a></li>
        {/if}
        {#if hasAccess('StockExchange')}
            <li><a href="#" on:click={() => currentTab = 'StockExchange'} class="{currentTab === 'StockExchange' ? 'font-bold' : ''}">Bourse</a></li>
        {/if}
    </ul>
</nav>