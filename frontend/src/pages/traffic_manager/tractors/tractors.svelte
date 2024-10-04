<script lang="ts">
    import Navbar from '@components/Navbar.svelte';
    import TrafficManagerNavbar from '@components/TrafficManagerNavbar.svelte';
    import {onMount} from "svelte";
    import type {Tractor} from "../../../interface/tractorInterface";
    import {userId, userRole} from "@stores/store";
    import axios from "axios";

    // Variables
    let title: string = 'Tractor management';
    let subtitle: string = 'Track the status of your tractors in real time.';
    const API_BASE_URL = import.meta.env.VITE_API_BASE_URL;
    let tractors: Tractor[]  = [];

    // Function to get tag color and text based on status
    function getStatusInfo(status: string): { color: string; text: string } {
        switch (status) {
            case 'available':
                return { color: 'bg-green-200 text-green-800', text: '◉ Available' };
            case 'in_transit':
                return { color: 'bg-orange-200 text-orange-800', text: '◉ On the way' };
            case 'on_market':
                return { color: 'bg-blue-200 text-blue-800', text: '◉ On the stock exchange' };
            default:
                return { color: 'bg-gray-200 text-gray-800', text: '🛇 Unknown' };
        }
    }

    onMount(() => {
        fetchTableInfo();
    });

    async function fetchTableInfo() {
        if($userRole !== "traffic_manager") {
            return;
        }
        await axios.get(`${API_BASE_URL}/tractors/trafficManager/${$userId}`)
            .then((response) => {
                tractors = response.data;
            }).catch((error) => {
                console.error('Error fetching tractors:', error.response);
            });
    }
</script>


<!-- Navbar -->
<Navbar/>
<TrafficManagerNavbar/>

<main class="p-10 mt-40">

    <section class="flex justify-between items-center mb-4">

        <!-- Title and subtitle -->
        <div>
            <h1 class="text-4xl font-bold mb-2">{title}</h1>
            <h2 class="text-2xl text-gray-600">{subtitle}</h2>
        </div>

        <button class="bg-gray-800 text-white font-bold px-4 py-2 rounded flex items-center hover:bg-gray-900 transition-colors self-end"
            on:click={fetchTableInfo}
        >
            <i class="fas fa-rotate-right mr-2"></i>
            Reload
        </button>

    </section>

    <!-- Table -->
    <div>
        <table class="table-auto w-full border-collapse border border-gray-300">
            <thead>
            <tr class="bg-gray-100">
                <th class="border p-2 text-center">Name</th>
                <th class="border p-2 text-center">Status</th>
                <th class="border p-2 text-center">Loading <span class="font-normal">(in m³)</span></th>
                <th class="border p-2 text-center">Location</th>
                <th class="border p-2 text-center">Route</th>
                <th class="border p-2 text-center">Actions</th>
            </tr>
            </thead>
            <tbody>
            {#each tractors as row, index}
                <tr class={index % 2 === 0 ? 'bg-gray-50' : 'bg-white'}>

                    <!-- Column 1 -->
                    <td class="border p-2 text-center">{row.id}</td>

                    <!-- Column 2 -->
                    <td class="border p-2 text-center">
                            <span class={`px-2 py-1 rounded ${getStatusInfo(row.state).color}`}>
                                {getStatusInfo(row.state).text}
                            </span>
                    </td>

                    <!-- Column 3 -->
                    <td class="border p-2 text-center">{row.current_units}/{row.max_units}</td>

                    <!-- Column 4 -->
                    <td class="border p-2 text-center">{row.current_checkpoint_id}</td>

                    <!-- Column 5 -->
                    <td class="border p-2 text-center">
                        {#if row.state === 'in_transit'}
                            <div class="flex flex-wrap justify-center space-x-2">
                                <button class="bg-red-200 text-red-600 px-4 py-2 flex items-center font-bold hover:bg-red-300 transition-colors rounded-md">
                                    <i class="fas fa-hand mr-2"></i>
                                    Stop
                                </button>
                            </div>
                        {:else if row.state === 'available'}
                            <div class="flex flex-wrap justify-center space-x-2">
                                <button class="bg-green-200 text-green-800 px-4 py-2 flex items-center font-bold hover:bg-green-300 transition-colors rounded-md">
                                    <i class="fas fa-truck mr-2"></i>
                                    Start
                                </button>
                                <button class="bg-blue-200 text-blue-800 px-4 py-2 flex items-center font-bold hover:bg-blue-300 transition-colors rounded-md">
                                    <i class="fas fa-plus mr-2"></i>
                                    Stock exchange
                                </button>
                                <button class="bg-gray-800 text-white px-4 py-2 flex items-center font-bold hover:bg-black transition-colors rounded-md">
                                    <i class="fas fa-right-from-bracket mr-2"></i>
                                    Remove
                                </button>
                            </div>
                        {:else}
                            <span class="text-gray-500">-</span>
                        {/if}
                    </td>
                </tr>
            {/each}
            </tbody>
        </table>
    </div>
</main>