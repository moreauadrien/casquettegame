import type { Result } from 'ts-results';
import { Err, Ok } from 'ts-results';

const WEBSOCKET_URL = `wss://${window.location.host}/ws`;

export type EventHandler = (data: any) => void;

export class WsClient {
	private socket?: WebSocket;
	private eventHandlers: Map<string, EventHandler>;

	public constructor() {
		this.eventHandlers = new Map<string, EventHandler>();
	}

	public connect() {
		this.socket = new WebSocket(WEBSOCKET_URL);

		return new Promise<Result<null, Error>>((resolve) => {
			const onOpen = () => {
				this.socket!.removeEventListener('open', onOpen);
				this.socket!.removeEventListener('error', onError);

				this.setupHandlers();

				resolve(Ok(null));
			};

			const onError = () => {
				this.socket = undefined;
				resolve(Err(new Error('failed to connect')));
			};

			this.socket!.addEventListener('open', onOpen);
			this.socket!.addEventListener('error', onError);
		});
	}

	public on(eventName: string, handler: EventHandler) {
		const handlerAlreadyExist = this.eventHandlers.has(eventName);
		if (handlerAlreadyExist) {
			throw new Error('this event already has an handler');
		}

		this.eventHandlers.set(eventName, handler);
	}

	public sendMessage(message: string) {
		this.socket?.send(message);
	}

	private setupHandlers() {
		this.setupMessageHandler();
		this.setupErrorHandler();
		this.setupCloseHandler();
	}

	private setupMessageHandler() {
		this.socket!.addEventListener('message', (e) => {
			console.log(e);
			const handler = this.eventHandlers.get('stateUpdate');

			if (handler !== undefined) {
				handler(JSON.parse(e.data));
			}
		});
	}

	private setupErrorHandler() {
		this.socket!.addEventListener('error', console.log);
	}

	private setupCloseHandler() {
		this.socket!.addEventListener('close', (e) => {
			console.log(e);
			this.socket = undefined;
		});
	}
}
