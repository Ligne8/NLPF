<script>
    import { Router, Link, Route } from "svelte-routing";
    import Register from "@pages/Authentification/Register.svelte";
    import Login from "@pages/Authentification/Login.svelte";
    import Home from '@pages/Home.svelte';
    import Lots from '@pages/lots/lots.svelte';
    import Tractors from '@pages/tractors/tractors.svelte';
    import TrafficManager from '@pages/traffic_manager/traffic_manager.svelte';
    import TrafficManagerRoutes from "@pages/traffic_manager/routes/routes.svelte";
    import TrafficManagerLots from "@pages/traffic_manager/lots/lots.svelte";
    import TrafficManagerTractors from "@pages/traffic_manager/tractors/tractors.svelte";
    import Trader from '@pages/trader/trader.svelte';
    import StockExchange from '@pages/stock_exchange/stock_exchange.svelte';
    import NotFound from '@pages/404.svelte'
    import './app.css';
    import '@fortawesome/fontawesome-free/css/all.css';
    export let url = "";
</script>

<Router {url} >
    <div>
        <Route path="/" component={Home} />
        <Route path="/register" component={Register} />
        <Route path="/login" component={Login} />
        <Route path="/lots" component={Lots} />
        <Route path="/tractors" component={Tractors} />
        <Route path="/traffic-manager" component={TrafficManager} />
        <Route path="/traffic-manager/routes" component={TrafficManagerRoutes} />
        <Route path="/traffic-manager/lots" component={TrafficManagerLots} />
        <Route path="/traffic-manager/tractors" component={TrafficManagerTractors} />
        <Route path="/trader" component={Trader} />
        <Route path="/stock-exchange" component={StockExchange} />
    </div>
</Router>