import { v4 as uuidv4 } from 'uuid';

export enum Team {
	BLUE,
	PURPLE,
	RED,
	YELLOW,
}

export function generateToken() {
	const preShuffled = uuidv4().replaceAll('-', '');
	return preShuffled
		.split('')
		.sort(function () {
			return 0.5 - Math.random();
		})
		.join('');
}

export function generateId() {
	return uuidv4();
}
