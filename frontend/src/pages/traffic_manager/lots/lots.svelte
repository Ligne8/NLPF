<script lang="ts">
    import Navbar from '@components/Navbar.svelte';
    import TrafficManagerNavbar from '@components/TrafficManagerNavbar.svelte';
    import { userId, userRole } from "@stores/store";
    import { onMount } from "svelte";
    import axios from "axios";
    import type { Lot } from 'src/interface/lotInterface';
    import type { Tractor } from 'src/interface/tractorInterface';

    // Variables
    let title: string = 'Lot management';
    let subtitle: string = 'Track the status of your lots in real time.';
    const API_BASE_URL = import.meta.env.VITE_API_BASE_URL;
    let lots: Lot[] = [];
    let compatibleTractorsMap: Map<number, Tractor[]> = new Map();
    let isModalOpen = false;
    let selectedLotId: number = null;
    let selectedStatus: string = 'all';
    let sortOption: string = 'none';

    // Function to get tag color and text based on status
    function getStatusInfo(state: string): { color: string; text: string } {
        switch (state) {
            case 'available':
                return { color: 'bg-green-200 text-green-800', text: '◉ Available' };
            case 'pending':
                return { color: 'bg-yellow-200 text-yellow-800', text: '◉ Pending' };
            case 'in_transit':
                return { color: 'bg-orange-200 text-orange-800', text: '◉ In transit' };
            case 'on_market':
                return { color: 'bg-blue-200 text-blue-800', text: '◉ On market' };
            case 'at_trader':
                return { color: 'bg-purple-200 text-purple-800', text: '◉ At trader' };
            case 'archived':
                return { color: 'bg-gray-200 text-gray-800', text: '◉ Archived' };
            default:
                return { color: 'bg-gray-200 text-gray-800', text: '🛇 Unknown' };
        }
    }

    // Function to open modal
    function openModal(lotId: number) {
        selectedLotId = lotId;
        isModalOpen = true;
    }

    // Function to close modal
    function closeModal() {
        isModalOpen = false;
    }

    // Fetch table info
    async function fetchTableInfo() {
        if($userRole !== "traffic_manager") {
            return;
        }
        await axios.get(`${API_BASE_URL}/lots/traffic_manager/${$userId}`)
            .then((response) => {
                lots = response.data;
            }).catch((error) => {
                console.error('Error fetching lots:', error.response);
            });
    }

    // Fetch table info
    async function getCompatibleTractors(lotId: number) {
        if ($userRole !== "traffic_manager") {
            return;
        }

        try {
            const response = await axios.get(`${API_BASE_URL}/lots/tractors/compatible/${$userId}/${lotId}`);
            const updatedMap = new Map(compatibleTractorsMap);
            updatedMap.set(lotId, response.data);
            compatibleTractorsMap = updatedMap;
        } catch (error) {
            console.error('Error fetching compatible tractors:', error.response);
        }
    }

    // Assign tractor to lot
    async function assignTractor(lotId: number, tractorId: number) {
        try {
            await axios.put(`${API_BASE_URL}/lots/assign/tractor`, {
                lot_id: lotId,
                tractor_id: tractorId
            });
            const updatedTractors = compatibleTractorsMap.get(lotId)?.filter(tractor => tractor.id !== tractorId) || [];
            const updatedMap = new Map(compatibleTractorsMap);
            if (updatedTractors.length === 0)
                updatedMap.delete(lotId);
            else
                updatedMap.set(lotId, updatedTractors);
            compatibleTractorsMap = updatedMap;
            if (!updatedTractors.length)
                closeModal();
            await fetchTableInfo();
        } catch (error) {
            console.error('Error assigning tractor:', error.response);
        }
    }

    async function assignLotToTrader(lotId: string) {
         await axios.put(`${API_BASE_URL}/lots/assign/${lotId}/trader`)
            .then((response) => {
                fetchTableInfo();
            }).catch((error) => {
                console.error('Error assigning lot to trader:', error.response);
            });
    }

    onMount(() => {
        fetchTableInfo();
    });

    // Update compatible tractors map
    $: {
        for (const lot of sortedData) {
            if (!compatibleTractorsMap.has(lot.id))
                getCompatibleTractors(lot.id);
        }
    }

    // Update data depending on filters
    $: sortedData = (() => {
        let data = selectedStatus === 'all' ? lots : lots.filter(lot => lot.state === selectedStatus);

        switch (sortOption) {
            case 'volume_asc':
                return data.sort((a, b) => a.volume - b.volume);
            case 'volume_desc':
                return data.sort((a, b) => b.volume - a.volume);
            case 'location_asc':
                return data.sort((a, b) => a.current_checkpoint.name.localeCompare(b.current_checkpoint.name));
            case 'location_desc':
                return data.sort((a, b) => b.current_checkpoint.name.localeCompare(a.current_checkpoint.name));
            default:
                return data;
        }
    })();

    const assignLotToTractor = async (lotId: number, tractorId: number) => {
        try {
            await axios.post(`${API_BASE_URL}/lots/tractors/assign`, {
              lot_id : lotId,
              tractor_id : tractorId
            });
            closeModal();
            fetchTableInfo();
        } catch (error: any) {
            console.error('Error assigning lot to tractor:', error.response);
        }
    };

