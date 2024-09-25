<script lang="ts">
    import Navbar from '../../components/Navbar.svelte';

    // Variables
    let title: string = 'Gestion des Lots';
    let subtitle: string = 'Suivez l’état de vos lots en temps réel.';

    // Function to get tag color and text based on status
    function getStatusInfo(status: string): { color: string; text: string } {
        switch (status) {
            case 'AVAILABLE':
                return { color: 'bg-green-200 text-green-800', text: '◉ Disponible' };
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

    // Example data
    const tableData = [
        { name: 'Lot 1', status: 'ON_THE_WAY', volume: 16, location: 'Paris', startCheckpoint: 'Lyon', endCheckpoint: 'Montpellier', trafficManager: ['Traffic manager 1'] },
        { name: 'Lot 2', status: 'ON_THE_STOCK_EXCHANGE', volume: 3, location: 'Lyon', startCheckpoint: 'Lyon', endCheckpoint: 'Paris', trafficManager: ['Traffic manager 4'] },
        { name: 'Lot 3', status: 'AVAILABLE', volume: 4, location: 'Marseille', startCheckpoint: 'Marseille', endCheckpoint: 'Montpellier', trafficManager: ['Traffic manager 2', 'Traffic manager 3', 'Traffic manager 4'] },
        { name: 'Lot 4', status: 'ARCHIVED', volume: 8, location: 'Montpellier', startCheckpoint: 'Paris', endCheckpoint: 'Montpellier', trafficManager: ['Traffic manager 3'] },
    ];
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
        <button class="bg-blue-500 text-white font-bold px-4 py-2 rounded flex items-center hover:bg-blue-600 transition-colors self-end">
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
                    <th class="border p-2 text-center">Volume <span class="font-normal">(en m³)</th>
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
                            {#if row.status === 'AVAILABLE'}
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