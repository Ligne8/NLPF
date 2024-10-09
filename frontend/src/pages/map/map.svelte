<script lang="ts">
    import { onMount } from 'svelte';
    import 'leaflet/dist/leaflet.css';
    import '@fortawesome/fontawesome-free/css/all.css';
    import L from 'leaflet';
    import 'leaflet-arrowheads';
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
    let showCheckpoints = true;
    let showLots = true;
    let showTractors = true;
    let showRoutes = true;

    // Function to initialize the map with a minimalist style
    function initializeMap() {
        map = L.map(mapContainer).setView([44.9068, 3.9598], 5);

        // Minimalist tile layer (CartoDB Positron)
        L.tileLayer('https://{s}.basemaps.cartocdn.com/light_all/{z}/{x}/{y}{r}.png', {
            attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
        }).addTo(map);
    }

    // Clean up existing markers and polylines
    function cleanMarkers() {
        map.eachLayer((layer) => {
            if (layer instanceof L.Marker || layer instanceof L.Polyline) {
                map.removeLayer(layer);
            }
        });
    }

    // Add markers with FontAwesome icons
    async function addMarkers(data) {

        // Loop through data to add markers dynamically
        for (const elt of data)
        {
            let marker = null;

            // Checkpoints
            if (elt.type == MarkerType.CHECKPOINT && showCheckpoints)
            {
                let icon = L.divIcon({
                    className: `fa fa-xl text-purple-500`,
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

            // Lots
            else if (elt.type == MarkerType.LOT && showLots)
            {
                let icon = L.divIcon({
                    className: `fa fa-3x text-yellow-500`,
                    iconSize: [32, 32],
                    iconAnchor: [16, 16],
                    html: `<div style="text-align: center;">
                            <i class="fa fa-box"></i>
                        </div>`,
                });
                console.log(elt)
                marker = L.marker([elt.current_checkpoint.latitude, elt.current_checkpoint.longitude], {icon})
                    .addTo(map)
                    .bindPopup(`<p style="font-weight: bold; color: gray; text-align: center">Lot : ${elt.resource_type}<br>
                                    <span style="font-weight: normal;">
                                        ${elt.current_checkpoint.name} (${elt.current_checkpoint.country})
                                    </span>
                                </p>
                                <p style="font-weight: bold; color: gray; text-align: left; margin: 0;">Departure :
                                    <span style="font-weight: normal;">
                                        ${elt.start_checkpoint.name} (${elt.start_checkpoint.country})
                                    </span>
                                </p>
                                <p style="font-weight: bold; color: gray; text-align: left; margin: 0;">Arrival :
                                    <span style="font-weight: normal;">
                                        ${elt.end_checkpoint.name} (${elt.end_checkpoint.country})
                                    </span>
                                </p>`);
            }

            // Routes and tractors
            else if (elt.type == MarkerType.TRACTOR)
            {
                // Routes
                if (elt.route !== null && elt.state === 'in_transit' && showRoutes)
                {
                    const checkpoints = await getCheckpointsByRouteID(elt.route_id);
                    let coords = [];
                    for (const c of checkpoints)
                    {
                        coords.push([c.checkpoint.latitude, c.checkpoint.longitude]);
                    }
                    const polyline = L.polyline(coords, { color: '#3b82f6', weight: 2 }).addTo(map);
                    polyline.arrowheads({
                        size: '15px',
                        frequency: 'allvertices',
                        fill: true,
                        color: '#3b82f6'
                    });
                }

                // Tractors
                if (showTractors)
                {
                    let icon = L.divIcon({
                        className: `fa fa-3x text-gray-800`,
                        iconSize: [32, 32],
                        iconAnchor: [16, 16],
                        html: `<div style="text-align: center;">
                                <i class="fa fa-truck"></i>
                            </div>`,
                    });
                    marker = L.marker([elt.current_checkpoint.latitude, elt.current_checkpoint.longitude], {icon})
                        .addTo(map)
                        .bindPopup(`<p style="font-weight: bold; color: gray; text-align: center">${elt.name}<br>
                                        <span style="font-weight: normal;">
                                            ${elt.current_checkpoint.name} (${elt.current_checkpoint.country})
                                        </span>
                                    </p>
                                    <p style="font-weight: bold; color: gray; text-align: left; margin: 0;">Departure :
                                        <span style="font-weight: normal;">
                                            ${elt.start_checkpoint.name} (${elt.start_checkpoint.country})
                                        </span>
                                    </p>
                                    <p style="font-weight: bold; color: gray; text-align: left; margin: 0;">Arrival :
                                        <span style="font-weight: normal;">
                                            ${elt.end_checkpoint.name} (${elt.end_checkpoint.country})
                                        </span>
                                    </p>`);
                }
            }

            // Animation
            if (marker && elt.state === 'available')
                animateMarker(marker);
        }
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

    // Fetch lots depending on user ID
    async function fetchLots() {
        let route: string = "";
        if ($userRole === "traffic_manager")
            route = `${API_BASE_URL}/lots/traffic_manager/${$userId}`
        else if ($userRole === "client")
            route = `${API_BASE_URL}/lots/owner/${$userId}`
        else
            return;
        try {
            const response = await axios.get(route);
            const lots = response.data.map(lot => ({
                ...lot,
                type: MarkerType.LOT
            }));
            return lots;
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

    // Get checkpoints of the route by ID
    async function getCheckpointsByRouteID(routeId: number) {
        try {
            const response = await axios.get(`${API_BASE_URL}/routes/${routeId}/checkpoints`);
            return response.data;
        } catch (err) {
            console.error(err);
        }
    }

    // Update markers based on the current map bounds
    async function updateMarkers() {
        cleanMarkers();
        const checkpoints = await fetchCheckpoints();
        const lots = await fetchLots();
        const tractors = await fetchTractors();
        addMarkers(checkpoints);
        addMarkers(lots);
        addMarkers(tractors);
    }

    // Toggle show checkpoints state
    function toggleCheckpoints() {
        showCheckpoints = !showCheckpoints;
        updateMarkers();
    }

    // Toggle show lots state
    function toggleLots() {
        showLots = !showLots;
        updateMarkers();
    }

    // Toggle show tractors state
    function toggleTractors() {
        showTractors = !showTractors;
        updateMarkers();
    }

    // Toggle show tractors state
    function toggleRoutes() {
        showRoutes = !showRoutes;
        updateMarkers();
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

    // Function to animate marker
    function animateMarker(marker) {
        const iconElement = marker.getElement().querySelector('.fa');
        let scale = 1;
        let direction = 1;
        setInterval(() => {
            scale += 0.05 * direction;
            if (scale >= 1.2) direction = -1;
            if (scale <= 0.7) direction = 1;
            iconElement.style.transform = `scale(${scale})`;
        }, 50);
    }

</script>


<!-- Navbar -->
<div class="navbar">
    <Navbar/>
</div>

<!-- Map container -->
<div id="map" bind:this={mapContainer}></div>

<!-- Layer Controls -->
<div class="absolute bottom-5 right-5 bg-white border border-gray-300 p-3 rounded shadow-lg z-100 flex flex-col space-y-2">
    <button 
        class={`text-white py-2 px-4 rounded hover:opacity-80 transition ${showCheckpoints ? 'bg-purple-500' : 'bg-red-500'}`}
        on:click={toggleCheckpoints}>
        <i class={`fa ${showCheckpoints ? 'fa-eye' : 'fa-eye-slash'} mr-2`}></i>
        Checkpoints
    </button>
    <button 
        class={`text-white py-2 px-4 rounded hover:opacity-80 transition ${showLots ? 'bg-yellow-500' : 'bg-red-500'}`}
        on:click={toggleLots}>
        <i class={`fa ${showLots ? 'fa-eye' : 'fa-eye-slash'} mr-2`}></i>
        Lots
    </button>
    <button 
        class={`text-white py-2 px-4 rounded hover:opacity-80 transition ${showTractors ? 'bg-gray-800' : 'bg-red-500'}`}
        on:click={toggleTractors}>
        <i class={`fa ${showTractors ? 'fa-eye' : 'fa-eye-slash'} mr-2`}></i>
        Tractors
    </button>
    <button 
        class={`text-white py-2 px-4 rounded hover:opacity-80 transition ${showRoutes ? 'bg-blue-500' : 'bg-red-500'}`}
        on:click={toggleRoutes}>
        <i class={`fa ${showRoutes ? 'fa-eye' : 'fa-eye-slash'} mr-2`}></i>
        Routes
    </button>
</div>


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