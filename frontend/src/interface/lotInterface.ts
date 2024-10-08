import type { Tractor } from "./tractorInterface";

export type Lot = {
  id: string;
  resource_type: string;
  volume: number;

  start_checkpoint_id?: string;
  start_checkpoint?: any; // TODO: add interface checkpoint

  end_checkpoint_id?: string;
  end_checkpoint?: any;

  tractor_id?: string;
  tractor?: Tractor;

  created_at: string;

  current_checkpoint_id?: string;
  current_checkpoint?: any;

  offer?: any;

  owner_id: string;
  owner?: any; // TODO: add interface User

  state: string;

  max_price_by_km: number;

  traffic_manager_id?: string;
  traffic_manager?: string;

  trader_id?: string;
  trader?: string;
};
