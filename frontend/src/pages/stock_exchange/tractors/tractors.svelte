<script lang="ts">
    import Navbar from '@components/Navbar.svelte';
    import StockExchangeNavbar from '@components/StockExchangeNavbar.svelte';
    import { userRole, userId } from '@stores/store';
    import axios from 'axios';
    import type { Tractor } from 'src/interface/tractorInterface';
    import { onMount } from 'svelte';

    // Variables
    let title: string = 'Tractor market';
    let subtitle: string = 'Negotiate available freight space according to market needs.';
    const API_BASE_URL = import.meta.env.VITE_API_BASE_URL;
    let tractors: Tractor[] = [];
    let isModalOpen = false;
    let priceValue: number = 1.0;
    let minPriceValue: number = 1.0;
    let maxPriceValue: number = 100.0;
    let volumeValue: number = 1.0;
    let minVolumeValue: number = 1.0;
    let maxVolumeValue: number = 50.0;
    let sortOption: string = 'none';
    let current_offer_id: number;

    // Function to format timestamp into DD/MM/YYYY
    const formatDate = (timestamp: number) => {
        const date = new Date(timestamp);
        return `${date.getDate().toString().padStart(2, '0')}/${(date.getMonth() + 1).toString().padStart(2, '0')}/${date.getFullYear()}`;
    };

    // Function to increase volume
    function increaseVolume() {
    if (volumeValue < maxVolumeValue)
        volumeValue += 1;
}
    // Function to decrease volume
    function decreaseVolume() {
        if (volumeValue > minVolumeValue)
            volumeValue -= 1;
    }

    // Function to open tractors modal
    function openModal(currentPrice: any, offer_id: any, min_price_by_km: number, bid_owner_id: string) {
      console.log(bid_owner_id)
        current_offer_id = offer_id;
        priceValue = min_price_by_km;
        minPriceValue = min_price_by_km;
        isModalOpen = true;
    }

    // Function to close tractors modal
    function closeModal() {
        current_offer_id = 0;
        isModalOpen = false;
    }

    // Function to bid
    function bid() {
        console.log('Bid:', priceValue);
        console.log('VolumeValue:', volumeValue);
        console.log('offerId:', current_offer_id);

        const payload = {
            offer_id: current_offer_id,
            bid: priceValue,
            volume: volumeValue,
            owner_id: $userId
        }
        axios.post(`${API_BASE_URL}/stock_exchange/tractor/bid`, payload).then((response) => {
          fetchTractors();
          closeModal();
        }).catch((error) => {
            console.error('Error bidding:', error.response);
        });
    }

    // Fetch table info
    async function fetchTractors() {
        await axios.get(`${API_BASE_URL}/stock_exchange/tractor_offers`)
            .then((response) => {
                tractors = response.data;
            }).catch((error) => {
                console.error('Error fetching tractors:', error.response);
            });
    }

    // Fetch all data
    onMount(() => {
        fetchTractors();
    });

    // Update data depending on filters
    $: sortedData = (() => {
        let data = tractors;

        switch (sortOption) {
            case 'loading_asc':
                return data.sort((a, b) => (a.current_units / a.max_units) - (b.current_units / b.max_units));
            case 'loading_desc':
                return data.sort((a, b) => (b.current_units / b.max_units) - (a.current_units / a.max_units));
            case 'remaining_volume_asc':
                return data.sort((a, b) => (a.max_units - a.current_units) - (b.max_units - b.current_units));
            case 'remaining_volume_desc':
                return data.sort((a, b) => (b.max_units - b.current_units) - (a.max_units - a.current_units));
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

            <!-- Sort by name, volume and location -->
            <select bind:value={sortOption} class="border border-gray-300 rounded px-2 py-1">
                <option value="none" disabled selected>Sort by</option>
                <option value="loading_asc">Loading (Ascending)</option>
                <option value="loading_desc">Loading (Descending)</option>
                <option value="remaining_volume_asc">Remaining volume (Ascending)</option>
                <option value="remaining_volume_desc">Remaining volume (Descending)</option>
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

    <table class="table-auto w-full border-collapse border border-gray-300">
        <thead>
            <tr class="bg-gray-100">
                <th class="border p-2 text-center">Expiration date</th>
                <th class="border p-2 text-center">Type</th>
                <th class="border p-2 text-center">Loading<br><span class="font-normal">(in m³)</span></th>
                <th class="border p-2 text-center">Minimum price<br><span class="font-normal">(in €/km)</span></th>
                <th class="border p-2 text-center">Current price<br><span class="font-normal">(in €/km)</span></th>
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
                    <td class="border p-2 text-center">{row.current_units}/{row.max_units}</td>
                    
                    <!-- Column 4 -->
                    <td class="border p-2 text-center">{row.min_price_by_km.toFixed(2)}</td>
                    
                    <!-- Column 5 -->
                    <td class="border p-2 text-center">{row.current_price.toFixed(2)}</td>
                    
                    <!-- Column 6 -->
                    {#if $userRole === "client"}
                        <td class="border p-2 text-center">
                            <div class="flex flex-wrap justify-center space-x-2 space-y-2">
                                <button class="bg-blue-200 text-blue-800 px-4 py-2 flex items-center font-bold hover:bg-blue-300 transition-colors rounded-md"
                                    on:click={() => openModal(row.current_price, row.offer_id, row.min_price_by_km, row.id)}
                                >
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
            <h2 class="text-2xl font-bold mb-6">Bid on the tractor</h2>

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

                <!-- Volume -->
                <div class="mb-4">
                    <label class="block text-gray-700 text-lg font-bold">Volume <span class="font-normal">(in m³)</span></label>
                    <div class="flex items-center justify-between">
                        <button
                            type="button"
                            class="bg-gray-200 px-3 py-2 rounded disabled:opacity-50"
                            on:click|stopPropagation={decreaseVolume}
                            disabled={volumeValue === minVolumeValue}
                        >
                            <i class="fas fa-minus"></i>
                        </button>

                        <input
                            type="number"
                            min={minVolumeValue}
                            max={maxVolumeValue}
                            bind:value={volumeValue}
                            class="text-2xl font-bold mx-4 text-gray-700 w-16 text-center"
                            on:input|preventDefault={() => {
                                if (volumeValue < minVolumeValue) volumeValue = minVolumeValue;
                                if (volumeValue > maxVolumeValue) volumeValue = maxVolumeValue;
                            }}
                        />

                        <button
                            type="button"
                            class="bg-gray-200 px-3 py-2 rounded disabled:opacity-50"
                            on:click|stopPropagation={increaseVolume}
                            disabled={volumeValue === maxVolumeValue}
                        >
                            <i class="fas fa-plus"></i>
                        </button>
                    </div>
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

    input[type="number"] {
        -moz-appearance: textfield;
    }

    input[type="number"]::-webkit-outer-spin-button,
    input[type="number"]::-webkit-inner-spin-button {
        -webkit-appearance: none;
        margin: 0;
    }

</style>