</script>


<!-- Navbar -->
<Navbar/>
<TrafficManagerNavbar/>

<main class="p-10 mt-40">

    <!-- Title and subtitle -->
    <div class="mb-2">
        <h1 class="text-4xl font-bold mb-4">{title}</h1>
        <h2 class="text-2xl text-gray-600">{subtitle}</h2>
    </div>

    <section class="flex justify-between items-center mb-4">

        <div class="flex justify-between items-center self-end">

            <!-- Filter by status -->
            <select bind:value={selectedStatus} class="mr-2 border border-gray-300 rounded px-2 py-1">
                <option value="all" disabled selected>Filter by status</option>
                <option value="all">All</option>
                <option value="available">Available</option>
                <option value="pending">Pending</option>
                <option value="in_transit">In transit</option>
                <option value="on_market">On market</option>
                <option value="archived">Archived</option>
            </select>

            <!-- Sort by volume and location -->
            <select bind:value={sortOption} class="border border-gray-300 rounded px-2 py-1">
                <option value="none" disabled selected>Sort by</option>
                <option value="volume_asc">Volume (Ascending)</option>
                <option value="volume_desc">Volume (Descending)</option>
                <option value="location_asc">Location (A-Z)</option>
                <option value="location_desc">Location (Z-A)</option>
            </select>

        </div>

    </section>

    <!-- Table -->
    <div>
        <table class="table-auto w-full border-collapse border border-gray-300">
            <thead>
            <tr class="bg-gray-100">
                <th class="border p-2 text-center">Status</th>
                <th class="border p-2 text-center">Volume <span class="font-normal">(in m³)</span></th>
                <th class="border p-2 text-center">Type</th>
                <th class="border p-2 text-center">Location</th>
                <th class="border p-2 text-center">Departure / Arrival</th>
                <th class="border p-2 text-center">Tractor</th>
                <th class="border p-2 text-center">Actions</th>
            </tr>
            </thead>
            <tbody>
            {#each sortedData as row, index}
                <tr class={index % 2 === 0 ? 'bg-gray-50' : 'bg-white'}>

                    <!-- Column 1 -->
                    <td class="border p-2 text-center">
                            <span class={`px-2 py-1 rounded ${getStatusInfo(row.state).color}`}>
                                {getStatusInfo(row.state).text}
                            </span>
                    </td>

                    <!-- Column 2 -->
                    <td class="border p-2 text-center">{row.volume}</td>

                    <!-- Column 3 -->
                    <td class="border p-2 text-center">{row.resource_type}</td>

                    <!-- Column 4 -->
                    <td class="border p-2 text-center">{row.current_checkpoint.name}</td>

                    <!-- Column 5 -->
                    <td class="border p-2 text-center">
                        {row.start_checkpoint.name} / {row.end_checkpoint.name}
                    </td>

                    <!-- Colonne 6 -->
                    <td class="border p-2 text-center">
                        {#if row.tractor }
                            <span class="px-2 py-1 mx-auto w-4/5 block">
                                {row.tractor.name}
                            </span>
                        {:else}
                            {#if compatibleTractorsMap.get(row.id) && compatibleTractorsMap.get(row.id).length > 0 && row.state === 'pending'}
                                <div class="flex flex-wrap justify-center space-x-2">
                                    <button class="bg-green-200 text-green-800 px-4 py-2 flex items-center font-bold hover:bg-green-300 transition-colors rounded-md"
                                        on:click={() => openModal(row.id)}>
                                        <i class="fas fa-truck mr-2"></i>
                                        Assign
                                    </button>
                                </div>
                            {:else}
                                <span class="px-2 py-1 mx-auto w-4/5 block text-gray-500">None</span>
                            {/if}
                        {/if}
                    </td>

                    <!-- Column 6 -->
                    <td class="border p-2 text-center">
                        {#if row.state === 'pending'}
                            <div class="flex flex-wrap justify-center space-x-2 space-y-2">
                                <button class="bg-blue-200 text-blue-800 px-4 py-2 flex items-center font-bold hover:bg-blue-300 transition-colors rounded-md"
                                on:click={() => assignLotToTrader(row.id)} >
                                    <i class="fas fa-plus mr-2"></i>
                                    Stock exchange
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

<!-- svelte-ignore a11y-no-static-element-interactions -->
<!-- svelte-ignore a11y-no-static-element-interactions -->
<!-- svelte-ignore a11y-click-events-have-key-events -->
<!-- svelte-ignore a11y-label-has-associated-control -->

{#if isModalOpen}

    <div class="fixed inset-0 bg-black bg-opacity-50 flex justify-center items-center z-50" on:click={closeModal}>

        <div class="bg-white p-6 rounded-lg shadow-lg w-4/5" on:click|stopPropagation>

            <!-- Close Button -->
            <button class="absolute top-2 right-2 text-gray-500 hover:text-gray-800" on:click={closeModal}>
                &times;
            </button>

            <!-- Modal Title -->
            <h2 class="text-2xl font-bold mb-6">Assign a tractor</h2>

            <table class="table-auto w-full border-collapse border border-gray-300">
                <thead>
                <tr class="bg-gray-100">
                    <th class="border p-2 text-center">Name</th>
                    <th class="border p-2 text-center">Status <span class="font-normal">(in m³)</span></th>
                    <th class="border p-2 text-center">Location</th>
                    <th class="border p-2 text-center">Route</th>
                    <th class="border p-2 text-center">Actions</th>
                </tr>
                </thead>
                <tbody>
                    {#each compatibleTractorsMap.get(selectedLotId) as tractor}
                        <tr>

                            <!-- Column 1 -->
                            <td class="border p-2 text-center max-w-11">{tractor.name}</td>

                            <!-- Column 2 -->
                            <td class="border p-2 text-center max-w-16">
                                    <span class={`px-2 py-1 rounded ${getStatusInfo(tractor.state).color}`}>
                                        {getStatusInfo(tractor.state).text}
                                    </span>
                            </td>


                            <!-- Column 4 -->
                            <td class="border p-2 text-center">{tractor.current_checkpoint.name}</td>

                            <!-- Column 5 -->
                            <td class="border p-2 text-center">
                                {#if tractor.route}
                                    {tractor.route.name}
                                {:else}
                                    <span class="px-2 py-1 mx-auto w-4/5 block text-gray-500">None</span>
                                {/if}
                            </td>

                            <!-- Column 6 -->
                            <td class="border p-2 text-center">
                                <div class="flex flex-wrap justify-center space-x-2">
                                    <button class="bg-gray-200 text-gray-600 px-4 py-2 flex items-center font-bold hover:bg-green-200 hover:text-green-800 transition-colors rounded-md"
                                        on:click={() => { assignLotToTractor(selectedLotId, tractor.id) }}
                                    >
                                        <i class="fas fa-hand-pointer mr-2 icon-default"></i>
                                        <i class="fas fa-check mr-2 icon-hover hidden"></i>
                                        Select
                                    </button>
                                </div>
                            </td>

                        </tr>
                    {/each}
                </tbody>
            </table>
        </div>
    </div>
{/if}

<style>
    button:hover .icon-default {
        display: none;
    }

    button:hover .icon-hover {
        display: inline-block;
    }

    .icon-hover {
        display: none;
    }
</style>