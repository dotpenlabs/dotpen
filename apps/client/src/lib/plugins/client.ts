// interaction code for user <-> PluginKit

import { toast } from 'svelte-sonner';
import { permissions } from './handlers/permissions';
import { setItem } from '../idb';
import { _PluginKit } from './head';

const database_head = 'https://raw.githubusercontent.com/dotpenlabs/PluginKit/main';

type PluginKitDatabase = {
	// global plugin id (e.g. nl.bijsven.dropbox)
	[key: string]: {
		manifest: string; // global plugin lookup (e.g. /plugins/nl.bijsven.dropbox/manifest.json)
	};
};

type PluginManifest = {
	id: string; // plugin id (e.g. nl.bijsven.dropbox)
	name: string; // plugin name (e.g. Dropbox)
	description: string; // plugin description (e.g. Dropbox integration)
	author: string; // plugin author (e.g. bijsven)

	hash: string; // plugin hash (e.g. sha256:1234567890abcdef)
	permissions: string[]; // plugin permissions (e.g. ["app:toast_info"])
	version: string; // plugin version (e.g. 1.0.0)
	app_version: string[]; // supported dotpen versions (NOT REQUIRED) (e.g. ["0.0.9", "0.1.0"])
};

class PluginKitClient {
	private stripBOM(str: string): string {
		if (str.charCodeAt(0) === 0xfeff) {
			return str.slice(1);
		}
		return str;
	}

	async download(url: string, hash: string): Promise<string> {
		return await fetch(url).then(async (resp) => {
			if (!resp.ok) {
				throw new Error(`[PluginKitClient] download failed: ${resp.statusText}`);
			}
			const text = await resp.text();
			const encoder = new TextEncoder();
			const data = encoder.encode(text);
			const hashBuffer = await window.crypto.subtle.digest('SHA-256', data);
			const hashArray = Array.from(new Uint8Array(hashBuffer));
			const hashHex = hashArray.map((b) => b.toString(16).padStart(2, '0')).join('');
			if (hashHex !== hash) {
				toast.error('Corrupted plugin', {
					description: 'The plugin you tried to install may be corrupted, malicious or damaged.'
				});
				throw new Error(`[PluginKitClient] download failed: hash mismatch`);
			}
			return text;
		});
	}

	async Register(pluginId: string): Promise<void> {
		const loading_toast = toast.loading('Installing plugin...', {
			description: 'This may take a few seconds...'
		});
		try {
			const resp = await fetch(database_head + '/database.json');
			if (!resp.ok) {
				throw new Error(`[PluginKitClient] database_head failed to load: ${resp.statusText}`);
			}

			const text = await resp.text();
			const database = JSON.parse(this.stripBOM(text).trim()) as PluginKitDatabase;

			if (!(pluginId in database)) {
				throw new Error(`[PluginKitClient] "${pluginId}" not found in PluginKit database.`);
			}

			const manifestUrl = database_head + database[pluginId].manifest;
			const manifestResp = await fetch(manifestUrl);
			if (!manifestResp.ok) {
				throw new Error(`[PluginKitClient] failed to fetch manifest: ${manifestResp.statusText}`);
			}

			const manifest = (await manifestResp.json()) as PluginManifest;
			console.log(`[PluginKitClient] "${pluginId}" loaded manifest:`, manifest);

			if (manifest.id !== pluginId) {
				toast.error('Installation failed', {
					description: 'The plugin you tried to install does not exist in our repository.'
				});
				throw new Error(
					`[PluginKitClient] manifest id "${manifest.id}" does not match plugin id "${pluginId}"`
				);
			}

			if (manifest.app_version) {
				if (APP_VERSION in manifest.app_version) {
					console.warn('[PluginKitClient] plugin is may be outdated, that can cause issues.');
					toast.warning(manifest.name, {
						description:
							'This plugin may be not compatible with the current version of Dotpen. Check in with the plugin author.'
					});
				} else {
					console.info(
						`[PluginKitClient] plugin "${pluginId}" is compatible with Dotpen ${APP_VERSION}`
					);
				}
			} else {
				console.warn(
					'[PluginKitClient] PluginKit does not know which version the plugin is compatible with. It is recommended to add the "app_version" property to the manifest to prevent issues.'
				);
			}

			if (!manifest.permissions) {
				console.warn(
					"[PluginKitClient] This plugin does not have any permissions assigned to it. You can't use any features that require interactions with Dotpen."
				);
			}

			if (manifest.permissions) {
				console.info(
					`[PluginKitClient] plugin "${pluginId}" has permissions:`,
					manifest.permissions
				);

				manifest.permissions
					.filter((permission) => !(permission in permissions))
					.forEach((permission) =>
						console.warn(
							`[PluginKitClient] Plugin "${pluginId}" has a permission "${permission}" that is not supported by PluginKit.`
						)
					);
			}

			console.info('[PluginKitClient] All checks passed, downloading plugin...');

			const pluginData = await this.download(
				database_head + database[pluginId].manifest.replace('/manifest.json', '') + '/plugin.html',
				manifest.hash
			);

			await setItem('plugin:' + pluginId, pluginData);
			console.info(
				`[PluginKitClient] Plugin "${pluginId}" downloaded and stored in "plugin:${pluginId}"`
			);

			_PluginKit.register({
				id: pluginId,
				name: manifest.name,
				content: pluginData,
				perms: manifest.permissions
			})
		} catch (error) {
			console.error('[PluginKitClient] error during plugin registration: ', error);
		} finally {
			toast.dismiss(loading_toast);
		}
	}
}

export const PluginKit = new PluginKitClient();
