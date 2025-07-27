import { toast } from 'svelte-sonner';
import type { PluginConfig } from '../../head';
import { canExecute, permissions } from '../permissions';

export function toast_info(plugin: PluginConfig, data: any): void {
	canExecute(plugin, permissions.APP_TOAST_INFO);

	toast.info(data.message, {
		description: 'Sent from ' + plugin.id
	});
}

export function toast_error(plugin: PluginConfig, data: any): void {
	canExecute(plugin, permissions.APP_TOAST_ERROR);

	toast.error(data.message, {
		description: 'Sent from ' + plugin.id
	});
}

export function toast_success(plugin: PluginConfig, data: any): void {
	canExecute(plugin, permissions.APP_TOAST_SUCCESS);

	toast.success(data.message, {
		description: 'Sent from ' + plugin.id
	});
}

export function toast_warning(plugin: PluginConfig, data: any): void {
	canExecute(plugin, permissions.APP_TOAST_WARNING);

	toast.warning(data.message, {
		description: 'Sent from ' + plugin.id
	});
}
