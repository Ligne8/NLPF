<script lang="ts">
    import {onMount} from 'svelte';
    import Navbar from '@components/Navbar.svelte';
    import {userId} from "@stores/store";
    import axios from "axios";
    import type { Lot } from 'src/interface/lotInterface';
    import HistoryNavbar from '@components/HistoryNavbar.svelte';

    let title: string = 'History of lots bids';
    let subtitle: string = 'Find your lot bid history.';
    const API_BASE_URL = import.meta.env.VITE_API_BASE_URL;
    let lots: Lot[] = [];
    let selectedState: string = 'all';
    let sortOption: string = 'none';

    // Function to format timestamp into DD/MM/YYYY
    const formatDate = (timestamp: number) => {
        const date = new Date(timestamp);
        return `${date.getDate().toString().padStart(2, '0')}/${(date.getMonth() + 1).toString().padStart(2, '0')}/${date.getFullYear()}`;
    };

    // Function to get tag color and text based on status
    function getStateInfo(state: string): { color: string; text: string } {
        switch (state) {
            case 'accepted':
                return {color: 'bg-green-200 text-green-800', text: '◉ Accepted'};
            case 'rejected':
                return {color: 'bg-red-200 text-red-800', text: '◉ Rejected'};
            case 'in_progress':
                return {color: 'bg-blue-200 text-blue-800', text: '◉ In progress'};
            case 'denied':
                return {color: 'bg-orange-200 text-orange-800', text: '◉ Denied'};
            default:
                return {color: 'bg-gray-200 text-gray-800', text: '🛇 Unknown'};
        }
    }

    // Fetch lots bid
    async function fetchLots() {
        await axios.get(`${API_BASE_URL}/lots/bids/${$userId}`)
            .then((response) => {
                lots = response.data;
            }).catch((error) => {
                console.error('Error fetching lots:', error.response);
            });
    }

    // Fetch all data
    onMount(async () => {
        await fetchLots();
    });

    // Update data depending on filters
    $: sortedData = (() => {
        let data = selectedState === 'all' ? lots : lots.filter(lot => lot.state === selectedState);

        switch (sortOption) {
            case 'max_price_asc':
                return data.sort((a, b) => a.max_price_by_km - b.max_price_by_km);
            case 'max_price_desc':
                return data.sort((a, b) => b.max_price_by_km - a.max_price_by_km);
            case 'current_price_asc':
                return data.sort((a, b) => b.current_price - a.current_price);
            case 'current_price_desc':
                return data.sort((a, b) => b.current_price - a.current_price);
            default:
                return data;
        }
    })();

</script>


<!-- Navbar -->
<Navbar/>
<HistoryNavbar/>

<main class="p-10 mt-40">

    <!-- Title and subtitle -->
    <div class="mb-2">
        <h1 class="text-4xl font-bold mb-2">{title}</h1>
        <h2 class="text-2xl text-gray-600">{subtitle}</h2>
    </div>

    <section class="flex justify-between items-center mb-4">

        <div class="flex justify-between items-center self-end">

            <!-- Filter by state -->
            <select bind:value={selectedState} class="mr-2 border border-gray-300 rounded px-2 py-1">
                <option value="all" disabled selected>Filter by state</option>
                <option value="all">All</option>
                <option value="accepted">Accepted</option>
                <option value="rejected">Rejected</option>
                <option value="in_progress">In progress</option>
                <option value="denied">Denied</option>
            </select>

            <!-- Sort by volume and location -->
            <select bind:value={sortOption} class="border border-gray-300 rounded px-2 py-1">
                <option value="none" disabled selected>Sort by</option>
                <option value="max_price_asc">Maximum price (Ascending)</option>
                <option value="max_price_desc">Maximum price (Descending)</option>
                <option value="current_price_asc">Current price (Ascending)</option>
                <option value="current_price_desc">Current price (Descending)</option>
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
                <th class="border p-2 text-center">State</th>
                <th class="border p-2 text-center">Date</th>
                <th class="border p-2 text-center">Maximum price<br><span class="font-normal">(in €/km)</span></th>
                <th class="border p-2 text-center">Current price<br><span class="font-normal">(in €/km)</span></th>
            </tr>
            </thead>
            <tbody>
            {#each sortedData as row, index}
                <tr class={index % 2 === 0 ? 'bg-gray-50' : 'bg-white'}>

                    <!-- Column 1 -->
                    <td class="border p-2 text-center">
                        <span class={`px-2 py-1 rounded ${getStateInfo(row.state).color}`}>
                            {getStateInfo(row.state).text}
                        </span>
                    </td>

                    <!-- Column 2 -->
                    <td class="border p-2 text-center">{formatDate(row.limit_date)}</td>

                    <!-- Column 3 -->
                    <td class="border p-2 text-center">{row.max_price_by_km.toFixed(2)}</td>

                    <!-- Column 4 -->
                    <td class="border p-2 text-center">{row.current_price.toFixed(2)}</td>

                </tr>
            {/each}
            </tbody>
        </table>
    </div>
</main>