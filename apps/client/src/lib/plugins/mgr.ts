// PluginMGR

type PluginCfg = {
	name: string;
	url: string;
	location: string;
};

class PluginMgr {
	private plugins: PluginCfg[] = [];

	register(plugin: PluginCfg) {
		this.plugins.push(plugin);
	}

	unregister(name: string) {
		this.plugins = this.plugins.filter((plugin) => plugin.name !== name);
	}

	getPlugins() {
		return this.plugins;
	}

	getPluginsAt(location: string) {
		return this.plugins.filter((plugin) => plugin.location === location);
	}
}

export const plgmgr = new PluginMgr();
