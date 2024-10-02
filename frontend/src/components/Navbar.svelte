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
            simulationDate = response.data.simulation_date; // Get datetime format YYYY-MM-DD
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
<nav class="bg-gray-800 p-4 text-white shadow-md">
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
                            Tracteurs
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
                            Bourse
                        </a>
                    </li>
                {/if}
            </ul>
        </div>

        <!-- Centered User Role Title -->
        <div class="text-center flex-1 hidden lg:block">
    <span class="text-lg font-semibold text-gray-300 hover:text-white transition-colors duration-300">
        {#if $userRole === 'traffic_manager'}
            TRAFFIC MANAGER
        {:else if $userRole === 'client'}
            CLIENT
        {:else if $userRole === 'trader'}
            TRADER
        {:else if $userRole === 'admin'}
            ADMIN
        {/if}
    </span>
        </div>

        <!-- Right section: buttons -->
        <div class="flex items-center space-x-4">
            {#if simulationDate}
                <div class="text-lg font-semibold text-gray-300">
                    Simulation : <span class="text-blue-400">{simulationDate}</span>
                </div>
            {:else}
                <div class="text-lg font-semibold text-red-500">
                    Erreur de date
                </div>
            {/if}

            <!-- Update Simulation Date Button -->
            <button class="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded-full transition-all duration-300 transform hover:scale-105" on:click={updateSimulationDate}>
                +
            </button>

            <!-- Logout Button -->
            <button on:click={() => {
                localStorage.clear();
                window.location.reload();
            }} class="bg-red-500 hover:bg-red-700 text-white font-bold py-2 px-4 rounded-lg transition-all duration-300 transform hover:scale-105">
                <i class="fa-solid fa-right-from-bracket"></i>
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

    /* Hover and transition effects */
    a {
        transition: color 0.3s ease;
    }

    a:hover {
        color: #60a5fa; /* Tailwind's blue-400 */
    }

    .rounded-lg {
        border-radius: 8px;
    }

    /* Flex adjustments */
    .flex-1 {
        flex-grow: 1;
    }

    /* Ensure spacing and margin are balanced */
    .ml-16 {
        margin-left: 4rem;
    }
</style>
