<script lang="ts">
    import Navbar from '@components/Navbar.svelte';
    import TraderNavbar from '@components/TraderNavbar.svelte';
    import axios from 'axios';
    import {userId} from "@stores/store";
    import { onMount } from 'svelte';
    import type { Tractor } from 'src/interface/tractorInterface';

    // Variables
    let title: string = 'Tractor offers';
    let subtitle: string = 'Create tractor offers in real time.';
    const API_BASE_URL = import.meta.env.VITE_API_BASE_URL;
    let selectedStatus: string = 'all';
    let sortOption: string = 'none';
    let tractors: Tractor[] = [];

    // Function to format timestamp into DD/MM/YYYY
    const formatDate = (timestamp: number) => {
        const date = new Date(timestamp);
        return `${date.getDate().toString().padStart(2, '0')}/${(date.getMonth() + 1).toString().padStart(2, '0')}/${date.getFullYear()}`;
    };

    // Function to get tag color and text based on status
    function getStateInfo(state: string): { color: string; text: string } {
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
            case 'archive':
                return { color: 'bg-gray-200 text-gray-800', text: '◉ Archived' };
            default:
                return { color: 'bg-gray-200 text-gray-800', text: '🛇 Unknown' };
        }
    }

    // Fetch all tractors of the trader
    async function fetchTractors() {
        try {
            const response = await axios.get(`${API_BASE_URL}/tractors/trader/${$userId}`);
            tractors = response.data;
            console.log(tractors)
        } catch (err) {
            console.error(err);
        }
    }

    // Fetch all data
    onMount(async () => {
        await fetchTractors();
    });

    // Update data depending on filters
    $: sortedData = (() => {
        let data = selectedStatus === 'all' ? tractors : tractors.filter(tractor => tractor.state === selectedStatus);

        switch (sortOption) {
            case 'name_asc':
                return data.sort((a, b) => a.name.localeCompare(b.name));
            case 'name_desc':
                return data.sort((a, b) => b.name.localeCompare(a.name));
            case 'loading_asc':
                return data.sort((a, b) => (a.current_units / a.max_units) - (b.current_units / b.max_units));
            case 'loading_desc':
                return data.sort((a, b) => (b.current_units / b.max_units) - (a.current_units / a.max_units));
            case 'remaining_volume_asc':
                return data.sort((a, b) => (a.max_units - a.current_units) - (b.max_units - b.current_units));
            case 'remaining_volume_desc':
                return data.sort((a, b) => (b.max_units - b.current_units) - (a.max_units - a.current_units));
            case 'location_asc':
                return data.sort((a, b) => a.current_checkpoint.name.localeCompare(b.current_checkpoint.name));
            case 'location_desc':
                return data.sort((a, b) => b.current_checkpoint.name.localeCompare(a.current_checkpoint.name));
            default:
                return data.sort((a, b) => a.name.localeCompare(b.name));
        }
    })();

    const publish = (rowId: string)=>{
      axios.patch(`${API_BASE_URL}/tractors/updateState`, {id: rowId, state: "on_market"}).then((response)=>{
        fetchTractors();
      }).catch((err)=>{
        console.error(err);
      });

    }

</script>


<!-- Navbar -->
<Navbar/>
<TraderNavbar/>

<main class="p-10 mt-40">

    <!-- Title and subtitle -->
    <section>
        <h1 class="text-4xl font-bold mb-4">{title}</h1>
        <h2 class="text-2xl mb-8 text-gray-600">{subtitle}</h2>
    </section>

    <!-- Table -->
    <div>
        <table class="table-auto w-full border-collapse border border-gray-300">
            <thead>
                <tr class="bg-gray-100">
                    <th class="border p-2 text-center">Status</th>
                    <th class="border p-2 text-center">Expiration date</th>
                    <th class="border p-2 text-center">Type</th>
                    <th class="border p-2 text-center">Loading<br><span class="font-normal">(in m³)</span></th>
                    <th class="border p-2 text-center">Departure / Arrival</th>
                    <th class="border p-2 text-center">Minimum price<br><span class="font-normal">(in €/km)</span></th>
                    <th class="border p-2 text-center">Current price<br><span class="font-normal">(in €/km)</span></th>
                    <th class="border p-2 text-center">Actions</th>
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
                        <td class="border p-2 text-center">
                            {formatDate(row.offer.limit_date)}
                        </td>
    
                        <!-- Column 3 -->
                        <td class="border p-2 text-center">{row.resource_type}</td>
                        
                        <!-- Column 4 -->
                        <td class="border p-2 text-center">{row.current_units}/{row.max_units}</td>

                        <!-- Column 5 -->
                        <td class="border p-2 text-center">
                            {row.start_checkpoint.name} / {row.end_checkpoint.name}
                        </td>
                        
                        <!-- Column 6 -->
                        <td class="border p-2 text-center">{row.min_price_by_km.toFixed(2)}</td>
                        
                        <!-- Column 7 -->
                        <td class="border p-2 text-center">{row.current_price.toFixed(2)}</td>
                        
                        <!-- Column 8 -->
                        <td class="border p-2 text-center">
                            <div class="flex flex-wrap justify-center space-x-2 space-y-2">
                                {#if row.state === "at_trader"}
                                    <button on:click={()=>{publish(row.id)}} class="bg-blue-200 text-blue-800 px-4 py-2 flex items-center font-bold hover:bg-blue-300 transition-colors rounded-md">
                                        <i class="fas fa-plus mr-2"></i>
                                        Publish
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