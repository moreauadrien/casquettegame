import type { Result } from 'ts-results';
import { Ok, Err } from 'ts-results';
import type { Team } from './utils';

export type PlayerInfos = {
	username: string;
	id: string;
	team: Team;
};

export async function createRoom(hostUsername: string): Promise<Result<{ roomId: string }, Error>> {
	const init = {
		method: 'POST',
		headers: {
			'Content-Type': 'application/json',
		},
		body: JSON.stringify({
			username: hostUsername,
		}),
	};

	try {
		const resp = await fetch('/api/create', init);
		const json = await resp.json();

		if (json.status !== 'ok') {
			return Err(new Error(json.message ?? 'an error has occurred'));
		}

		return Ok({ roomId: json.roomId });
	} catch (e) {
		console.error(e);
		return Err(new Error(''));
	}
}

export async function joinRoom(username: string, roomId: string): Promise<Result<null, Error>> {
	const init = {
		method: 'POST',
		headers: {
			'Content-Type': 'application/json',
		},
		body: JSON.stringify({
			username: username,
			roomId: roomId,
		}),
	};

	try {
		const resp = await fetch('/api/join', init);
		const json = await resp.json();

		if (json.status !== 'ok') {
			return Err(new Error(json.message ?? 'an error has occurred'));
		}

		return Ok(null);
	} catch (e) {
		console.error(e);
		return Err(new Error(''));
	}
}
