<script lang="ts">
    import { onMount } from 'svelte';
    import 'leaflet/dist/leaflet.css';
    import '@fortawesome/fontawesome-free/css/all.css';
    import L from 'leaflet';
    import Navbar from "@components/Navbar.svelte";

    // Enum for marker types
    enum MarkerType {
        CHECKPOINT,
        TRACTOR,
        LOT
    }

    // Variables
    let map;
    let mapContainer;

    // Function to initialize the map with a minimalist style
    function initializeMap() {
        map = L.map(mapContainer).setView([48.8566, 2.3522], 5);

        // Minimalist tile layer (CartoDB Positron)
        L.tileLayer('https://{s}.basemaps.cartocdn.com/light_all/{z}/{x}/{y}{r}.png', {
            attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
        }).addTo(map);
    }

    // Add markers with FontAwesome icons
    function addMarkers(data) {

        // Clean up existing markers before adding new ones
        map.eachLayer((layer) => {
            if (layer instanceof L.Marker) {
                map.removeLayer(layer);
            }
        });

        // Loop through data to add markers dynamically
        data.forEach(point => {

            if (point.type == MarkerType.CHECKPOINT)
            {
                let icon = L.divIcon({
                    className: `fa fa-2x text-red-500`,
                    iconSize: [32, 32],
                    iconAnchor: [16, 16],
                    html: `<div style="text-align: center;">
                            <i class="fa fa-location-pin"></i>
                        </div>`,
                });
                L.marker([point.lat, point.lon], {icon})
                    .addTo(map)
                    .bindPopup(`<span style="font-weight: bold; color: gray;">${point.name}</span>`);
            }
        });
    }

    // Simulate API call to fetch data for current map bounds
    async function fetchDynamicData() {
        const dummyData = [
            {
                name: 'Checkpoint 1', type: MarkerType.CHECKPOINT, lat: 48.8566, lon: 2.3522
            }
        ];

        // Simulate API delay
        return new Promise(resolve => {
            setTimeout(() => resolve(dummyData), 1000);
        });
    }

    // Update markers based on the current map bounds
    async function updateMarkers() {
        const data = await fetchDynamicData();
        addMarkers(data);
    }

    // Initialize the map when the component is mounted
    onMount(() => {
        initializeMap();

        // Fetch data when the map is moved or zoomed
        map.on('moveend', updateMarkers);
        map.on('zoomend', updateMarkers);

        // Initial load of markers
        updateMarkers();
    });

</script>


<!-- Navbar -->
<div class="navbar">
    <Navbar/>
</div>

<!-- Map container -->
<div id="map" bind:this={mapContainer}></div>


<style>
    #map {
        height: 100vh;
        width: 100%;
        position: relative;
        z-index: 0;
    }

    /* Ensure the navbar stays above the map */
    .navbar {
        position: fixed;
        top: 0;
        left: 0;
        right: 0;
        z-index: 50;
    }

</style>