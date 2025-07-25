import type { PluginConfig } from '../../head';

export function reload(plugin: PluginConfig, data: any): void {
    if (plugin.perms?.includes('app:reload')) {
        location.reload();
    }
}
