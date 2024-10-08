<script lang="ts">
    import Navbar from '@components/Navbar.svelte';
    import TraderNavbar from '@components/TraderNavbar.svelte';
    import { userId } from '@stores/store';
    import axios from 'axios';
    import type { Lot } from 'src/interface/lotInterface';
    import { onMount } from 'svelte';

    // Variables
    let title: string = 'Lot offers';
    let subtitle: string = 'Create lot offers in real time.';
    const API_BASE_URL = import.meta.env.VITE_API_BASE_URL;
    let lots: Lot[] = [];
    let selectedStatus: string = 'all';
    let sortOption: string = 'none';

    // Function to format timestamp into DD/MM/YYYY
    const formatDate = (timestamp: number) => {
        const date = new Date(timestamp);
        return `${date.getDate().toString().padStart(2, '0')}/${(date.getMonth() + 1).toString().padStart(2, '0')}/${date.getFullYear()}`;
    };

    // Function to get tag color and text based on status
    function getStatusInfo(state: string): { color: string; text: string } {
        switch (state) {
            case 'available':
                return {color: 'bg-green-200 text-green-800', text: '◉ Available'};
            case 'pending':
                return {color: 'bg-yellow-200 text-yellow-800', text: '◉ Pending'};
            case 'in_transit':
                return {color: 'bg-orange-200 text-orange-800', text: '◉ In transit'};
            case 'on_market':
                return {color: 'bg-blue-200 text-blue-800', text: '◉ On market'};
            case 'archived':
                return {color: 'bg-gray-200 text-gray-800', text: '◉ Archived'};
            case 'at_trader':
                return {color: 'bg-purple-200 text-purple-800', text: '◉ At trader'};
            default:
                return {color: 'bg-gray-200 text-gray-800', text: '🛇 Unknown'};
        }
    }

    // Fetch all lots of the trader
    async function fetchLots() {
        try {
            const response = await axios.get(`${API_BASE_URL}/lots/trader/${$userId}`);
            lots = response.data;
            console.log(lots);
        } catch (err) {
            console.error(err);
        }
    }

    // Fetch all data
    onMount(async () => {
        await fetchLots();
    });

    // Update data depending on filters
    $: sortedData = (() => {
        let data = selectedStatus === 'all' ? lots : lots.filter(lot => lot.state === selectedStatus);

        switch (sortOption) {
            case 'volume_asc':
                return data.sort((a, b) => a.volume - b.volume);
            case 'volume_desc':
                return data.sort((a, b) => b.volume - a.volume);
            case 'location_asc':
                return data.sort((a, b) => a.current_checkpoint.localeCompare(b.current_checkpoint));
            case 'location_desc':
                return data.sort((a, b) => b.current_checkpoint.localeCompare(a.current_checkpoint));
            default:
                return data;
        }
    })();

</script>


<!-- Navbar -->
<Navbar/>
<TraderNavbar/>

<main class="p-10 mt-40">

    <!-- Title and subtitle -->
    <section class="mb-2">
        <h1 class="text-4xl font-bold mb-2">{title}</h1>
        <h2 class="text-2xl text-gray-600">{subtitle}</h2>
    </section>

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

        <div class="flex justify-between items-center self-end">

            <!-- Reload button -->
            <button class="bg-gray-800 text-white font-bold px-4 py-2 rounded flex items-center hover:bg-gray-900 transition-colors self-end"
                    on:click={fetchLots}
            >
                <i class="fas fa-rotate-right mr-2"></i>
                Reload
            </button>

        </div>

    </section>

    <!-- Table -->
    <div>
        <table class="table-auto w-full border-collapse border border-gray-300">
            <thead>
                <tr class="bg-gray-100">
                    <th class="border p-2 text-center">Status</th>
                    <th class="border p-2 text-center">Expiration date</th>
                    <th class="border p-2 text-center">Type</th>
                    <th class="border p-2 text-center">Volume<br><span class="font-normal">(in m³)</span></th>
                    <th class="border p-2 text-center">Departure / Arrival</th>
                    <th class="border p-2 text-center">Maximum price<br><span class="font-normal">(in €/km)</span></th>
                    <th class="border p-2 text-center">Current price<br><span class="font-normal">(in €/km)</span></th>
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
                        <td class="border p-2 text-center">
                            {#if row.offer}
                                {formatDate(row.offer.limit_date)}
                            {:else}
                                <span class="px-2 py-1 mx-auto w-4/5 block text-gray-500">None</span>
                            {/if}
                        </td>

                        <!-- Column 3 -->
                        <td class="border p-2 text-center">{row.resource_type}</td>
                        
                        <!-- Column 4 -->
                        <td class="border p-2 text-center">{row.volume}</td>

                        <!-- Column 5 -->
                        <td class="border p-2 text-center">
                            {row.start_checkpoint.name} / {row.end_checkpoint.name}
                        </td>
                        
                        <!-- Column 6 -->
                        <td class="border p-2 text-center">{row.max_price_by_km.toFixed(2)}</td>
                        
                        <!-- Column 7 -->
                        <td class="border p-2 text-center">{row.current_price.toFixed(2)}</td>
                        
                        <!-- Column 8 -->
                        <td class="border p-2 text-center">
                            <div class="flex flex-wrap justify-center space-x-2 space-y-2">
                                {#if !row.offer && row.state === "at_trader"}
                                    <button class="bg-blue-200 text-blue-800 px-4 py-2 flex items-center font-bold hover:bg-blue-300 transition-colors rounded-md">
                                        <i class="fas fa-plus mr-2"></i>
                                        Offer
                                    </button>
                                {:else}
                                    <span class="text-gray-500">-</span>
                                {/if}
                            </div>
                        </td>

                    </tr>
                {/each}
            </tbody>
        </table>
    </div>
</main>