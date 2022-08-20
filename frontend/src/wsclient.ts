import { generateId } from './utils';

import Joi from 'joi';

interface Event {
	type: string;
	to?: string;
	data?: object;
}

export interface Payload {
	event: Event;
	messageId?: string;
}

const payloadSchema = Joi.object({
	event: {
		type: Joi.string().required(),
		to: Joi.string(),
		data: Joi.object(),
	},
	messageId: Joi.string(),
});

export type EventHandler = (data: any) => object | void;
export type ResponseHandler = ((data: any) => void) | (() => void);

const WEBSOCKET_URL = `ws://${window.location.host}/ws`;

export type Credentials = {
	username: string;
	token: string;
	id: string;
};

export class WsClient {
	private socket?: WebSocket;
	private credentials: Credentials;

	private eventHandlers: Map<string, EventHandler>;
	private responseHandlers: Map<string, ResponseHandler>;

	public constructor(credentials: Credentials) {
		this.credentials = credentials;
		this.eventHandlers = new Map<string, EventHandler>();
		this.responseHandlers = new Map<string, ResponseHandler>();
	}

	public connect() {
		this.socket = new WebSocket(WEBSOCKET_URL);

		return new Promise<string>((resolve, reject) => {
			const onOpen = () => {
				this.socket!.removeEventListener('open', onOpen);
				this.socket!.removeEventListener('error', onError);

				this.setupHandlers();

				resolve('connected');
			};

			const onError = () => {
				this.socket = undefined;
				reject('failed to connect');
			};

			this.socket!.addEventListener('open', onOpen);
			this.socket!.addEventListener('error', onError);
		});
	}

	public sendEvent(event: Event, handler?: ResponseHandler) {
		if (this.socket === undefined) throw new Error('websocket is not connected');

		let payload;

		if (event.type === 'response') {
			payload = {
				event,
			};
		} else {
			const messageId = generateId();
			payload = {
				messageId,
				credentials: this.credentials,
				event,
			};

			if (handler !== undefined) this.onResponse(messageId, handler);
		}

		this.socket.send(JSON.stringify(payload));
	}

	public on(eventName: string, handler: EventHandler) {
		const handlerAlreadyExist = this.eventHandlers.has(eventName);
		if (handlerAlreadyExist) {
			throw new Error('this event already has an handler');
		}

		this.eventHandlers.set(eventName, handler);
	}

	private onResponse(responseId: string, handler: ResponseHandler) {
		const handlerAlreadyExist = this.responseHandlers.has(responseId);
		if (handlerAlreadyExist) {
			throw new Error('this event already has an handler');
		}

		this.responseHandlers.set(responseId, handler);

		setTimeout(() => {
			this.responseHandlers.delete(responseId);
		}, 5000);
	}

	private setupHandlers() {
		this.setupMessageHandler();
		this.setupErrorHandler();
		this.setupCloseHandler();
	}

	private runEventHandler(payload: Payload) {
		if (payload.messageId === undefined) {
			throw new Error('"messageId" field is required');
		}

		const handler = this.eventHandlers.get(payload.event.type);

		if (handler !== undefined) {
			const responseData = handler(payload.event.data);

			if (responseData !== undefined) {
				const responseEvent = {
					type: 'response',
					to: payload.messageId,
					data: responseData,
				};

				this.sendEvent(responseEvent);
			}
		}
	}

	private runResponseHandler(payload: Payload) {
		if (payload.event.to === undefined) {
			throw new Error('"to" field is required on a response event');
		}

		const handler = this.responseHandlers.get(payload.event.to!);
		this.responseHandlers.delete(payload.event.to!);

		if (handler !== undefined) {
			handler(payload.event.data);
		}
	}

	private setupMessageHandler() {
		this.socket!.addEventListener('message', (e) => {
			const rawPayload = JSON.parse(e.data);

			const value = payloadSchema.validate(rawPayload);
			if (value.error !== undefined) {
				throw new Error(value.error.message);
			}

			const payload = <Payload>rawPayload;

			if (payload.event.type === 'response') {
				this.runResponseHandler(payload);
			} else {
				this.runEventHandler(payload);
			}
		});
	}

	private setupErrorHandler() {
		this.socket!.addEventListener('error', (e) => {});
	}

	private setupCloseHandler() {
		this.socket!.addEventListener('close', (e) => {
			this.socket = undefined;
		});
	}
}
