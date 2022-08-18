import { derived, writable } from 'svelte/store';
import { client } from './wsclient';

const writableHost = writable('');
const writablePlayers = writable<{ username: string; team: number }[]>([]);

client.on('playerJoin', (data) => {
	console.log(data);
	writableHost.set(data?.host);
	writablePlayers.set(data?.players);
});

client.on('infos', (data) => {
	console.log(data);
	writableHost.set(data?.host);
	writablePlayers.set(data?.players);
});

export const players = derived(writablePlayers, (s) => s);
export const host = derived(writableHost, (s) => s);
