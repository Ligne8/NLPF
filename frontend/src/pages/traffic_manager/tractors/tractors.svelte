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
    let routes: Route[]  = [];
    let routesLoaded: boolean = false;
    let isStockExchangeModalOpen = false;
    let selectedLotId: number = null;
    let selectedStatus: string = 'all';
    let sortOption: string = 'none';
    let limitDate: string = '';
    let minDate: string = '';

    interface Route {
        name: string;
        id: string;
        route_path: string;
    }

    // Open select date modal
    const openStockExchangeModal = (lotId: number)=>{
        selectedLotId = lotId;  
        isStockExchangeModalOpen = true; 
    }

    // Close select date modal
    const closeStockExchangeModal = ()=>{
        fetchRoutes();
        fetchTableInfo();
        fetchLimitDate();
        isStockExchangeModalOpen = false; 
    }

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
            case 'return_from_market':
                return {color: 'bg-fuchsia-200 text-fuchsia-800', text: '◉ Return from market'};
            default:
                return { color: 'bg-gray-200 text-gray-800', text: '🛇 Unknown' };
        }
    }

    // Fetch limit date from backend
    async function fetchLimitDate() {
        try {
            const response = await axios.get(`${API_BASE_URL}/simulations/date`);
            const date = new Date(response.data.simulation_date);
            date.setDate(date.getDate() + 1);
            minDate = date.toISOString().split('T')[0];
        } catch (err) {
            console.error('Error fetching limit date:', err);
        }
    }

    onMount(() => {
        fetchRoutes();
        fetchTableInfo();
        fetchLimitDate();
    });


    const fetchRoutes = async () => {
        const response = await fetch(`${API_BASE_URL}/routes/traffic_manager/parsed/${$userId}`);
        if (response.ok) {
            const data = await response.json();
            if (data != null){
                routes = data;
                routesLoaded = true;
            }
        } else {
            console.error('Failed to fetch routes:', response.status);
        }
    }

    async function fetchTableInfo() {
        if($userRole !== "traffic_manager") {
            return;
        }
        await axios.get(`${API_BASE_URL}/tractors/trafficManager/${$userId}`)
            .then((response) => {
                tractors = response.data.map(tractor => ({
                    ...tractor,
                    selected_route: '',
                }));
            }).catch((error) => {
                console.error('Error fetching tractors:', error.response);
            });
    }

    const addRoute = async (t: Tractor)=>{
      axios.post(`${API_BASE_URL}/tractors/route`, {tractor_id: t.id ,route_id: t.selected_route.id})
        .then((response) => {
          fetchTableInfo()
        }).catch((error) => {
          console.error('Error adding route:', error.response);
        });
    }
    const removeRoute = async (t: Tractor)=>{
      axios.delete(`${API_BASE_URL}/tractors/route`, {data: {tractor_id: t.id}})
        .then((response) => {
          fetchTableInfo()
        }).catch((error) => {
          console.error('Error removing route:', error.response);
        });
    }

    const startTractor = async (t: Tractor)=>{
      if (t.route_id == null){
        alert("Please select a route first")
        return
      }
      axios.patch(`${API_BASE_URL}/tractors/updateState`, {id: t.id, state: "in_transit"})
        .then((response) => {
          fetchTableInfo()
        }).catch((error) => {
          console.error('Error starting tractor:', error.response);
        });
    }

    const stopTractor = async (t: Tractor)=>{
      axios.patch(`${API_BASE_URL}/tractors/updateState`, {id: t.id, state: "pending"})
        .then((response) => {
          fetchTableInfo()
        }).catch((error) => {
          console.error('Error starting tractor:', error.response);
        });
    }

    // Function to filter routes by current checkpoint
    function getMatchingRoutes(currentCheckpointName: string, endCheckointName: string): Route[] {
        return routes.filter(route => {
            return route.route_path.split(' ')[0].toLocaleLowerCase() === currentCheckpointName.toLocaleLowerCase();
        }).filter(route =>{
          return route.route_path.split(' ')[route.route_path.split(' ').length-1].toLocaleLowerCase() === endCheckointName.toLocaleLowerCase();
        });
    }

    // Assign tractor to trader
    async function assignTractorToTrader(tractorId: string) {
        if (limitDate == '') {
            alert('Veuillez sélectionner une date limite.');
            return;
        }
        const offerData = {
            limit_date: new Date(limitDate).toISOString(),
        };
        await axios.post(`${API_BASE_URL}/tractors/assign/${tractorId}/trader`, offerData)
        .then((response) => {
            fetchTableInfo();
            limitDate = ''; 
            closeStockExchangeModal();
        }).catch((error) => {
            console.error('Error assigning tractor to trader:', error.response);
        });
        }

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

</script>


