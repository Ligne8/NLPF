<script lang="ts">
    import Navbar from '@components/Navbar.svelte';
    import StockExchangeNavbar from '@components/StockExchangeNavbar.svelte';
    import { userRole } from '@stores/store';
    import axios from 'axios';
    import type { Lot } from 'src/interface/lotInterface';
    import { onMount } from 'svelte';

    // Variables
    let title: string = 'Lot market';
    let subtitle: string = 'Explore a wide selection of lots with dynamic volumes and prices.';
    const API_BASE_URL = import.meta.env.VITE_API_BASE_URL;
    let isModalOpen = false;
    let priceValue: number = 1.0;
    let minPriceValue: number = 1.0;
    let maxPriceValue: number = 10.0;
    let minVolumeValue: number = 1.0;
    let current_offer_id: string = '';
    let lots: Lot[] = [];
    let selectedStatus: string = 'all';
    let sortOption: string = 'none';

    // Function to format timestamp into DD/MM/YYYY
    const formatDate = (timestamp: number) => {
        const date = new Date(timestamp);
        return `${date.getDate().toString().padStart(2, '0')}/${(date.getMonth() + 1).toString().padStart(2, '0')}/${date.getFullYear()}`;
    };


    // Function to open tractors modal
    function openModal(currentPrice: number, offer_id: string, max_price_by_km: number) {
        priceValue = currentPrice;
        current_offer_id = offer_id;
        minPriceValue = 1;
        maxPriceValue = max_price_by_km;
        isModalOpen = true;
    }

    // Function to close tractors modal
    function closeModal() {
        isModalOpen = false;
    }

    // Function to bid
    function bid() {
        const payload = {
            bid: priceValue,
            offer_id: current_offer_id
        };
        axios.post(`${API_BASE_URL}/stock_exchange/lot/bid`, payload).then((response) => {
          fetchLots();
          closeModal();
        }).catch((error) => {
            console.error('Error bidding:', error.response);
        });
    }

    // Fetch table info
    async function fetchLots() {
        await axios.get(`${API_BASE_URL}/stock_exchange/lot_offers`)
            .then((response) => {
                lots = response.data;
            }).catch((error) => {
                console.error('Error fetching lots:', error.response);
            });
    }

    // Fetch all data
    onMount(() => {
        fetchLots();
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
                return data.sort((a, b) => a.current_checkpoint.name.localeCompare(b.current_checkpoint.name));
            case 'location_desc':
                return data.sort((a, b) => b.current_checkpoint.name.localeCompare(a.current_checkpoint.name));
            default:
                return data;
        }
    })();

</script>


<!-- Navbar -->
<Navbar/>
<StockExchangeNavbar/>

<main class="p-10 mt-40">

    <!-- Title and subtitle -->
    <div class="mb-2">
        <h1 class="text-4xl font-bold mb-2">{title}</h1>
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

    <table class="table-auto w-full border-collapse border border-gray-300">
        <thead>
            <tr class="bg-gray-100">
                <th class="border p-2 text-center">Expiration date</th>
                <th class="border p-2 text-center">Type</th>
                <th class="border p-2 text-center">Volume<br><span class="font-normal">(in m³)</span></th>
                <th class="border p-2 text-center">Maximum price<br><span class="font-normal">(in €/km)</span></th>
                <th class="border p-2 text-center">Minimum Bid<br><span class="font-normal">(in €/km)</span></th>
                {#if $userRole === "client"}
                    <th class="border p-2 text-center">Actions</th>
                {/if}
            </tr>
        </thead>
        <tbody>
            {#each sortedData as row, index}
                <tr class={index % 2 === 0 ? 'bg-gray-50' : 'bg-white'}>
                    
                    <!-- Column 1 -->
                    <td class="border p-2 text-center">{formatDate(row.limit_date)}</td>

                    <!-- Column 2 -->
                    <td class="border p-2 text-center">{row.resource_type}</td>
                    
                    <!-- Column 3 -->
                    <td class="border p-2 text-center">{row.volume}</td>
                    
                    <!-- Column 4 -->
                    <td class="border p-2 text-center">{row.max_price_by_km.toFixed(2)}</td>
                    
                    <!-- Column 5 -->
                    <td class="border p-2 text-center">{row.current_price.toFixed(2)}</td>
                    
                    <!-- Column 6 -->
                    {#if $userRole === "client"}
                        <td class="border p-2 text-center">
                            <div class="flex flex-wrap justify-center space-x-2 space-y-2">
                                <button class="bg-blue-200 text-blue-800 px-4 py-2 flex items-center font-bold hover:bg-blue-300 transition-colors rounded-md"
                                    on:click={() => openModal(row.current_price, row.offer_id, row.max_price_by_km)}>
                                    <i class="fas fa-coins mr-2"></i>
                                    Bid
                                </button>
                            </div>
                        </td>
                    {/if}

                </tr>
            {/each}
        </tbody>
    </table>
</main>

<!-- svelte-ignore a11y-no-static-element-interactions -->
<!-- svelte-ignore a11y-no-static-element-interactions -->
<!-- svelte-ignore a11y-click-events-have-key-events -->
<!-- svelte-ignore a11y-label-has-associated-control -->

{#if isModalOpen}

    <div class="fixed inset-0 bg-black bg-opacity-50 flex justify-center items-center z-50" on:click={closeModal}>

        <div class="bg-white text-center p-6 rounded-lg shadow-lg w-1/4" on:click|stopPropagation>

            <!-- Close Button -->
            <button class="absolute top-2 right-2 text-gray-500 hover:text-gray-800" on:click={closeModal}>
                &times;
            </button>

            <!-- Modal Title -->
            <h2 class="text-2xl font-bold mb-6">Bid on the lot</h2>

            <!-- Form -->
            <form on:submit|preventDefault={bid}>

                <!-- Price -->
                <div class="mb-4">
                    <label class="block text-gray-700 text-lg font-bold">Price</label>
                    <p class="text-3xl font-bold text-gray-700">{priceValue} <span class="font-normal">€/km</span></p>
                    <input
                        type="range"
                        min={minPriceValue}
                        max={maxPriceValue}
                        step="0.1"
                        bind:value={priceValue}
                        class="range w-full mt-2"
                    />
                </div>


                <!-- Validate button -->
                <div class="flex justify-center mt-8">
                    <button type="submit" class="bg-blue-500 text-white px-6 py-2 rounded hover:bg-blue-600">
                        <i class="fas fa-check"></i>
                        <span class="font-bold">Validate</span>
                    </button>
                </div>
            </form>
        </div>
    </div>
{/if}


<!-- Style -->
<style>
    .range {
        -webkit-appearance: none;
        appearance: none;
        width: 100%;
        height: 15px;
        border-radius: 5px;
        background: #e2e8f0;
        outline: none;
        transition: background 0.3s;
    }

    .range:hover {
        background: #cbd5e1;
    }

    .range::-webkit-slider-thumb {
        -webkit-appearance: none;
        appearance: none;
        width: 20px;
        height: 20px;
        border-radius: 50%;
        background: #374151;
        cursor: pointer;
    }

    .range::-moz-range-thumb {
        width: 20px;
        height: 20px;
        border-radius: 50%;
        background: #374151;
        cursor: pointer;
    }

    .range::-ms-thumb {
        width: 20px;
        height: 20px;
        border-radius: 50%;
        background: #374151;
        cursor: pointer;
    }

</style>
