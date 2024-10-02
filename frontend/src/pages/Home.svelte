<script lang="ts">
  import { userId, userRole } from '../stores/store';
  import { faUser, faTruck, faChartLine, faLock } from '@fortawesome/free-solid-svg-icons';
  import { FontAwesomeIcon } from '@fortawesome/svelte-fontawesome';

  // Variables
  let title: string = 'Accueil';
  let subtitle: string = 'Sélectionnez un rôle pour accéder à votre espace.';

  // Roles data
  const roles = [
    { role: 'client', title: 'Client', icon: faUser, path: '/lots' },
    { role: 'trafficManager', title: 'Traffic Manager', icon: faTruck, path: '/traffic-manager' },
    { role: 'trader', title: 'Trader', icon: faChartLine, path: '/trader' },
    { role: 'admin', title: 'Admin', icon: faLock, path: '/lots' }
  ];

  // Function to handle role click
  const handleRoleClick = async (role: string, path: string) => {
    await userRole.set(role as any);
    await userId.set('1');
    window.location.href = path;
  };
</script>


<!-- svelte-ignore a11y-click-events-have-key-events -->
<!-- svelte-ignore a11y-no-static-element-interactions -->

<main class="px-10 flex flex-col items-center justify-center min-h-screen">
  
  <!-- Title and subtitle -->
  <section class="text-center">
    <h1 class="text-6xl font-bold mb-4">{title}</h1>
    <h2 class="text-3xl mb-10 text-gray-600">{subtitle}</h2>
  </section>

  <!-- Blocks -->
  <div class="container mx-auto">
    <div class="flex flex-wrap justify-between items-center h-full gap-6" style="color: rgba(255, 255, 255, 0.9);">
      {#each roles as { role, title, icon, path }, i}
        <div
          class="flex flex-col items-center px-6 py-20 bg-white rounded-lg shadow-md cursor-pointer hover:shadow-lg transition-transform transform hover:-translate-y-2 w-full md:w-[23%] text-center"
          style="background-color: hsl(210, 100%, {60 - i * 15}%);"
          on:click={() => handleRoleClick(role, path)}
        >
          <div class="text-8xl mb-4">
            <FontAwesomeIcon icon={icon} />
          </div>
          <div class="text-2xl font-bold">{title}</div>
        </div>
      {/each}
    </div>
  </div>
</main>