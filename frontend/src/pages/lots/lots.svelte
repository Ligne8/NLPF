<script lang="ts">
    import { onMount } from 'svelte';
    import Navbar from '@components/Navbar.svelte';
    import {userId} from "@stores/store";

    const API_BASE_URL = import.meta.env.VITE_API_BASE_URL;

    interface Checkpoint {
      id: string;
      name: string;
    }

    let title: string = 'Lot management';
    let subtitle: string = 'Track the status of your lots in real time.';
    let isModalOpen = false;
    let checkpoints: Checkpoint[] = [];
    let types = ['Bulk', 'Solid', 'Liquid'];
    let selectedType: string = types[0];
    let volume: string = '';
    let maxPrice: string = '';
    let selectedDeparture: Checkpoint;
    let selectedArrival: Checkpoint;
    let tableData: LotTable[] = [];
    let trafficManagers: TrafficManager[] = [];
    let selectedStatus: string = 'all';
    let sortOption: string = 'none';
    
    interface TrafficManager {
      id: string;
      name: string;
    }

    interface LotTable {
      id: string
      state: string;
      volume: number;
      currentCheckpoint: string;
      startCheckpoint: string;
      endCheckpoint: string;
      trafficManager: TrafficManager;
    }

    const fetchAllData = async () => {
        await fetchLots();
        await fetchTrafficManagers();
        await fetchCheckpoints();
    };

    // Fetch all data
    onMount(async () => {
      await fetchAllData();
    });

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
            case 'archived':
                return { color: 'bg-gray-200 text-gray-800', text: '◉ Archived' };
            case 'at_trader':
                return { color: 'bg-purple-200 text-purple-800', text: '◉ At trader' };
            default:
                return { color: 'bg-gray-200 text-gray-800', text: '🛇 Unknown' };
        }
    }

    // Function to validate input
    function validateVolume(event: Event) {
        const input = event.target as HTMLInputElement;
        input.value = input.value.replace(/[^0-9.]/g, '');
        if ((input.value.match(/\./g) || []).length > 1)
            input.value = input.value.replace(/\.+$/, '');
        volume = input.value;
    }

    // Function to validate input
    function validateMaxPrice(event: Event) {
        const input = event.target as HTMLInputElement;
        input.value = input.value.replace(/[^0-9.]/g, '');
        if ((input.value.match(/\./g) || []).length > 1)
            input.value = input.value.replace(/\.+$/, '');
        maxPrice = input.value;
    }

    async function fetchTrafficManagers(){
      try {
        const response = await fetch(`${API_BASE_URL}/users/traffic_managers`);
        if (response.ok) {
          const data = await response.json();
          trafficManagers = data.map((trafficManager: any) => ({id: trafficManager.id, name: `${trafficManager.username}`}));
        } else {
          console.error('Failed to fetch traffic managers:', response.status);
        }
      } catch (error) {
        console.error('Error fetching traffic managers:', error);
      }
    }

    async function fetchLots() {
      try {
        const response = await fetch(`${API_BASE_URL}/lots/owner/${$userId}`);
        if (response.ok) {
          const data = await response.json();
          tableData = data.map((lot: any) => ({
            id: lot.id,
            state: lot.state,
            volume: lot.volume,
            currentCheckpoint: lot.current_checkpoint.name,
            startCheckpoint: lot.start_checkpoint.name,
            endCheckpoint: lot.end_checkpoint.name,
            trafficManager: lot.traffic_manager == null ? null : lot.traffic_manager.username,
            createdAt: new Date(lot.created_at)
          })).sort((a:any, b:any) => b.createdAt - a.createdAt);
        } else {
          console.error('Failed to fetch lots:', response.status);
        }
      } catch (error) {
        console.error('Error fetching lots:', error);
      }
    }

    async function fetchCheckpoints(){
        try {
            const response = await fetch(`${API_BASE_URL}/checkpoints`);
            if (response.ok)
            {
                const data = await response.json();
                checkpoints = data.map((checkpoint: any) => ({name: checkpoint.name, id: checkpoint.id}));
                selectedDeparture = checkpoints[0];
                selectedArrival = checkpoints.length > 1 ? checkpoints[1] : checkpoints[0];
            }
            else
            {
                console.error('Failed to fetch checkpoints:', response.status);
            }
        } catch (error) {
            console.error('Error fetching checkpoints:', error);
        }
    }

    // Function to add a lot
    function addLot() {

        // Add a lot to the table
        const newLot = {
            resource_type: selectedType.toLowerCase(),
            volume: parseFloat(volume),
            max_price_by_km: parseFloat(maxPrice),
            current_checkpoint_id: selectedDeparture.id,
            start_checkpoint_id: selectedDeparture.id,
            end_checkpoint_id: selectedArrival.id,
            state: 'available',
            owner_id: $userId
        };

        selectedType = types[0];
        volume = '';
        maxPrice = '';
        selectedDeparture = checkpoints[0]; // Valeur par défaut
        selectedArrival = checkpoints[1]; // Valeur par défaut


        fetch(`${API_BASE_URL}/lots`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(newLot)
        }).then(response => {
          fetchLots();
        }).catch(error => {
            console.error('Error adding lot:', error);
        });
        closeModal();
    }

    // Function to open modal
    function openModal() {
        isModalOpen = true;
    }

    // Function to close modal
    function closeModal() {
        isModalOpen = false;
    }

    const assignToTrafficManager = (lotId: string, trafficManager: TrafficManager) => {
      if (trafficManager == null) {
        console.error('Traffic manager is null');
        alert('Veuillez sélectionner un traffic manager');
        return;
      }
      fetch(`${API_BASE_URL}/lots/traffic_manager`, {
          method: 'POST',
          headers: {
              'Content-Type': 'application/json'
          },
          body: JSON.stringify({ traffic_manager_id: trafficManager.id, lot_id: lotId })
      }).then(response => {
          fetchLots();
          alert('Lot attribué avec succès');
      }).catch(error => {
          console.error('Error assigning lot to traffic manager:', error);
          alert('Erreur lors de l\'attribution du lot');
      });
    }

    // Fucntion to delete a lot
    const deleteLot = (lotId: string) => {
      fetch(`${API_BASE_URL}/lots/${lotId}`, {
          method: 'DELETE'
      }).then(response => {
          fetchLots();
          alert('Lot supprimé avec succès');
      }).catch(error => {
          console.error('Error deleting lot:', error);
          alert('Erreur lors de la suppression du lot');
      });
    }

    // Update data depending on filters
    $: sortedData = (() => {
        let data = selectedStatus === 'all' ? tableData : tableData.filter(lot => lot.state === selectedStatus);

        switch (sortOption) {
            case 'volume_asc':
                return data.sort((a, b) => a.volume - b.volume);
            case 'volume_desc':
                return data.sort((a, b) => b.volume - a.volume);
            case 'location_asc':
                return data.sort((a, b) => a.currentCheckpoint.localeCompare(b.currentCheckpoint));
            case 'location_desc':
                return data.sort((a, b) => b.currentCheckpoint.localeCompare(a.currentCheckpoint));
            default:
                return data;
        }
    })();

