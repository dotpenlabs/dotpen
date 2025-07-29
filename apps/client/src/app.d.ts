/// <reference types="vite/client" />

declare module 'macy';

declare global {
	const APP_VERSION: string;

	interface Window {
		SetHydrating: (value: boolean) => void;
		SetColored: (isColored: 'default' | 'green' | 'red') => void;
	}

	namespace App {
		// interface Error {}
		// interface Locals {}
		// interface PageData {}
		// interface PageState {}
		// interface Platform {}
	}
}

export {};
