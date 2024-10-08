<script lang="ts">
    import Navbar from '@components/Navbar.svelte';
    import TrafficManagerNavbar from '@components/TrafficManagerNavbar.svelte';
    import { onMount } from 'svelte';
    import axios from 'axios';
    import {userId} from "@stores/store";
    const API_BASE_URL = import.meta.env.VITE_API_BASE_URL;

    interface Checkpoints {
        id: string;
        name: string;
    }

    // Variables
    let title: string = 'Route management';
    let subtitle: string = 'Manage available routes and itineraries.';
    let checkpoints: Checkpoints[] = [];
    let selectedCheckpoints: Checkpoints[] = [];
    let newRouteName: string = '';
    let tableData: Route[] = [];

    interface Route {
        name: string;
        id: string;
        route_path: string;
    }


    onMount(() => {
        fetchCheckpoints();
        fetchRoutes();
    });

    const fetchRoutes = async () => {
        const response = await fetch(`${API_BASE_URL}/routes/traffic_manager/parsed/${$userId}`);
        if (response.ok) {
            const data = await response.json();
            if (data != null){
              tableData = data;
            }
        } else {
            console.error('Failed to fetch routes:', response.status);
        }
    }

    async function fetchCheckpoints(){
        try {
            const response = await fetch(`${API_BASE_URL}/checkpoints`);
            if (response.ok)
            {
                const data = await response.json();
                checkpoints = data.map((checkpoint: any) => ({name: checkpoint.name, id: checkpoint.id}));
                selectedCheckpoints = [checkpoints[0]];
            }
            else
            {
                console.error('Failed to fetch checkpoints:', response.status);
            }
        } catch (error) {
            console.error('Error fetching checkpoints:', error);
        }
    }


    // Function to add a new checkpoint select
    function addCheckpoint() {
        const availableCheckpoints = getAvailableCheckpoints(selectedCheckpoints.length);
        if (availableCheckpoints.length > 0)
        {
            selectedCheckpoints = [...selectedCheckpoints, availableCheckpoints[0]];
        }
    }

    // Function to remove a checkpoint
    function removeCheckpoint(index: number) {
        if (index < 1)
            return;
        selectedCheckpoints = selectedCheckpoints.slice(0, index).concat(selectedCheckpoints.slice(index + 1));
        selectedCheckpoints = selectedCheckpoints.filter((value, i, arr) => i === 0 || value !== arr[i - 1]);
    }

    function getAvailableCheckpoints(currentIndex: number): Checkpoints[] {
        if (currentIndex === 0)
            return checkpoints;

        const previousCheckpoint = selectedCheckpoints[currentIndex - 1];
        if (currentIndex === checkpoints.length - 1)
        {
            return checkpoints.filter(cp => cp !== previousCheckpoint);
        }
        else
        {
            const nextCheckpoint = selectedCheckpoints[currentIndex + 1];
            return checkpoints.filter(cp => cp !== previousCheckpoint && cp !== nextCheckpoint);
        }
    }

    // Function to add a new route to the table
    function addRouteToTable() {
        const newRoute = {
            name: newRouteName,
            traffic_manager_id: $userId,
            route: selectedCheckpoints.map((cp, index) => 
            {
              const payload = {
                checkpoint_id: cp.id,
                position: index+1
              }
              return payload;
            }
          )
        };
        axios.post(`${API_BASE_URL}/routes`, newRoute)
            .then(() => {
                fetchRoutes();
            })
            .catch((error:any) => {
                console.error('Failed to add route:', error);
            });
    }

    // Function to validate the route
    function validateRoute() {
        const validCheckpoints = selectedCheckpoints.filter(cp => cp.name !== '');
        if (validCheckpoints.length < 2 || newRouteName.trim() === '')
            return;

        // Add the new route to the table
        addRouteToTable();

        // Reset the inputs after validation
        selectedCheckpoints = [checkpoints[0]];
        newRouteName = '';
    }

</script>

<!-- Navbar -->
<Navbar />
<TrafficManagerNavbar />

<main class="p-10 mt-40">

    <!-- Title and subtitle -->
    <section class="mb-8">
        <h1 class="text-4xl font-bold mb-2">{title}</h1>
        <h2 class="text-2xl text-gray-600">{subtitle}</h2>
    </section>

    <div class="flex">

        <!-- Left part -->
        <div class="w-2/3 pr-8 border-r border-gray-300">
            <h2 class="text-2xl text-gray-800 font-bold mb-4">
                <i class="fas fa-list mr-2"></i>
                List of routes
            </h2>
            <table class="w-full border-collapse border border-gray-300">
                <thead>
                <tr class="bg-gray-100">
                    <th class="border p-2 text-center">Name</th>
                    <th class="border p-2 text-center">Steps</th>
                </tr>
                </thead>
                <tbody>
                {#each tableData as row, index}
                    <tr class={index % 2 === 0 ? 'bg-gray-50' : 'bg-white'}>

                        <!-- Column 1 -->
                        <td class="border p-2 text-center">{row.name}</td>

                        <!-- Column 2 -->
                        <td class="border p-2 text-center">
                            {row.route_path}
                        </td>

                    </tr>
                {/each}
                </tbody>
            </table>
        </div>

        <!-- Right part -->
        <div class="w-1/3 pl-8">
            <h2 class="text-2xl text-gray-800 font-bold mb-4">
                <i class="fas fa-plus mr-2"></i>
                Add a route
            </h2>

            <!-- Route name input field -->
            <div class="mb-4">
                <input
                        type="text"
                        id="route-name"
                        class="border border-gray-300 rounded px-3 py-2 w-full"
                        bind:value={newRouteName}
                        placeholder="Enter route name"
                />
            </div>

            <!-- Checkpoints select inputs -->
            <div class="mb-4">
                {#each selectedCheckpoints as selected, index}
                    <div class="mb-1 flex items-center">
                        {#if index !== 0}
                            <button
                                    on:click={() => removeCheckpoint(index)}
                                    class="bg-red-500 text-white rounded-md w-8 h-8 hover:bg-red-600 flex items-center justify-center mr-2"
                                    title="Supprimer ce checkpoint"
                            >
                                <i class="fas fa-minus"></i>
                            </button>
                        {/if}

                        <select id="checkpoint-{index}" class="border border-gray-300 rounded px-3 py-2 w-full"
                                bind:value={selectedCheckpoints[index]}
                        >
                            {#each getAvailableCheckpoints(index) as checkpoint}
                                <option value={checkpoint}>
                                    {checkpoint.name}
                                </option>
                            {/each}
                        </select>
                    </div>
                {/each}
                <!-- Add checkpoint button -->
                <button
                        on:click={addCheckpoint}
                        class="bg-gray-800 text-white rounded px-4 py-2 w-full hover:bg-gray-900 transition-colors flex items-center justify-center"
                        disabled={selectedCheckpoints[selectedCheckpoints.length - 1] === null}
                        class:bg-gray-300={selectedCheckpoints[selectedCheckpoints.length - 1] === null}
                >
                    <i class="fas fa-plus"></i>
                </button>
            </div>

            <!-- Validate button -->
            {#if selectedCheckpoints.filter(cp => cp !== null).length >= 2 && newRouteName.trim() !== ''}
                <div class="flex justify-center mt-4">
                    <button
                            on:click={validateRoute}
                            class="bg-blue-500 text-white font-bold rounded px-6 py-3 hover:bg-blue-600 transition-colors"
                    >
                        <i class="fas fa-check mr-2"></i>
                        Validate
                    </button>
                </div>
            {/if}
        </div>
    </div>
</main>
