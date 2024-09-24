import { writable } from 'svelte/store';

// Types
export type UserRole = 'admin' | 'trafficManager' | 'trader' | 'client';

// Hardcoded data
export const userRole = writable<UserRole>('admin');
export const currentTab = writable<string>('');
export const currentTrafficManagerTab = writable<string>('');
