<script lang="ts">
    import Navbar from '@components/Navbar.svelte';

    // Variables
    let title: string = 'Title';
    let subtitle: string = 'This is a subtitle.';
    let isTractorsModalOpen = false;
    let isLotsModalOpen = false;

    // Example data
    const tractorsData = [
        { id: 1, expirationDate: 1695564000000, type: 'Bulk', spaceAvailable: 500, minPrice: 2.5, currentPrice: 3.2 },
        { id: 2, expirationDate: 1698242400000, type: 'Liquid', spaceAvailable: 800, minPrice: 1.8, currentPrice: 2.1 },
        { id: 3, expirationDate: 1700834400000, type: 'Solid', spaceAvailable: 300, minPrice: 3.0, currentPrice: 3.5 }
    ];
    const lotsData = [
        { id: 1, expirationDate: 1695564000000, type: 'Bulk', volume: 80, maxPrice: 1.5, currentPrice: 1.2 },
        { id: 2, expirationDate: 1698242400000, type: 'Liquid', volume: 50, maxPrice: 0.8, currentPrice: 1.7 },
        { id: 3, expirationDate: 1700834400000, type: 'Solid', volume: 30, maxPrice: 2.0, currentPrice: 0.5 }
    ];

    // Function to format timestamp into DD/MM/YYYY
    const formatDate = (timestamp: number) => {
        const date = new Date(timestamp);
        return `${date.getDate().toString().padStart(2, '0')}/${(date.getMonth() + 1).toString().padStart(2, '0')}/${date.getFullYear()}`;
    };

    // Function to open tractors modal
    function openTractorsModal() {
        isTractorsModalOpen = true;
    }

    // Function to close tractors modal
    function closeTractorsModal() {
        isTractorsModalOpen = false;
    }

    // Function to open lots modal
    function openLotsModal() {
        isTractorsModalOpen = true;
    }

    // Function to close lots modal
    function closeLotsModal() {
        isTractorsModalOpen = false;
    }

    // Function to add bid on tractors
    function bidTractors() {
        console.log("Form submitted!");
        closeTractorsModal();
    }

    // Function to add bid on lots
    function bidLots() {
        console.log("Form submitted!");
        closeLotsModal();
    }

</script>


<!-- Navbar -->
<Navbar/>

