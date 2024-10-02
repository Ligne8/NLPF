<script lang="ts">
    import { onMount } from 'svelte';
    import { userRole, currentTab } from '@stores/store';
    import type { UserRole } from '@stores/store';

    // Role permissions
    const rolePermissions: Record<UserRole, string[]> = {
        admin: ['Lots', 'Tractors', 'TrafficManager', 'Trader', 'StockExchange'],
        traffic_manager: ['TrafficManager'],
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

    // Set tab based on URL
    onMount(() => {
        const path = window.location.pathname;
        if (path.startsWith('/traffic-manager'))
            currentTab.set('TrafficManager');
        else if (path.startsWith('/lots'))
            currentTab.set('Lots');
        else if (path.startsWith('/tractors'))
            currentTab.set('Tractors');
        else if (path.startsWith('/trader'))
            currentTab.set('Trader');
        else if (path.startsWith('/stock-exchange'))
            currentTab.set('StockExchange');
        else
            currentTab.set('');
    });
</script>


<!-- Navbar -->
<nav class="bg-gray-800 p-4 text-white">
    <div class="flex items-center">

        <!-- Logo -->
        <a href="/" class="flex-shrink-0">
            <img src="/images/logo.png" alt="Logo" class="h-12 w-auto" />
        </a>
        <a href="/">
            <span class="ml-3 text-xl font-bold tracking-[8px]">
                LIGNE<span class="text-blue-400">8</span>
            </span>
        </a>

        <!-- Navigation Links -->
        <ul class="flex space-x-16 ml-16">
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
    </div>
</nav>
