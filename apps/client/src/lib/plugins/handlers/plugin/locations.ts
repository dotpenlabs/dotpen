import type { PluginConfig } from '../../head';
import { canExecute, permissions } from '../permissions';

export function locations(
	plugin: PluginConfig,
	data: {
		locations: string[];
	}
): void {
	canExecute(plugin, permissions.SET_LOCATION);

	plugin.locations = data.locations;

	console.info(
		`[PluginKit] Plugin "${plugin.id}" has been set to locations: ${data.locations.join(', ')}`
	);
}