<main class="p-10">

    <!-- Title and subtitle -->
    <section>
        <h1 class="text-4xl font-bold mb-4">{title}</h1>
        <h2 class="text-2xl mb-8 text-gray-600">{subtitle}</h2>
    </section>

    <div class="flex">

        <!-- Left part -->
        <div class="w-1/2 pr-8 border-r border-gray-300">
            <h2 class="text-2xl text-gray-800 font-bold mb-4">
                <i class="fas fa-truck mr-2"></i>
                Offre de tracteurs
            </h2>
            <table class="table-auto w-full border-collapse border border-gray-300">
                <thead>
                    <tr class="bg-gray-100">
                        <th class="border p-2 text-center">ID</th>
                        <th class="border p-2 text-center">Date d'expiration</th>
                        <th class="border p-2 text-center">Type</th>
                        <th class="border p-2 text-center">Espace disponible<br><span class="font-normal">(en m³)</span></th>
                        <th class="border p-2 text-center">Prix minimum<br><span class="font-normal">(en €/km)</span></th>
                        <th class="border p-2 text-center">Prix actuel<br><span class="font-normal">(en €/km)</span></th>
                        <th class="border p-2 text-center">Actions</th>
                    </tr>
                </thead>
                <tbody>
                    {#each tractorsData as row, index}
                        <tr class={index % 2 === 0 ? 'bg-gray-50' : 'bg-white'}>

                            <!-- Column 1 -->
                            <td class="border p-2 text-center">{row.id}</td>
                            
                            <!-- Column 2 -->
                            <td class="border p-2 text-center">{formatDate(row.expirationDate)}</td>

                            <!-- Column 3 -->
                            <td class="border p-2 text-center">{row.type}</td>
                            
                            <!-- Column 4 -->
                            <td class="border p-2 text-center">{row.spaceAvailable}</td>
                            
                            <!-- Column 5 -->
                            <td class="border p-2 text-center">{row.minPrice.toFixed(2)}</td>
                            
                            <!-- Column 6 -->
                            <td class="border p-2 text-center">{row.currentPrice.toFixed(2)}</td>
                            
                            <!-- Column 7 -->
                            <td class="border p-2 text-center">
                                <div class="flex flex-wrap justify-center space-x-2">
                                    <button class="bg-blue-200 text-blue-800 px-4 py-2 flex items-center font-bold hover:bg-blue-300 transition-colors rounded-md"
                                            on:click={openTractorsModal}
                                    >
                                        <i class="fas fa-coins mr-2"></i>
                                        Enchérir
                                    </button>
                                </div>
                            </td>

                        </tr>
                    {/each}
                </tbody>
            </table>
        </div>

        <!-- Right part -->
        <div class="w-1/2 pl-8">
            <h2 class="text-2xl text-gray-800 font-bold mb-4">
                <i class="fas fa-box-open mr-2"></i>
                Offres de lots
            </h2>

            <table class="table-auto w-full border-collapse border border-gray-300">
                <thead>
                    <tr class="bg-gray-100">
                        <th class="border p-2 text-center">ID</th>
                        <th class="border p-2 text-center">Date d'expiration</th>
                        <th class="border p-2 text-center">Type</th>
                        <th class="border p-2 text-center">Volume<br><span class="font-normal">(en m³)</span></th>
                        <th class="border p-2 text-center">Prix maximum<br><span class="font-normal">(en €/km)</span></th>
                        <th class="border p-2 text-center">Prix actuel<br><span class="font-normal">(en €/km)</span></th>
                        <th class="border p-2 text-center">Actions</th>
                    </tr>
                </thead>
                <tbody>
                    {#each lotsData as row, index}
                        <tr class={index % 2 === 0 ? 'bg-gray-50' : 'bg-white'}>

                            <!-- Column 1 -->
                            <td class="border p-2 text-center">{row.id}</td>
                            
                            <!-- Column 2 -->
                            <td class="border p-2 text-center">{formatDate(row.expirationDate)}</td>

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
                                <div class="flex flex-wrap justify-center space-x-2">
                                    <button class="bg-blue-200 text-blue-800 px-4 py-2 flex items-center font-bold hover:bg-blue-300 transition-colors rounded-md"
                                            on:click={openLotsModal}
                                    >
                                        <i class="fas fa-coins mr-2"></i>
                                        Enchérir
                                    </button>
                                </div>
                            </td>

                        </tr>
                    {/each}
                </tbody>
            </table>
        </div>

    </div>
</main>

<!-- svelte-ignore a11y-no-static-element-interactions -->
<!-- svelte-ignore a11y-no-static-element-interactions -->
<!-- svelte-ignore a11y-click-events-have-key-events -->
<!-- svelte-ignore a11y-label-has-associated-control -->

{#if isTractorsModalOpen}

    <div class="fixed inset-0 bg-black bg-opacity-50 flex justify-center items-center z-50" on:click={closeTractorsModal}>

        <div class="bg-white p-6 rounded-lg shadow-lg w-1/3" on:click|stopPropagation>

            <!-- Close Button -->
            <button class="absolute top-2 right-2 text-gray-500 hover:text-gray-800" on:click={closeTractorsModal}>
                &times;
            </button>

            <!-- Modal Title -->
            <h2 class="text-2xl font-bold mb-6">Associer un lot à l'enchère</h2>

            <!-- Form -->
            <form on:submit|preventDefault={bidTractors}>

                <!-- Current price -->
                <div class="mb-2 flex justify-left items-center">
                    <span class="text-xl mr-2 font-bold">Prix actuel :</span>
                    <span class="text-xl font-normal">12 €/km</span>
                </div>

                <!-- Add button -->
                <div class="flex justify-center mt-4">
                    <button type="submit" class="bg-blue-500 text-white px-6 py-2 rounded hover:bg-blue-600">
                        <i class="fas fa-check"></i>
                        <span class="font-bold">Valider</span>
                    </button>
                </div>
            </form>
        </div>
    </div>
{/if}

<!-- svelte-ignore a11y-no-static-element-interactions -->
<!-- svelte-ignore a11y-no-static-element-interactions -->
<!-- svelte-ignore a11y-click-events-have-key-events -->
<!-- svelte-ignore a11y-label-has-associated-control -->

{#if isLotsModalOpen}

    <div class="fixed inset-0 bg-black bg-opacity-50 flex justify-center items-center z-50" on:click={closeLotsModal}>

        <div class="bg-white p-6 rounded-lg shadow-lg w-1/3" on:click|stopPropagation>

            <!-- Close Button -->
            <button class="absolute top-2 right-2 text-gray-500 hover:text-gray-800" on:click={closeLotsModal}>
                &times;
            </button>

            <!-- Modal Title -->
            <h2 class="text-2xl font-bold mb-6">Ajouter un lot</h2>

            <!-- Form -->
            <form on:submit|preventDefault={bidLots}>

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