<!-- Navbar -->
<Navbar/>
<TrafficManagerNavbar/>

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
                <option value="at_trader">At trader</option>
                <option value="archive">Archived</option>
            </select>

            <!-- Sort by name, volume and location -->
            <select bind:value={sortOption} class="border border-gray-300 rounded px-2 py-1">
                <option value="none" disabled selected>Sort by</option>
                <option value="name_asc">Name (A-Z)</option>
                <option value="name_desc">Name (Z-A)</option>
                <option value="loading_asc">Loading (Ascending)</option>
                <option value="loading_desc">Loading (Descending)</option>
                <option value="remaining_volume_asc">Remaining volume (Ascending)</option>
                <option value="remaining_volume_desc">Remaining volume (Descending)</option>
                <option value="location_asc">Location (A-Z)</option>
                <option value="location_desc">Location (Z-A)</option>
            </select>

        </div>

        <div class="flex justify-between items-center self-end">

            <!-- Reload button -->
            <button class="bg-gray-800 text-white font-bold px-4 py-2 rounded flex items-center hover:bg-gray-900 transition-colors self-end"
                    on:click={fetchTableInfo}
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
                <th class="border p-2 text-center max-w-16">Name</th>
                <th class="border p-2 text-center max-w-16">Status</th>
                <th class="border p-2 text-center max-w-16">Loading <span class="font-normal">(in m³)</span></th>
                <th class="border p-2 text-center max-w-16">Location</th>
                <th class="border p-2 text-center max-w-16">Departure / Arrival</th>
                <th class="border p-2 text-center">Route</th>
                <th class="border p-2 text-center w-60">Actions</th>
            </tr>
            </thead>
            <tbody>
            {#each sortedData as row, index}
                <tr class={index % 2 === 0 ? 'bg-gray-50' : 'bg-white'}>

                    <!-- Column 1 -->
                    <td class="border p-2 text-center">{row.name}</td>

                    <!-- Column 2 -->
                    <td class="border p-2 text-center max-w-32">
                            <span class={`px-2 py-1 rounded ${getStateInfo(row.state).color}`}>
                                {getStateInfo(row.state).text}
                            </span>
                    </td>

                    <!-- Column 3 -->
                    <td class="border p-2 text-center">{row.current_units}/{row.max_units}</td>

                    <!-- Column 4 -->
                    <td class="border p-2 text-center">{row.current_checkpoint.name}</td>

                    <!-- Column 5 -->
                    <td class="border p-2 text-center">{row.start_checkpoint.name} / {row.end_checkpoint.name}</td>

                    <!-- Column 6 -->
                    <td class="border p-2 text-center max-w-28">
                        {#if row.state === 'pending' && row.route_id == null}
                            {#if routesLoaded && getMatchingRoutes(row.current_checkpoint.name, row.end_checkpoint.name).length > 0}
                                <select bind:value={row.selected_route}
                                        class="border border-gray-300 rounded px-2 py-1 mx-auto w-4/5"
                                        on:change={() => addRoute(row)}>
                                    <option value="" disabled selected>Select a route</option>
                                    {#each getMatchingRoutes(row.current_checkpoint.name, row.end_checkpoint.name) as routeOption}
                                        <option value={routeOption}>{routeOption.route_path}</option>
                                    {/each}
                                </select>
                            {:else if routesLoaded}
                                <span class="px-2 py-1 mx-auto w-4/5 block text-gray-500">None</span>
                            {/if}
                        {:else}
                            {#if row.route}
                                <span class="px-2 py-1 mx-auto w-4/5 block">
                                    {row.route.name}
                                </span>
                            {:else}
                                <span class="px-2 py-1 mx-auto w-4/5 block text-gray-500">None</span>
                            {/if}
                        {/if}
                    </td>

                    <!-- Column 6 -->
                    <td class="border p-2 text-center">
                        {#if row.state === 'in_transit'}
                            <div class="flex flex-wrap justify-center space-x-2 space-y-2">
                                <button on:click={()=>(stopTractor(row))} class="bg-red-200 text-red-600 px-4 py-2 flex items-center font-bold hover:bg-red-300 transition-colors rounded-md">
                                    <i class="fas fa-hand mr-2"></i>
                                    Stop
                                </button>
                            </div>
                        {:else if row.state === 'pending'}
                            <div class="flex flex-wrap justify-center space-x-2 space-y-2">
                                <button on:click={()=>{startTractor(row)}} class="bg-green-200 text-green-800 px-4 py-2 flex items-center font-bold hover:bg-green-300 transition-colors rounded-md">
                                    <i class="fas fa-truck mr-2"></i>
                                    Start
                                </button>
                                <button class="bg-blue-200 text-blue-800 px-4 py-2 flex items-center font-bold hover:bg-blue-300 transition-colors rounded-md"
                                    on:click={openStockExchangeModal(row.id)} >
                                    <i class="fas fa-plus mr-2"></i>
                                    Stock exchange
                                </button>
                                {#if row.route_id !== null}
                                    <button on:click={()=>{removeRoute(row)}} class="bg-red-200 text-red-600 px-4 py-2 flex items-center font-bold hover:bg-red-300 transition-colors rounded-md">
                                        <i class="fas fa-eraser mr-2"></i>
                                        Remove route
                                    </button>
                                {/if}
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
{#if isStockExchangeModalOpen}
    <div class="fixed inset-0 bg-black bg-opacity-50 flex justify-center items-center z-50"
         on:click={closeStockExchangeModal}>
        <div class="bg-white p-6 rounded-lg shadow-lg w-1/3" on:click|stopPropagation>
            <!-- Close Button -->
            <button class="absolute top-2 right-2 text-gray-500 hover:text-gray-800" on:click={closeStockExchangeModal}>
                &times;
            </button>

            <!-- Modal Title -->
            <h2 class="text-2xl font-bold mb-6">Stock Exchange</h2>

            <!-- Form -->
            <form on:submit|preventDefault={assignTractorToTrader(selectedLotId)}>

                <!-- Limit Date -->
                <div class="mb-2">
                    <label class="block text-gray-700 text-sm font-bold mb-2">Limit Date :</label>
                    <input type="date"
                        class="w-full border border-gray-300 p-2 rounded"
                        bind:value={limitDate}
                        min={minDate}
                        required
                    />
                </div>

                <!-- Submit button -->
                <div class="flex justify-center mt-4">
                    <button type="submit" class="bg-blue-500 text-white px-6 py-2 rounded hover:bg-blue-600">
                        <i class="fas fa-check"></i>
                        <span class="font-bold">Submit</span>
                    </button>
                </div>
            </form>
        </div>
    </div>
{/if}