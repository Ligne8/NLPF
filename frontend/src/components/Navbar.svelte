<script lang="ts">
    import { onMount } from 'svelte';

    import { userRole, currentTab } from '@stores/store.js';
    import type { UserRole } from '@stores/store.js';
    import axios from 'axios';

    // Variables
    const API_BASE_URL = import.meta.env.VITE_API_BASE_URL;

    // Role permissions
    const rolePermissions: Record<UserRole, string[]> = {
        admin: ['Lots', 'Tractors', 'TrafficManager', 'Trader', 'StockExchange', 'Map'],
        traffic_manager: ['TrafficManager', 'Map'],
        trader: ['Trader', 'StockExchange'],
        client: ['Lots', 'Tractors', 'StockExchange', 'Map']
    };

    // State for simulation date
    let simulationDate: string = '';
    let error: string = '';

    // Function to format timestamp into DD/MM/YYYY
    const formatDate = (timestamp: number) => {
        const date = new Date(timestamp);
        return `${date.getDate().toString().padStart(2, '0')}/${(date.getMonth() + 1).toString().padStart(2, '0')}/${date.getFullYear()}`;
    };

    // Fetch simulation date from the backend
    async function fetchSimulationDate() {
        try {
            const response = await axios.get(`${API_BASE_URL}/simulations/date`);
            simulationDate = response.data.simulation_date; // Date retrieval in YYYY-MM-DD format
        } catch (err) {
            error = 'Error retrieving simulation date';
            console.error(err);
        }
    }

    async function updateSimulationDate() {
        try {
            await axios.patch(`${API_BASE_URL}/simulations/date`, {});
            await axios.get(`${API_BASE_URL}/simulations/move_tractors`, {});
            await fetchSimulationDate(); // Re-fetch the date after updating
        } catch (err) {
            error = 'Error updating simulation date';
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
        else if (path.startsWith('/map'))
            currentTab.set('Map');
        else
            currentTab.set(''); // Reset if on root or unknown path

        fetchSimulationDate();
    });
</script>

<!-- Centered User Role Title -->
<div class="text-center flex-1 block lg:block bg-yellow-500 fixed top-0 left-0 w-full z-50">
    <p class="text-lg text-gray-800">
        <span class="font-bold">
            {#if $userRole === 'traffic_manager'}
                Traffic Manager
            {:else if $userRole === 'client'}
                Client
            {:else if $userRole === 'trader'}
                Trader
            {:else if $userRole === 'admin'}
                Admin
            {/if}
        </span>
        view
    </p>
</div>

<!-- Navbar -->
<nav class="bg-gray-800 p-4 text-white shadow-md fixed top-0 left-0 w-full z-50 mt-7">
    <div class="flex items-center justify-between">

        <div class="flex items-center">
            <!-- Logo -->
            <a href="/" class="flex-shrink-0">
                <img src="/images/logo.png" alt="Logo" class="h-12 w-auto transition-transform duration-300 hover:scale-105" />
            </a>
            <a href="/">
                <span class="ml-3 text-xl font-bold tracking-widest hover:text-blue-400 transition-colors duration-300">
                    LIGNE<span class="text-blue-400">8</span>
                </span>
            </a>

            <!-- Navigation Links -->
            <ul class="flex space-x-8 ml-16 text-base">
                {#if hasAccess('Lots')}
                    <li>
                        <a href="/lots" on:click={() => switchTab('Lots')} class="{$currentTab === 'Lots' ? 'font-bold text-blue-400' : 'hover:text-blue-400 transition-colors duration-300'}">
                            Lots
                        </a>
                    </li>
                {/if}
                {#if hasAccess('Tractors')}
                    <li>
                        <a href="/tractors" on:click={() => switchTab('Tractors')} class="{$currentTab === 'Tractors' ? 'font-bold text-blue-400' : 'hover:text-blue-400 transition-colors duration-300'}">
                            Tractors
                        </a>
                    </li>
                {/if}
                {#if hasAccess('TrafficManager')}
                    <li>
                        <a href="/traffic-manager" on:click={() => switchTab('TrafficManager')} class="{$currentTab === 'TrafficManager' ? 'font-bold text-blue-400' : 'hover:text-blue-400 transition-colors duration-300'}">
                            Traffic manager
                        </a>
                    </li>
                {/if}
                {#if hasAccess('Trader')}
                    <li>
                        <a href="/trader" on:click={() => switchTab('Trader')} class="{$currentTab === 'Trader' ? 'font-bold text-blue-400' : 'hover:text-blue-400 transition-colors duration-300'}">
                            Trader
                        </a>
                    </li>
                {/if}
                {#if hasAccess('StockExchange')}
                    <li>
                        <a href="/stock-exchange" on:click={() => switchTab('StockExchange')} class="{$currentTab === 'StockExchange' ? 'font-bold text-blue-400' : 'hover:text-blue-400 transition-colors duration-300'}">
                            Stock exchange
                        </a>
                    </li>
                {/if}
                {#if hasAccess('Map')}
                    <li>
                        <a href="/map" on:click={() => switchTab('Map')} class="{$currentTab === 'Map' ? 'font-bold text-blue-400' : 'hover:text-blue-400 transition-colors duration-300'}">
                            Map
                        </a>
                    </li>
                {/if}
            </ul>
        </div>

        <!-- Right section -->
        <div class="flex items-center space-x-4">
            {#if simulationDate}
                <div class="text-lg text-white">
                    Simulation date : <span class="font-bold text-blue-400">{formatDate(simulationDate)}</span>
                </div>
            {:else}
                <div class="text-lg font-normal text-red-600">
                    No date found
                </div>
            {/if}

            <!-- Update simulation date button -->
            <button class="bg-blue-500 hover:bg-blue-700 text-white font-bold rounded-full transition-all duration-300 transform hover:scale-105 w-10 h-10 flex items-center justify-center"
                    on:click={updateSimulationDate}
            >
                <i class="fas fa-plus"></i>
            </button>            

            <!-- Logout button -->
            <button on:click={() => {
                localStorage.clear();
                window.location.reload();
            }} class="bg-red-500 hover:bg-red-700 text-white font-bold py-2 px-4 rounded-md transition-all duration-300 transform hover:scale-105">
                <i class="fa-solid fa-right-from-bracket"></i>
            </button>
        </div>
    </div>
</nav>
