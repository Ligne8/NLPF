<script>
  import Navbar from '@components/Navbar.svelte';
  const userRole = localStorage.getItem('userRole');
  const isLogin = localStorage.getItem('isConnected');
  console.log(isLogin);
  if ( isLogin !== 'true' ) {
    window.location.href = '/login';
  }

</script>

<Navbar />