<script lang="ts">
    import { onMount } from 'svelte';
    import Navbar from '../../components/Navbar.svelte';

    const API_BASE_URL = import.meta.env.VITE_API_BASE_URL;

    // Variables
    let title: string = 'Gestion des Lots';
    let subtitle: string = 'Suivez l’état de vos lots en temps réel.';
    let isModalOpen = false;
    let checkpoints: string[] = [];
    let types = ['Bulk', 'Solid', 'Liquid'];
    let lotName: string = '';
    let selectedType: string = types[0];
    let volume: string = '';
    let maxPrice: string = '';
    let selectedDeparture: string = '';
    let selectedArrival: string = '';

    // Example data
    let tableData = [
        { name: 'Lot 1', status: 'ON_THE_WAY', volume: 16, location: 'Paris', startCheckpoint: 'Lyon', endCheckpoint: 'Montpellier', trafficManager: ['Traffic manager 1'] },
        { name: 'Lot 2', status: 'ON_THE_STOCK_EXCHANGE', volume: 3, location: 'Lyon', startCheckpoint: 'Lyon', endCheckpoint: 'Paris', trafficManager: ['Traffic manager 4'] },
        { name: 'Lot 3', status: 'PENDING', volume: 4, location: 'Marseille', startCheckpoint: 'Marseille', endCheckpoint: 'Montpellier', trafficManager: ['Traffic manager 2', 'Traffic manager 3', 'Traffic manager 4'] },
        { name: 'Lot 4', status: 'ARCHIVED', volume: 8, location: 'Montpellier', startCheckpoint: 'Paris', endCheckpoint: 'Montpellier', trafficManager: ['Traffic manager 3'] },
    ];

    // Fetch all data
    onMount(async () => {

        // GET checkpoints
        try {
            const response = await fetch(`${API_BASE_URL}/checkpoints`);
            if (response.ok)
            {
                const data = await response.json();
            
                // Extract checkpoint names
                checkpoints = data.map((checkpoint: { name: string }) => checkpoint.name);

                // Define default selected checkpoints
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
    });

    // Function to get tag color and text based on status
    function getStatusInfo(status: string): { color: string; text: string } {
        switch (status) {
            case 'PENDING':
                return { color: 'bg-green-200 text-green-800', text: '◉ En attente' };
            case 'ON_THE_WAY':
                return { color: 'bg-orange-200 text-orange-800', text: '◉ En route' };
            case 'ON_THE_STOCK_EXCHANGE':
                return { color: 'bg-yellow-200 text-yellow-800', text: '◉ En bourse' };
            case 'ARCHIVED':
                return { color: 'bg-gray-200 text-gray-800', text: '◉ Archivé' };
            default:
                return { color: 'bg-gray-200 text-gray-800', text: '🛇 Inconnu' };
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

    // Function to add lot
    function addLot() {

        // Add lot to the table
        const newLot = {
            name: lotName,
            status: 'PENDING',
            type: selectedType,
            volume: parseFloat(volume),
            maxPrice: parseFloat(maxPrice),
            location: selectedDeparture,
            startCheckpoint: selectedDeparture,
            endCheckpoint: selectedArrival,
            trafficManager: ['Traffic manager 1', 'Traffic manager 1', 'Traffic manager 3']
        };

        // Create new instance of the table
        tableData = [...tableData, newLot];

        // Reset form values
        lotName = '';
        selectedType = '';
        volume = '';
        maxPrice = '';
        selectedDeparture = checkpoints[0]; // Valeur par défaut
        selectedArrival = checkpoints[1]; // Valeur par défaut

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
</script>


<!-- Navbar -->
<Navbar/>

<main class="p-10">
    
    <section class="flex justify-between items-center mb-4">

        <!-- Title and subtitle -->
        <div>
            <h1 class="text-4xl font-bold mb-2">{title}</h1>
            <h2 class="text-2xl text-gray-600">{subtitle}</h2>
        </div>

        <!-- Create button -->
        <button class="bg-blue-500 text-white font-bold px-4 py-2 rounded flex items-center hover:bg-blue-600 transition-colors self-end"
                on:click={openModal}
        >
            <i class="fas fa-plus mr-2"></i>
            Ajouter un lot
        </button>

    </section>

    <!-- Table -->
    <div>
        <table class="table-auto w-full border-collapse border border-gray-300">
            <thead>
                <tr class="bg-gray-100">
                    <th class="border p-2 text-center">Nom</th>
                    <th class="border p-2 text-center">Status</th>
                    <th class="border p-2 text-center">Volume <span class="font-normal">(en m³)</span></th>
                    <th class="border p-2 text-center">Localisation</th>
                    <th class="border p-2 text-center">Départ / Arrivée</th>
                    <th class="border p-2 text-center">Traffic manager</th>
                </tr>
            </thead>
            <tbody>
                {#each tableData as row, index}
                    <tr class={index % 2 === 0 ? 'bg-gray-50' : 'bg-white'}>

                        <!-- Column 1 -->
                        <td class="border p-2 text-center">{row.name}</td>

                        <!-- Column 2 -->
                        <td class="border p-2 text-center">
                            <span class={`px-2 py-1 rounded ${getStatusInfo(row.status).color}`}>
                                {getStatusInfo(row.status).text}
                            </span>
                        </td>

                        <!-- Column 3 -->
                        <td class="border p-2 text-center">{row.volume}</td>

                        <!-- Column 4 -->
                        <td class="border p-2 text-center">{row.location}</td>

                        <!-- Column 5 -->
                        <td class="border p-2 text-center">
                            {row.startCheckpoint} / {row.endCheckpoint}
                        </td>

                        <!-- Column 6 -->
                        <td class="border p-2 text-center">
                            {#if row.status === 'PENDING'}
                                <select class="border border-gray-300 rounded px-2 py-1 mx-auto w-4/5">
                                    {#each row.trafficManager as trafficManagerOption}
                                        <option>{trafficManagerOption}</option>
                                    {/each}
                                </select>
                            {:else}
                                <span class="px-2 py-1 mx-auto w-4/5 block">
                                    {row.trafficManager[0]}
                                </span>
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
            <h2 class="text-2xl font-bold mb-6">Ajouter un lot</h2>

            <!-- Form -->
            <form on:submit|preventDefault={addLot}>

                <!-- Name -->
                <div class="mb-2">
                    <label class="block text-gray-700 text-sm font-bold mb-2">Nom :</label>
                    <input type="text" class="w-full border border-gray-300 p-2 rounded" placeholder="Entrez le nom du lot" bind:value={lotName} required>
                </div>

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
                           placeholder="Entrez le volume (en m³)" 
                           on:input={validateVolume} 
                           value={volume}
                           required
                    >
                </div>

                <!-- Max price -->
                <div class="mb-2">
                    <label class="block text-gray-700 text-sm font-bold mb-2">Prix maximum :</label>
                    <input type="text" 
                           class="w-full border border-gray-300 p-2 rounded" 
                           placeholder="Entrez le prix maximum (par km)" 
                           on:input={validateMaxPrice} 
                           value={volume}
                           required
                    >
                </div>

                <!-- Départ -->
                <div class="mb-2">
                    <label class="block text-gray-700 text-sm font-bold mb-2">Départ :</label>
                    <select class="w-full border border-gray-300 p-2 rounded" bind:value={selectedDeparture}>
                        {#each checkpoints as checkpoint}
                            {#if checkpoint !== selectedArrival}
                                <option value={checkpoint}>{checkpoint}</option>
                            {/if}
                        {/each}
                    </select>
                </div>

                <!-- Arrivée -->
                <div class="mb-2">
                    <label class="block text-gray-700 text-sm font-bold mb-2">Arrivée :</label>
                    <select class="w-full border border-gray-300 p-2 rounded" bind:value={selectedArrival}>
                        {#each checkpoints as checkpoint}
                            {#if checkpoint !== selectedDeparture}
                                <option value={checkpoint}>{checkpoint}</option>
                            {/if}
                        {/each}
                    </select>
                </div>

                <!-- Add button -->
                <div class="flex justify-center mt-4">
                    <button type="submit" class="bg-blue-500 text-white px-6 py-2 rounded hover:bg-blue-600">
                        <i class="fas fa-plus"></i>
                        <span class="font-bold">Ajouter</span>
                    </button>
                </div>
            </form>
        </div>
    </div>
{/if}