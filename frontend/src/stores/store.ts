import { writable } from 'svelte/store';

// Types
export type UserRole = 'admin' | 'trafficManager' | 'trader' | 'client';

// Hardcoded data
export const userRole = writable<UserRole>('admin');
export const currentTab = writable<string>('');
export const currentTrafficManagerTab = writable<string>('Roads');

// Exemple data
export const tableData = writable([
    { name: 'Camion 1', status: 'ON_THE_WAY', currentCapacity: 120, totalCapacity: 120, location: 'Paris', road: ['Paris - Lyon', 'Paris - Montpellier', 'Paris - Marseille'] },
    { name: 'Camion 2', status: 'ON_THE_STOCK_EXCHANGE', currentCapacity: 38, totalCapacity: 154, location: 'Lyon', road: ['Lyon - Paris', 'Lyon - Montpellier'] },
    { name: 'Camion 3', status: 'AVAILABLE', currentCapacity: 52, totalCapacity: 86, location: 'Marseille', road: ['Marseille - Montpellier', 'Marseille - Lyon', 'Marseille - Marseille'] },
    { name: 'Camion 4', status: 'AVAILABLE', currentCapacity: 0, totalCapacity: 94, location: 'Montpellier', road: ['Montpellier - Marseille', 'Montpellier - Paris', 'Montpellier - Lyon', 'Montpellier - Perpignan'] },
]);