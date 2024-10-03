<script lang="ts">
    import Navbar from '@components/Navbar.svelte';
    import TrafficManagerNavbar from '@components/TrafficManagerNavbar.svelte';

    // Variables
    let title: string = 'Lot management';
    let subtitle: string = 'Track the status of your lots in real time.';

    // Function to get tag color and text based on status
    function getStatusInfo(status: string): { color: string; text: string } {
        switch (status) {
            case 'PENDING':
                return { color: 'bg-green-200 text-green-800', text: '◉ Pending' };
            case 'ON_THE_WAY':
                return { color: 'bg-orange-200 text-orange-800', text: '◉ On the road' };
            case 'ON_THE_STOCK_EXCHANGE':
                return { color: 'bg-yellow-200 text-yellow-800', text: '◉ On the stock exchange' };
            case 'ARCHIVED':
                return { color: 'bg-gray-200 text-gray-800', text: '◉ Archived' };
            default:
                return { color: 'bg-gray-200 text-gray-800', text: '🛇 Unknown' };
        }
    }

    // Example data
    const tableData = [
        { name: 'Lot 1', status: 'ON_THE_WAY', volume: 16, location: 'Paris', startCheckpoint: 'Lyon', endCheckpoint: 'Montpellier', tractor: ['tractor 1'] },
        { name: 'Lot 2', status: 'ON_THE_STOCK_EXCHANGE', volume: 3, location: 'Lyon', startCheckpoint: 'Lyon', endCheckpoint: 'Paris', tractor: ['tractor 4'] },
        { name: 'Lot 3', status: 'PENDING', volume: 4, location: 'Marseille', startCheckpoint: 'Marseille', endCheckpoint: 'Montpellier', tractor: ['tractor 2', 'tractor 3', 'tractor 4'] },
        { name: 'Lot 4', status: 'ARCHIVED', volume: 8, location: 'Montpellier', startCheckpoint: 'Paris', endCheckpoint: 'Montpellier', tractor: ['tractor 3'] },
    ];
</script>


<!-- Navbar -->
<Navbar/>
<TrafficManagerNavbar/>

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
                <th class="border p-2 text-center">Name</th>
                <th class="border p-2 text-center">Status</th>
                <th class="border p-2 text-center">Volume <span class="font-normal">(in m³)</span></th>
                <th class="border p-2 text-center">Location</th>
                <th class="border p-2 text-center">Departure / Arrival</th>
                <th class="border p-2 text-center">Tractor</th>
                <th class="border p-2 text-center">Actions</th>
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
                                {#each row.tractor as tractorOption}
                                    <option>{tractorOption}</option>
                                {/each}
                            </select>
                        {:else}
                                <span class="px-2 py-1 mx-auto w-4/5 block">
                                    {row.tractor[0]}
                                </span>
                        {/if}
                    </td>

                    <!-- Column 7 -->
                    <td class="border p-2 text-center">
                        {#if row.status === 'PENDING'}
                            <div class="flex flex-wrap justify-center space-x-2">
                                <button class="bg-blue-200 text-blue-800 px-4 py-2 flex items-center font-bold hover:bg-blue-300 transition-colors rounded-md">
                                    <i class="fas fa-plus mr-2"></i>
                                    Stock exchange
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