</script>


<!-- Navbar -->
<Navbar/>

<main class="p-10 pt-40">

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

        <div class="flex justify-between items-center self-end">

            <!-- Create button -->
            <button class="bg-blue-500 text-white mr-2 font-bold px-4 py-2 rounded flex items-center hover:bg-blue-600 transition-colors self-end"
                    on:click={openModal}
            >
                <i class="fas fa-plus mr-2"></i>
                Add a lot
            </button>

            <!-- Reload button -->
            <button class="bg-gray-800 text-white font-bold px-4 py-2 rounded flex items-center hover:bg-gray-900 transition-colors self-end"
                    on:click={fetchAllData}
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
                <th class="border p-2 text-center">Volume <span class="font-normal">(in m³)</span></th>
                <th class="border p-2 text-center">Location</th>
                <th class="border p-2 text-center">Departure / Arrival</th>
                <th class="border p-2 text-center">Traffic manager</th>
                <th class="border p-2 text-center">Actions</th>
            </tr>
            </thead>
            <tbody>
            {#each sortedData as row, index}
                <tr class={index % 2 === 0 ? 'bg-gray-50' : 'bg-white'}>


                    <!-- Column 2 -->
                    <td class="border p-2 text-center">
                            <span class={`px-2 py-1 rounded ${getStatusInfo(row.state).color}`}>
                                {getStatusInfo(row.state).text}
                            </span>
                    </td>

                    <!-- Column 3 -->
                    <td class="border p-2 text-center">{row.volume}</td>

                    <!-- Column 4 -->
                    <td class="border p-2 text-center">{row.currentCheckpoint}</td>

                    <!-- Column 5 -->
                    <td class="border p-2 text-center">
                        {row.startCheckpoint} / {row.endCheckpoint}
                    </td>

                    <!-- Column 6 -->
                    <td class="border p-2 text-center">
                        {#if row.state === 'available'}
                            <select bind:value={row.trafficManager}  class="border border-gray-300 rounded px-2 py-1 mx-auto w-4/5">
                                {#each trafficManagers as trafficManagerOption}
                                    <option value={trafficManagerOption}>{trafficManagerOption.name}</option>
                                {/each}
                            </select>
                        {:else}
                            <span class="text-gray-500">{row.trafficManager}</span>
                        {/if}
                    </td>

                    <td class="border p-2 text-center">
                        {#if row.state === 'available'}
                            <div class="flex flex-wrap justify-center space-x-2">
                                {#if row.trafficManager}
                                    <button class="bg-green-200 text-green-800 px-4 py-2 flex items-center font-bold hover:bg-green-300 transition-colors rounded-md"
                                            on:click={()=>{assignToTrafficManager(row.id, row.trafficManager)}}
                                    >
                                        <i class="fas fa-truck mr-2"></i>
                                        Assign
                                    </button>
                                {/if}
                                <button class="bg-blue-200 text-blue-800 px-4 py-2 flex items-center font-bold hover:bg-blue-300 transition-colors rounded-md">
                                    <i class="fas fa-plus mr-2"></i>
                                    Stock exchange
                                </button>
                                <button on:click={()=>{deleteLot(row.id)}} class="bg-gray-800 text-white px-4 py-2 flex items-center font-bold hover:bg-black transition-colors rounded-md">
                                    <i class="fas fa-right-from-bracket mr-2"></i>
                                    Retirer
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

        <div class="bg-white p-6 rounded-lg shadow-lg w-1/3" on:click|stopPropagation>

            <!-- Close Button -->
            <button class="absolute top-2 right-2 text-gray-500 hover:text-gray-800" on:click={closeModal}>
                &times;
            </button>

            <!-- Modal Title -->
            <h2 class="text-2xl font-bold mb-6">Add a lot</h2>

            <!-- Form -->
            <form on:submit|preventDefault={addLot}>

                <!-- Type -->
                <div class="mb-2">
                    <label class="block text-gray-700 text-sm font-bold mb-2">Type :</label>
                    <select class="w-full border border-gray-300 p-2 rounded" bind:value={selectedType}>
                        {#each types as type}
                            <option value={type}>{type}</option>
                        {/each}
                    </select>
                </div>

                <!-- Volume -->
                <div class="mb-2">
                    <label class="block text-gray-700 text-sm font-bold mb-2">Volume :</label>
                    <input type="text"
                           class="w-full border border-gray-300 p-2 rounded"
                           placeholder="Enter volume (in m³)"
                           on:input={validateVolume}
                           value={volume}
                           required
                    >
                </div>

                <!-- Max price -->
                <div class="mb-2">
                    <label class="block text-gray-700 text-sm font-bold mb-2">Maximum price :</label>
                    <input type="text"
                           class="w-full border border-gray-300 p-2 rounded"
                           placeholder="Enter maximum price (per km)"
                           on:input={validateMaxPrice}
                           value={maxPrice}
                           required
                    >
                </div>

                <!-- Departure -->
                <div class="mb-2">
                    <label class="block text-gray-700 text-sm font-bold mb-2">Departure :</label>
                    <select class="w-full border border-gray-300 p-2 rounded" bind:value={selectedDeparture}>
                        {#each checkpoints as checkpoint}
                            {#if checkpoint.id !== selectedArrival.id}
                                <option value={checkpoint}>{checkpoint.name}</option>
                            {/if}
                        {/each}
                    </select>
                </div>

                <!-- Arrival -->
                <div class="mb-2">
                    <label class="block text-gray-700 text-sm font-bold mb-2">Arrival :</label>
                    <select class="w-full border border-gray-300 p-2 rounded" bind:value={selectedArrival}>
                        {#each checkpoints as checkpoint}
                            {#if checkpoint.id !== selectedDeparture.id}
                                <option value={checkpoint}>{checkpoint.name}</option>
                            {/if}
                        {/each}
                    </select>
                </div>

                <!-- Add button -->
                <div class="flex justify-center mt-4">
                    <button type="submit" class="bg-blue-500 text-white px-6 py-2 rounded hover:bg-blue-600">
                        <i class="fas fa-plus"></i>
                        <span class="font-bold">Add</span>
                    </button>
                </div>
            </form>
        </div>
    </div>
{/if}