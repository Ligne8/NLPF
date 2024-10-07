<script lang="ts">
    import Navbar from '@components/Navbar.svelte';
    import TraderNavbar from '@components/TraderNavbar.svelte';

    // Variables
    let title: string = 'Lot offers';
    let subtitle: string = 'Create lot offers in real time.';

    // Example data
    const tableData = [
        { id: 1, expirationDate: null, type: 'Bulk', volume: 500, maxPrice: 2.5, currentPrice: 3.2 },
        { id: 2, expirationDate: 1698242400000, type: 'Liquid', volume: 800, maxPrice: 1.8, currentPrice: 2.1 },
        { id: 3, expirationDate: null, type: 'Solid', volume: 300, maxPrice: 3.0, currentPrice: 3.5 }
    ];

    // Function to format timestamp into DD/MM/YYYY
    const formatDate = (timestamp: number) => {
        const date = new Date(timestamp);
        return `${date.getDate().toString().padStart(2, '0')}/${(date.getMonth() + 1).toString().padStart(2, '0')}/${date.getFullYear()}`;
    };

</script>


<!-- Navbar -->
<Navbar/>
<TraderNavbar/>

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
                    <th class="border p-2 text-center">ID</th>
                    <th class="border p-2 text-center">Expiration date</th>
                    <th class="border p-2 text-center">Type</th>
                    <th class="border p-2 text-center">Volume<br><span class="font-normal">(in m³)</span></th>
                    <th class="border p-2 text-center">Maximum price<br><span class="font-normal">(in €/km)</span></th>
                    <th class="border p-2 text-center">Current price<br><span class="font-normal">(in €/km)</span></th>
                    <th class="border p-2 text-center">Actions</th>
                </tr>
            </thead>
            <tbody>
                {#each tableData as row, index}
                    <tr class={index % 2 === 0 ? 'bg-gray-50' : 'bg-white'}>

                        <!-- Column 1 -->
                        <td class="border p-2 text-center">{row.id}</td>
                        
                        <!-- Column 2 -->
                        <td class="border p-2 text-center">
                            {#if !row.expirationDate}
                                <input type="date" bind:value={row.expirationDate} class="border p-1 text-center"/>
                            {:else}
                                {formatDate(row.expirationDate)}
                            {/if}
                        </td>

                        <!-- Column 3 -->
                        <td class="border p-2 text-center">{row.type}</td>
                        
                        <!-- Column 4 -->
                        <td class="border p-2 text-center">{row.volume}</td>
                        
                        <!-- Column 5 -->
                        <td class="border p-2 text-center">{row.maxPrice.toFixed(2)}</td>
                        
                        <!-- Column 6 -->
                        <td class="border p-2 text-center">{row.currentPrice.toFixed(2)}</td>
                        
                        <!-- Column 7 -->
                        <td class="border p-2 text-center">
                            <div class="flex flex-wrap justify-center space-x-2 space-y-2">
                                {#if row.expirationDate}
                                    <button class="bg-blue-200 text-blue-800 px-4 py-2 flex items-center font-bold hover:bg-blue-300 transition-colors rounded-md">
                                        <i class="fas fa-plus mr-2"></i>
                                        Offer
                                    </button>
                                {:else}
                                    <span class="text-gray-500">-</span>
                                {/if}
                            </div>
                        </td>

                    </tr>
                {/each}
            </tbody>
        </table>
    </div>
</main>