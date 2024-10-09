<script>
    import axios from "axios";
    import { userId, userRole } from "@stores/store";

    const API_BASE_URL = import.meta.env.VITE_API_BASE_URL;

    let username = '';
    let password = '';
    let errorMessage = '';

    async function LoginUser() {
        errorMessage = '';
        try {
            const response = await axios.post(`${API_BASE_URL}/auth/login`, {
                username,
                password,
            });
            if (response.status === 200) {
                userId.set(response.data.user.id);
                userRole.set(response.data.user.role);
                window.location.href = '/';
            }
        } catch (error) {
            errorMessage = 'Incorrect username or password.';
            console.error(error);
        }
    }
</script>

<main class="flex items-center justify-center h-screen bg-gray-200">
    <div class="bg-white p-8 rounded-lg shadow-lg w-full max-w-md">
        <h1 class="text-2xl font-bold text-center mb-6">Log in</h1>
        <form on:submit|preventDefault={LoginUser}>
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
            {#if errorMessage}
                <div class="text-red-500 text-sm text-center mb-4">{errorMessage}</div>
            {/if}
            <button type="submit" class="w-full bg-blue-500 hover:bg-blue-600 text-white font-bold py-2 px-4 rounded focus:outline-none focus:ring focus:ring-blue-300">Login</button>
            <button type="button" on:click={() => window.location.href = '/register'} class="w-full bg-transparent border border-blue-500 hover:bg-blue-500 hover:text-white text-blue-500 font-bold py-2 px-4 rounded mt-2">Sign in</button>
        </form>
    </div>
</main>