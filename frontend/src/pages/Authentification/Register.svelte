<script lang="ts">
    import {onMount} from "svelte";
    import axios from "axios";
    import {currentTab} from "@stores/store";
    import {faChartLine, faLock, faTruck, faUser} from "@fortawesome/free-solid-svg-icons";

    const API_BASE_URL = import.meta.env.VITE_API_BASE_URL;

    let username = '';
    let password = '';
    let role = '';

    // Roles data
    const roles = [
        { role: 'client', title: 'Client', icon: faUser},
        { role: 'traffic_manager', title: 'Traffic Manager', icon: faTruck },
        { role: 'trader', title: 'Trader', icon: faChartLine},
        { role: 'admin', title: 'Admin', icon: faLock}
    ];

    async function registerUser() {
        try {
            const response = await axios.post(`${API_BASE_URL}/auth/register`, {
                username,
                password,
                role
            });
            if (response.status === 200) {
                window.location.href = '/login';
            }
        } catch (error) {
            console.error(error);
        }
    }
</script>

<main class="p-10">
    <form>
        <div class="mb-4">
            <label for="username" class="block text-gray-700 text-sm font-bold mb-2">Username</label>
            <input type="text" id="username" name="username" class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline" bind:value={username}>
        </div>
        <div class="mb-6">
            <label for="password" class="block text-gray-700 text-sm font-bold mb-2">Password</label>
            <input type="password" id="password" name="password" class="shadow appearance-none border border-red rounded w-full py-2 px-3 text-gray-700 mb-3 leading-tight focus:outline-none focus:shadow-outline" bind:value={password}>
        </div>
        <div class="mb-6">
            <label for="role" class="block text-gray-700 text-sm font-bold mb-2">Role</label>
            <select id="role" name="role" class="shadow appearance-none border border-red rounded w-full py-2 px-3 text-gray-700 mb-3 leading-tight focus:outline-none focus:shadow-outline" bind:value={role}>
                {#each roles as { role, title }}
                    <option value={role}>{title}</option>
                {/each}
            </select>
        </div>
        <button type="button" on:click={registerUser} class="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline">Register</button>
        <button type="button" on:click={() => window.location.href = '/login'} class="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline">Already register</button>
    </form>
</main>