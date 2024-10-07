<script lang="ts">
    import Navbar from '@components/Navbar.svelte';
    import TrafficManagerNavbar from '@components/TrafficManagerNavbar.svelte';

    // Variables
    let title: string = 'Lot management';
    let subtitle: string = 'Track the status of your lots in real time.';
    let selectedStatus: string = 'all';
    let sortOption: string = 'none';

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
            default:
                return { color: 'bg-gray-200 text-gray-800', text: '🛇 Unknown' };
        }
    }

    // Example data
    const tableData = [
        { state: 'in_transit', volume: 16, currentCheckpoint: 'Paris', startCheckpoint: 'Lyon', endCheckpoint: 'Montpellier', tractor: ['tractor 1'] },
        { state: 'on_market', volume: 3, currentCheckpoint: 'Lyon', startCheckpoint: 'Lyon', endCheckpoint: 'Paris', tractor: ['tractor 4'] },
        { state: 'pending', volume: 4, currentCheckpoint: 'Marseille', startCheckpoint: 'Marseille', endCheckpoint: 'Montpellier', tractor: ['tractor 2', 'tractor 3', 'tractor 4'] },
        { state: 'archived', volume: 8, currentCheckpoint: 'Montpellier', startCheckpoint: 'Paris', endCheckpoint: 'Marseille', tractor: ['tractor 3'] },
        { state: 'available', volume: 4, currentCheckpoint: 'Lyon', startCheckpoint: 'Marseille', endCheckpoint: 'Montpellier', tractor: ['tractor 2'] },
        { state: 'available', volume: 2, currentCheckpoint: 'Montpellier', startCheckpoint: 'Lyon', endCheckpoint: 'Lyon', tractor: ['tractor 3'] }
    ];

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
<TrafficManagerNavbar/>

<main class="p-10 mt-40">

    <!-- Title and subtitle -->
    <div class="mb-2">
        <h1 class="text-4xl font-bold mb-4">{title}</h1>
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

    <!-- Table -->
    <div>
        <table class="table-auto w-full border-collapse border border-gray-300">
            <thead>
            <tr class="bg-gray-100">
                <th class="border p-2 text-center">Status</th>
                <th class="border p-2 text-center">Volume <span class="font-normal">(in m³)</span></th>
                <th class="border p-2 text-center">Location</th>
                <th class="border p-2 text-center">Departure / Arrival</th>
                <th class="border p-2 text-center">Tractor</th>
                <th class="border p-2 text-center">Actions</th>
            </tr>
            </thead>
            <tbody>
            {#each sortedData as row, index}
                <tr class={index % 2 === 0 ? 'bg-gray-50' : 'bg-white'}>

                    <!-- Column 1 -->
                    <td class="border p-2 text-center">
                            <span class={`px-2 py-1 rounded ${getStatusInfo(row.state).color}`}>
                                {getStatusInfo(row.state).text}
                            </span>
                    </td>

                    <!-- Column 2 -->
                    <td class="border p-2 text-center">{row.volume}</td>

                    <!-- Column 3 -->
                    <td class="border p-2 text-center">{row.currentCheckpoint}</td>

                    <!-- Column 4 -->
                    <td class="border p-2 text-center">
                        {row.startCheckpoint} / {row.endCheckpoint}
                    </td>

                    <!-- Column 5 -->
                    <td class="border p-2 text-center">
                        {#if row.state === 'pending'}
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

                    <!-- Column 6 -->
                    <td class="border p-2 text-center">
                        {#if row.state === 'pending'}
                            <div class="flex flex-wrap justify-center space-x-2 space-y-2">
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