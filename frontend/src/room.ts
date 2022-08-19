import type { Credentials, Payload } from './wsclient';
import { WsClient } from './wsclient';

import type { Writable } from 'svelte/store';
import { writable } from 'svelte/store';

export enum GameState {
	NotConnected,
	WaitingRoom,
}

export class Room {
	private socket?: WsClient;
	private credentials?: Credentials;
	private gameState: Writable<GameState>;
	private roomId?: string;

	public constructor() {
		this.gameState = writable(GameState.NotConnected);
	}

	public join(roomId: string) {
		this.errorIfNotConnected();

		const joinEvent = {
			type: 'join',
			data: {
				roomId,
			},
		};

		return new Promise<void>((resolve, reject) => {
			const timeoutId = setTimeout(() => {
				reject();
			}, 5000);

			this.socket!.sendEvent(joinEvent, (payload: Payload) => {
				clearTimeout(timeoutId);
				resolve();
			});
		});
	}

	public GameState() {
		return {
			subscribe: this.gameState.subscribe,
		};
	}

	private setGameState(gameState: GameState) {
		this.gameState.set(gameState);
	}

	public connect(credentials: Credentials) {
		this.credentials = credentials;

		this.socket = new WsClient(this.credentials);
		return this.socket.connect();
	}

	private errorIfNotConnected() {
		if (this.socket === undefined) {
			throw new Error('the websocket is not connected');
		}
	}
}
