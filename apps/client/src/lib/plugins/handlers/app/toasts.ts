import { toast } from 'svelte-sonner';
import type { PluginConfig } from '../../head';

export function toast_info(plugin: PluginConfig, data: any): void {
	toast.info(data.message, {
		description: 'Sent from ' + plugin.id
	});
}

export function toast_error(plugin: PluginConfig, data: any): void {
	toast.error(data.message, {
		description: 'Sent from ' + plugin.id
	});
}

export function toast_success(plugin: PluginConfig, data: any): void {
	toast.success(data.message, {
		description: 'Sent from ' + plugin.id
	});
}

export function toast_warning(plugin: PluginConfig, data: any): void {
	toast.warning(data.message, {
		description: 'Sent from ' + plugin.id
	});
}
