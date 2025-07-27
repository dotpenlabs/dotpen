import type { PluginConfig } from '../head';

export const permissions = {
	APP_RELOAD: 'app:reload',
	APP_TOAST_INFO: 'app:toast_info',
	APP_TOAST_ERROR: 'app:toast_error',
	APP_TOAST_SUCCESS: 'app:toast_success',
	APP_TOAST_WARNING: 'app:toast_warning',
	BOOKMARKS_ADD: 'bookmarks:add',
	BOOKMARKS_REMOVE: 'bookmarks:remove',
	BOOKMARKS_UPDATE: 'bookmarks:update',
	BOOKMARKS_GET: 'bookmarks:get',
	BOOKMARKS_GET_ALL: 'bookmarks:get_all',
	COLLECTIONS_GET: 'collections:get',
	COLLECTIONS_GET_ALL: 'collections:get_all',
	COLLECTIONS_CREATE: 'collections:create',
	COLLECTIONS_UPDATE: 'collections:update',
	COLLECTIONS_DELETE: 'collections:delete',
	AI_TAGS: 'ai:tags',
	AI_SEARCH: 'ai:search',
	AI_QUERY_ALL_TAGS: 'ai:query_all_tags',
	AI_QUERY_TAGS_IN_OPEN_COLLECTIONS: 'ai:query_tags_in_open_collections',
	SET_LOCATION: 'plugin:set_visible_locations'
};

export function canExecute(plugin: PluginConfig, permission: string) {
	if (!plugin.perms?.includes(permission)) {
		throw new Error(
			`[PluginKit] Plugin "${plugin.id}" tried to execute a permission it does not have ("${permission}")`
		);
	}
}
