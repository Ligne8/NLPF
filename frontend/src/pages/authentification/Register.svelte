<script lang="ts">
    import axios from "axios";
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

<main class="flex items-center justify-center h-screen bg-gray-200">
    <div class="bg-white p-8 rounded-lg shadow-lg w-full max-w-md">
        <h1 class="text-2xl font-bold text-center mb-6">Sign in</h1>
        <form on:submit|preventDefault={registerUser}>
            <div class="mb-4">
                <label for="username" class="block text-gray-700 text-sm font-bold mb-2">Username</label>
                <div class="relative">
                    <i class="fas fa-user absolute left-3 top-2.5 text-gray-400"></i>
                    <input type="text" id="username" name="username" placeholder="Username" bind:value={username} required 
                           class="shadow appearance-none border rounded w-full py-2 px-10 text-gray-700 leading-tight focus:outline-none focus:ring focus:ring-blue-300">
                </div>
            </div>
            <div class="mb-4">
                <label for="password" class="block text-gray-700 text-sm font-bold mb-2">Password</label>
                <div class="relative">
                    <i class="fas fa-lock absolute left-3 top-2.5 text-gray-400"></i>
                    <input type="password" id="password" name="password" placeholder="Password" bind:value={password} required 
                           class="shadow appearance-none border rounded w-full py-2 px-10 text-gray-700 leading-tight focus:outline-none focus:ring focus:ring-blue-300">
                </div>
            </div>
            <div class="mb-6">
                <label for="role" class="block text-gray-700 text-sm font-bold mb-2">Role</label>
                <select id="role" name="role" class="shadow border rounded w-full py-2 px-3 text-gray-700 mb-3 leading-tight focus:outline-none focus:ring focus:ring-blue-300" bind:value={role}>
                    <option value="" disabled selected>Select a role</option>
                    {#each roles as { role, title }}
                        <option value={role}>{title}</option>
                    {/each}
                </select>
            </div>
            <button type="submit" class="w-full bg-blue-500 hover:bg-blue-600 text-white font-bold py-2 px-4 rounded focus:outline-none focus:ring focus:ring-blue-300">Sign in</button>
            <button type="button" on:click={() => window.location.href = '/login'} class="w-full bg-transparent border border-blue-500 hover:bg-blue-500 hover:text-white text-blue-500 font-bold py-2 px-4 rounded mt-2">Log in</button>
        </form>
    </div>
</main>
