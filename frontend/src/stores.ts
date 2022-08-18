import { browser } from '$app/env';
import { writable } from 'svelte/store';
import { generateId, generateToken } from './utils';

function createUsernameStore() {
	const { subscribe, set } = writable('');

	if (browser) {
		let value = localStorage.getItem('username');
		if (value !== null) {
			set(value);
		}

		subscribe((value) => {
			localStorage.setItem('username', value);
		});
	}

	return {
		subscribe,
		set,
	};
}

export const username = createUsernameStore();

function createTokenStore() {
	const { subscribe, set } = writable('');

	if (browser) {
		let value = localStorage.getItem('token');
		if (value === null) {
			const token = generateToken();
			localStorage.setItem('token', token);
			set(token);
		} else {
			set(value);
		}
	}

	return {
		subscribe,
	};
}

function createIdStore() {
	const { subscribe, set } = writable('');

	if (browser) {
		let value = localStorage.getItem('id');
		if (value === null) {
			const id = generateId();
			localStorage.setItem('id', id);
			set(id);
		} else {
			set(value);
		}
	}

	return {
		subscribe,
	};
}

export const token = createTokenStore();
export const playerId = createIdStore();
