import { writable } from 'svelte/store';

// Types
export type UserRole = 'admin' | 'trafficManager' | 'trader' | 'client';

// Store data into local storage
const storedRole = localStorage.getItem('userRole') as UserRole || 'admin';
const storedId = localStorage.getItem('userId') || '0';

// Store variables into localStorage
export const userRole = writable<UserRole>(storedRole);
export const userId = writable<string>(storedId);

if (storedId === '0' && window.location.pathname !== '/register' && window.location.pathname !== '/login') {
    window.location.href = "/login";
}

// Update local storage on change event
userRole.subscribe(value => {
    localStorage.setItem('userRole', value);
});
userId.subscribe(value => {
    localStorage.setItem('userId', value);
});

export const currentTab = writable<string>('');
export const currentTrafficManagerTab = writable<string>('');