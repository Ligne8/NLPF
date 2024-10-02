export type Tractor = {
    id: string;
    resource_type: string;
    max_units: number;
    current_units: number;

    current_checkpoint_id?: string;
    // TODO: add current_checkpoint
    current_checkpoint?: string;

    state: string;
    created_at: string;
    owner_id?: string;

    // TODO: add owner
    owner?: string;
    min_price_by_km: number;

    traffic_manager_id?: string;
    // TODO: add traffic_manager
    traffic_manager?: string;

    trader_id?: string;
    // TODO: add trader
    trader?: string;

    route_id?: string;
    // TODO: add route
    route?: string;
};
