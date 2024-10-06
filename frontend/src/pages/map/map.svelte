<script lang="ts">
    import { onMount } from 'svelte';
    import 'leaflet/dist/leaflet.css';
    import '@fortawesome/fontawesome-free/css/all.css';
    import L from 'leaflet';
    import Navbar from "@components/Navbar.svelte";
    import axios from "axios";
    import { userRole, userId } from '@stores/store.js';

    // Enum for marker types
    enum MarkerType {
        CHECKPOINT,
        TRACTOR,
        LOT
    }

    // Variables
    let map;
    let mapContainer;
    const API_BASE_URL = import.meta.env.VITE_API_BASE_URL;

    // Function to initialize the map with a minimalist style
    function initializeMap() {
        map = L.map(mapContainer).setView([44.9068, 3.9598], 5);

        // Minimalist tile layer (CartoDB Positron)
        L.tileLayer('https://{s}.basemaps.cartocdn.com/light_all/{z}/{x}/{y}{r}.png', {
            attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
        }).addTo(map);
    }

    // Clean up existing markers 
    function cleanMarkers() {
        map.eachLayer((layer) => {
            if (layer instanceof L.Marker) {
                map.removeLayer(layer);
            }
        });
    }

    // Add markers with FontAwesome icons
    function addMarkers(data) {

        // Loop through data to add markers dynamically
        data.forEach(elt => {

            // Checkpoints
            if (elt.type == MarkerType.CHECKPOINT)
            {
                let icon = L.divIcon({
                    className: `fa fa-xl text-red-500`,
                    iconSize: [16, 16],
                    iconAnchor: [8, 8],
                    html: `<div style="text-align: center;">
                            <i class="fa fa-location-pin"></i>
                        </div>`,
                });
                L.marker([elt.latitude, elt.longitude], {icon})
                    .addTo(map)
                    .bindPopup(`<p style="font-weight: bold; color: gray;">${elt.name}
                                    <span style="font-weight: normal;">(${elt.country})</span>
                                </p>`);
            }

            // Tractors
            if (elt.type == MarkerType.TRACTOR)
            {
                if (elt.route !== null)
                {
                    // Trace the line between start_checkpoint and end_checkpoint
                    const startCoords = [elt.start_checkpoint.latitude, elt.start_checkpoint.longitude];
                    const endCoords = [elt.end_checkpoint.latitude, elt.end_checkpoint.longitude];

                    // Create the polyline
                    const polyline = L.polyline([startCoords, endCoords], { color: 'blue', weight: 3 }).addTo(map);
                }

                let icon = L.divIcon({
                    className: `fa fa-3x text-gray-800`,
                    iconSize: [32, 32],
                    iconAnchor: [16, 16],
                    html: `<div style="text-align: center;">
                            <i class="fa fa-truck"></i>
                        </div>`,
                });
                L.marker([elt.current_checkpoint.latitude, elt.current_checkpoint.longitude], {icon})
                    .addTo(map)
                    .bindPopup(`<p style="font-weight: bold; color: gray; text-align: center">${elt.name}<br>
                                    <span style="font-weight: normal;">
                                        ${elt.current_checkpoint.name} (${elt.current_checkpoint.country})
                                    </span>
                                </p>`);
            }
        });
    }

    // Fetch all checkpoints
    async function fetchCheckpoints() {
        try {
            const response = await axios.get(`${API_BASE_URL}/checkpoints`);
            const checkpoints = response.data.map(point => ({
                ...point,
                type: MarkerType.CHECKPOINT
            }));
            return checkpoints;
        } catch (err) {
            console.error(err);
        }
    }

    // Fetch tractors depending on user ID
    async function fetchTractors() {
        let route: string = "";
        if ($userRole === "traffic_manager")
            route = `${API_BASE_URL}/tractors/trafficManager/${$userId}`
        else if ($userRole === "client")
            route = `${API_BASE_URL}/tractors/owner/${$userId}`
        else
            return;
        try {
            const response = await axios.get(route);
            const tractors = response.data.map(tractor => ({
                ...tractor,
                type: MarkerType.TRACTOR
            }));
            return tractors;
        } catch (err) {
            console.error(err);
        }
    }

    // Update markers based on the current map bounds
    async function updateMarkers() {
        cleanMarkers();
        const checkpoints = await fetchCheckpoints();
        const tractors = await fetchTractors();
        addMarkers(checkpoints);
        addMarkers(tractors);
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