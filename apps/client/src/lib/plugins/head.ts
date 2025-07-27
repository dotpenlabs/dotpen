import type { HandlerFn } from './handlers';
import { handlers } from './handlers';

export type PluginConfig = {
	id: string;
	name: string;
	locations?: string[];
	content: string;
	render?: HTMLIFrameElement;
	perms: string[];
};

type Subscriber = () => void;

class PluginKitMgr {
	private plugins: PluginConfig[] = [];
	private handlers: Map<string, HandlerFn>;
	private subscribers: Set<Subscriber> = new Set();

	constructor() {
		this.handlers = new Map(Object.entries(handlers));
		window.addEventListener('message', this.#receive.bind(this));
	}

	subscribe(callback: Subscriber): () => void {
		this.subscribers.add(callback);
		return () => this.subscribers.delete(callback);
	}

	#notify(): void {
		for (const callback of this.subscribers) {
			callback();
		}
	}

	register(plugin: PluginConfig): void {
		if (this.plugins.some((p) => p.id === plugin.id)) {
			throw new Error(`[PluginKit] Plugin with id "${plugin.id}" is already registered.`);
		}
		this.plugins.push(plugin);
		this.#notify();
	}

	deregister(id: string): void {
		this.plugins = this.plugins.filter((plugin) => plugin.id !== id);
		this.#notify();
	}

	getByLocation(location: string): PluginConfig[] {
		return this.plugins.filter((plugin) => {
			const hasCorrectLocation = location in plugin.locations;
			const hasPermission = plugin.perms?.includes(`app:${location}`);

			if (hasCorrectLocation && !hasPermission) {
				console.error(
					`[PluginKit] Plugin "${plugin.id ?? 'unknown'}" tried to load on "${location}" without correct permissions ("app:${location}")`
				);
				return false;
			}

			return hasCorrectLocation && hasPermission;
		});
	}

	bind(id: string, iframe: HTMLIFrameElement): void {
		const plugin = this.plugins.find((p) => p.id === id);
		if (plugin) plugin.render = iframe;
	}

	emit(id: string, message: any): void {
		const plugin = this.plugins.find((p) => p.id === id);
		plugin?.render?.contentWindow?.postMessage(message, window.location.origin);
	}

	#receive(event: MessageEvent): void {
		console.debug('[PluginKit] Received message event:', event);
		for (const plugin of this.plugins) {
			if (plugin.render?.contentWindow === event.source) {
				console.debug('[PluginKit] Message accepted for plugin:', plugin.id);
				this.#handle(event.data, plugin);
			} else {
				console.debug('[PluginKit] Message ignored for plugin:', plugin.id);
			}
		}
	}

	#handle(data: any, plugin: PluginConfig): void {
		if (!data?.action) return;

		const handler = this.handlers.get(data.action);
		if (handler) {
			handler(plugin, data);
		} else {
			console.warn(`[PluginKit] No handler for action "${data.action}"`);
		}
	}
}

export const _PluginKit = (() => {
	if (typeof window !== 'undefined') {
		return new PluginKitMgr();
	} else {
		return {
			register: () => {},
			deregister: () => {},
			getByLocation: () => [],
			bind: () => {},
			emit: () => {},
			subscribe: () => () => {}
		};
	}
})();
