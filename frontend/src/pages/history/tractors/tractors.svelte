<script lang="ts">
    import {onMount} from 'svelte';
    import Navbar from '@components/Navbar.svelte';
    import {userId} from "@stores/store";
    import axios from "axios";
    import type { Tractor } from 'src/interface/tractorInterface';
    import HistoryNavbar from '@components/HistoryNavbar.svelte';

    let title: string = 'History';
    let subtitle: string = 'Find your tractor bid history.';
    const API_BASE_URL = import.meta.env.VITE_API_BASE_URL;
    let tractors: Tractor[] = [];
    let selectedState: string = 'all';
    let sortOption: string = 'none';

    // Function to format timestamp into DD/MM/YYYY
    const formatDate = (timestamp: number) => {
        const date = new Date(timestamp);
        return `${date.getDate().toString().padStart(2, '0')}/${(date.getMonth() + 1).toString().padStart(2, '0')}/${date.getFullYear()}`;
    };

    // Fetch tractors bid
    async function fetchTractors() {
        await axios.get(`${API_BASE_URL}/tractors/bids/${$userId}`)
            .then((response) => {
                tractors = response.data;
                console.log(tractors)
            }).catch((error) => {
                console.error('Error fetching tractors:', error.response);
            });
    }

    // Fetch all data
    onMount(async () => {
        await fetchTractors();
    });

    // Update data depending on filters
    $: sortedData = (() => {
        let data = selectedState === 'all' ? tractors : tractors.filter(tractor => tractor.state === selectedState);

        switch (sortOption) {
            case 'min_price_asc':
                return data.sort((a, b) => a.min_price_by_km - b.min_price_by_km);
            case 'min_price_desc':
                return data.sort((a, b) => b.min_price_by_km - a.min_price_by_km);
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
                <option value="min_price_asc">Maximum price (Ascending)</option>
                <option value="min_price_desc">Maximum price (Descending)</option>
                <option value="current_price_asc">Current price (Ascending)</option>
                <option value="current_price_desc">Current price (Descending)</option>
            </select>

        </div>

        <div class="flex justify-between items-center self-end">

            <!-- Reload button -->
            <button class="bg-gray-800 text-white font-bold px-4 py-2 rounded flex items-center hover:bg-gray-900 transition-colors self-end"
                    on:click={fetchTractors}
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
                    <td class="border p-2 text-center">{row.state}</td>

                    <!-- Column 2 -->
                    <td class="border p-2 text-center">{formatDate(row.offer.limit_date)}</td>

                    <!-- Column 3 -->
                    <td class="border p-2 text-center">{row.min_price_by_km.toFixed(2)}</td>

                    <!-- Column 4 -->
                    <td class="border p-2 text-center">{row.current_price.toFixed(2)}</td>

                </tr>
            {/each}
            </tbody>
        </table>
    </div>
</main>