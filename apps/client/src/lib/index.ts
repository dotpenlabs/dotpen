import { PUBLIC_API_ENDPOINT } from '$env/static/public';

import PocketBase from 'pocketbase';
// import { PluginKit } from 'pluginkit';

export const pb = new PocketBase(PUBLIC_API_ENDPOINT);

// export const pk = new PluginKit({
// 	Database: 'https://raw.githubusercontent.com/dotpenlabs/repository/main/database.json',
// 	Plugins: 'https://raw.githubusercontent.com/dotpenlabs/repository/main/plugins/',
// 	Permissions: ['widget:sidebar:bottom', 'pluginkit:allow-same-origin', 'app:load'],
// 	ResumeAfterLoad: true,
// 	AllowInsecureSameOrigin: true
// });

const DB_NAME = 'dotpen-db';
const STORE_NAME = 'keyval';
const DB_VERSION = 1;

export class kv {
	static async open(): Promise<IDBDatabase> {
		return new Promise((resolve, reject) => {
			const request = indexedDB.open(DB_NAME, DB_VERSION);
			request.onupgradeneeded = () => {
				const db = request.result;
				if (!db.objectStoreNames.contains(STORE_NAME)) {
					db.createObjectStore(STORE_NAME);
				}
			};
			request.onsuccess = () => resolve(request.result);
			request.onerror = () => reject(request.error);
		});
	}

	static async get(key: string): Promise<string> {
		const db = await this.open();
		return new Promise((resolve, reject) => {
			const tx = db.transaction(STORE_NAME, 'readonly');
			const store = tx.objectStore(STORE_NAME);
			const req = store.get(key);
			req.onsuccess = () => resolve(req.result as string);
			req.onerror = () => reject(req.error);
		});
	}

	static async set(key: string, value: string): Promise<void> {
		const db = await this.open();
		return new Promise((resolve, reject) => {
			const tx = db.transaction(STORE_NAME, 'readwrite');
			const store = tx.objectStore(STORE_NAME);
			const req = store.put(value, key);
			req.onsuccess = () => resolve();
			req.onerror = () => reject(req.error);
		});
	}

	static async delete(key: string): Promise<void> {
		const db = await this.open();
		return new Promise((resolve, reject) => {
			const tx = db.transaction(STORE_NAME, 'readwrite');
			const store = tx.objectStore(STORE_NAME);
			const req = store.delete(key);
			req.onsuccess = () => resolve();
			req.onerror = () => reject(req.error);
		});
	}

	static async clear(): Promise<void> {
		const db = await this.open();
		return new Promise((resolve, reject) => {
			const tx = db.transaction(STORE_NAME, 'readwrite');
			const store = tx.objectStore(STORE_NAME);
			const req = store.clear();
			req.onsuccess = () => resolve();
			req.onerror = () => reject(req.error);
		});
	}
}
