<script lang="ts">

    // Variables
    let title = 'Titre';
    let subtitle = 'Sous-titre';
    let currentTab = 'Tab1';
    let currentSubTab = 'Sous-onglet1';

    // Table data
    const tableData = [
        { nom: 'Tracteur 1', status: 'ON_THE_WAY' },
        { nom: 'Tracteur 2', status: 'ARRIVED' },
        { nom: 'Tracteur 3', status: 'PENDING' },
    ];

    // Function to get tag color and text based on status
    function getStatusInfo(status: string): { color: string; text: string } {
        switch (status) {
            case 'PENDING':
                return { color: 'bg-green-200 text-green-800', text: 'En attente' };
            case 'ON_THE_WAY':
                return { color: 'bg-yellow-200 text-yellow-800', text: 'En route' };
            case 'ARRIVED':
                return { color: 'bg-red-200 text-red-800', text: 'Arrivé' };
            default:
                return { color: 'bg-gray-200 text-gray-800', text: 'Inconnu' };
        }
    }

</script>


<!-- Navbar -->
<nav class="bg-gray-800 p-4 text-white">
    <ul class="flex space-x-8">
        <li><a href="#" on:click={() => currentTab = 'Tab1'} class="{currentTab === 'Tab1' ? 'font-bold' : ''}">Tab 1</a></li>
        <li><a href="#" on:click={() => currentTab = 'Tab2'} class="{currentTab === 'Tab2' ? 'font-bold' : ''}">Tab 2</a></li>
        <li><a href="#" on:click={() => currentTab = 'Tab3'} class="{currentTab === 'Tab3' ? 'font-bold' : ''}">Tab 3</a></li>
    </ul>
</nav>

<!-- Sub-Navbar -->
<nav class="bg-gray-200 p-2 text-gray-800">
    <ul class="flex space-x-8">
        <li><a href="#" on:click={() => currentSubTab = 'Sous-onglet1'} class="{currentSubTab === 'Sous-onglet1' ? 'font-bold' : ''}">Sous-onglet 1</a></li>
        <li><a href="#" on:click={() => currentSubTab = 'Sous-onglet2'} class="{currentSubTab === 'Sous-onglet2' ? 'font-bold' : ''}">Sous-onglet 2</a></li>
        <li><a href="#" on:click={() => currentSubTab = 'Sous-onglet3'} class="{currentSubTab === 'Sous-onglet3' ? 'font-bold' : ''}">Sous-onglet 3</a></li>
    </ul>
</nav>

<main class="p-10">
    <!-- Title and subtitle -->
    <section>
        <h1 class="text-4xl font-bold mb-4">{title}</h1>
        <h2 class="text-2xl mb-4 text-gray-600">{subtitle}</h2>
    </section>

    <!-- Table -->
    <div>
        <table class="table-auto w-full border-collapse border border-gray-300">
            <thead>
                <tr class="bg-gray-100">
                    <th class="border p-2 text-center">Nom</th>
                    <th class="border p-2 text-center">Status</th>
                    <th class="border p-2 text-center">Bourse</th>
                    <th class="border p-2 text-center">Action</th>
                    <th class="border p-2 text-center">Assigner une route</th>
                </tr>
            </thead>
            <tbody>
                {#each tableData as row, index}
                    <tr class={index % 2 === 0 ? 'bg-gray-50' : 'bg-white'}>

                        <!-- Column 1 -->
                        <td class="border p-2 text-center">{row.nom}</td>

                        <!-- Column 2 -->
                        <td class="border p-2 text-center">
                            <span class={`px-2 py-1 rounded ${getStatusInfo(row.status).color}`}>
                                {getStatusInfo(row.status).text}
                            </span>
                        </td>

                        <!-- Column 3 -->
                        <td class="border p-2 text-center">
                            <button class="bg-blue-500 text-white rounded-full w-10 h-10 flex items-center justify-center mx-auto hover:bg-blue-600 transition-colors">
                                <i class="fas fa-plus"></i>
                            </button>
                        </td>

                        <!-- Column 4 -->
                        <td class="border p-2 text-center">
                            <button class="bg-gray-200 text-gray-800 px-4 py-2 flex items-center font-bold mx-auto hover:bg-gray-300 transition-colors rounded-md">
                                <i class="fas fa-truck mr-2"></i>
                                Démarrer
                            </button>
                        </td>

                        <!-- Column 5 -->
                        <td class="border p-2 text-center">
                            <select class="border border-gray-300 rounded px-2 py-1 mx-auto w-4/5">
                                <option>Lyon - Marseille</option>
                                <option>Paris - Marseille</option>
                                <option>Toulouse - Paris</option>
                            </select>
                        </td>
                    </tr>
                {/each}
            </tbody>
        </table>
    </div>
</main>