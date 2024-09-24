<script lang="ts">
    import Navbar from '../../../components/Navbar.svelte';
    import TrafficManagerNavbar from '../../../components/TrafficManagerNavbar.svelte';

    // Variables
    let title: string = 'Gestion des routes';
    let subtitle: string = 'Gérez les routes et les itinéraires disponibles.';
    let checkpoints = ['Paris', 'Lyon', 'Marseille', 'Toulouse', 'Nice', 'Nantes'];

    // Example data
    let tableData = [
        { name: 'Route 1', route: ['Paris', 'Montpellier', 'Marseille'] },
        { name: 'Route 2', route: ['Paris', 'Montpellier'] },
        { name: 'Route 3', route: ['Marseille', 'Lyon', 'Marseille'] },
        { name: 'Route 4', route: ['Montpellier', 'Paris', 'Lyon', 'Perpignan'] },
    ];

    // Store selected checkpoints
    let selectedCheckpoints: string[] = [checkpoints[0]]; // Default to the first checkpoint

    // Function to simulate fetching checkpoints from the backend
    function fetchCheckpoints() {
        return checkpoints;
    }

    // Function to add a new checkpoint select
    function addCheckpoint() {
        // Set the default value to the first checkpoint
        selectedCheckpoints = [...selectedCheckpoints, checkpoints[0]];
        checkpoints = fetchCheckpoints();
    }

    // Function to add a new route to the table
    function addRouteToTable() {
        const newRoute = {
            name: `Route ${tableData.length + 1}`, // Dynamically generate route name
            route: selectedCheckpoints.filter(cp => cp !== '') // Get selected checkpoints
        };

        tableData = [...tableData, newRoute]; // Update the table data
    }

    // Function to validate the route
    function validateRoute() {
        const validCheckpoints = selectedCheckpoints.filter(cp => cp !== '');
        if (validCheckpoints.length < 2)
            return;

        // Add the new route to the table
        addRouteToTable();

        // Reset the inputs after validation
        selectedCheckpoints = [checkpoints[0]]; // Reset to default after validation
    }
</script>


<!-- Navbar -->
<Navbar />
<TrafficManagerNavbar />


<main class="p-10">

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
                Liste des routes
            </h2>
            <table class="w-full border-collapse border border-gray-300">
                <thead>
                    <tr class="bg-gray-100">
                        <th class="border p-2 text-center">Nom</th>
                        <th class="border p-2 text-center">Étapes</th>
                    </tr>
                </thead>
                <tbody>
                    {#each tableData as row, index}
                        <tr class={index % 2 === 0 ? 'bg-gray-50' : 'bg-white'}>

                            <!-- Column 1 -->
                            <td class="border p-2 text-center">{row.name}</td>

                            <!-- Column 2 -->
                            <td class="border p-2 text-center">
                                {row.route.join(' - ')}
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
                Ajouter une route
            </h2>

            <!-- Checkpoints select inputs -->
            <div class="mb-4">
                {#each selectedCheckpoints as selected, index}
                    <div class="mb-1">
                        <select id="checkpoint-{index}" class="border border-gray-300 rounded px-3 py-2 w-full"
                                bind:value={selectedCheckpoints[index]}
                                disabled={index !== selectedCheckpoints.length - 1}
                        >
                            {#each checkpoints as checkpoint}
                                <option value={checkpoint}>
                                    {checkpoint}
                                </option>
                            {/each}
                        </select>
                    </div>
                {/each}

                <!-- Add checkpoint button -->
                <button 
                    on:click={addCheckpoint} 
                    class="bg-gray-800 text-white rounded px-4 py-2 w-full hover:bg-gray-900 transition-colors flex items-center justify-center"
                    disabled={selectedCheckpoints[selectedCheckpoints.length - 1] === ''}
                    class:bg-gray-300={selectedCheckpoints[selectedCheckpoints.length - 1] === ''}
                >
                    <i class="fas fa-plus"></i>
                </button>
            </div>

            <!-- Validate button -->
            {#if selectedCheckpoints.filter(cp => cp !== '').length >= 2}
                <div class="flex justify-center mt-4">
                    <button 
                        on:click={validateRoute}
                        class="bg-blue-500 text-white font-bold rounded px-6 py-3 hover:bg-blue-600 transition-colors"
                    >
                        <i class="fas fa-check mr-2"></i>
                        Valider la route
                    </button>
                </div>
            {/if}
        </div>
    </div>
</main>
