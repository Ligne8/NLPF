<script lang="ts">
    import { onMount } from 'svelte';
    import { userRole, currentTab } from '@stores/store.js';
    import type { UserRole } from '@stores/store.js';
    import axios from 'axios';

    // Role permissions
    const rolePermissions: Record<UserRole, string[]> = {
        admin: ['Lots', 'Tractors', 'TrafficManager', 'Trader', 'StockExchange'],
        traffic_manager: ['TrafficManager'],
        trader: ['Trader', 'StockExchange'],
        client: ['Lots', 'Tractors', 'StockExchange']
    };

    // State for simulation date
    let simulationDate: string = '';
    let error: string = '';
    // Fetch simulation date from the backend
    async function fetchSimulationDate() {
        try {
            const response = await axios.get('http://localhost:8080/api/v1/simulations/date');
            simulationDate = response.data.simulation_date; // Récupération de la date au format YYYY-MM-DD
        } catch (err) {
            error = 'Erreur lors de la récupération de la date de simulation';
            console.error(err);
        }
    }

    async function updateSimulationDate() {
        try {
            await axios.patch('http://localhost:8080/api/v1/simulations/date', {});
            await fetchSimulationDate(); // Re-fetch the date after updating
        } catch (err) {
            error = 'Erreur lors de la mise à jour de la date de simulation';
            console.error(err);
        }
    }


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
            currentTab.set(''); // Reset if on root or unknown path

        fetchSimulationDate();

    });
</script>


<!-- Navbar -->
<nav class="bg-gray-800 p-4 text-white">
    <div class="flex items-center justify-between">
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
    <div class="flex items-center space-x-4">
        {#if simulationDate}
            <!-- Affichage de la date de simulation -->
            <div class="text-lg font-semibold">
                Date de simulation: <span class="text-blue-400">{simulationDate}</span>
            </div>
        {:else}
            <!-- Message d'erreur si la date de simulation ne peut pas être récupérée -->
            <div class="text-lg font-semibold text-red-500">
                Unable to retrieve simulation date
            </div>
        {/if}

        <!-- Bouton + pour mettre à jour la date -->
        <button class="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded-full" on:click={updateSimulationDate}>
            +
        </button>
    </div>
    </div>
</nav>

<style>
    button {
        display: flex;
        justify-content: center;
        align-items: center;
        width: 40px;
        height: 40px;
        border-radius: 50%;
    }
</style>

