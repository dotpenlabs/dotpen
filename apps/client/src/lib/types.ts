export type Bookmark = {
	id: string; // (collection: bookmarks) the id of the bookmark
	label: string; //  (crawled) the title of the page
	link: string; // (provided) the link to the bookmark
	favicon: string; // (url local, crawled) points to the favicon
	cover: string; // (url local, crawled) points to the cover
	collection: string; // (reference to collections) the collection of the bookmark
	deleted: boolean; // (provided) whether the bookmark is deleted
	updated: Date; // (server-gen) The last update to this record
	created: Date; // (server-gen) The first update to this record
	_favicon_base64?: string; // (local-gen) cache of favicon in base64
	_cover_base64?: string; // (local-gen) cache of cover in base64
};

export type Collection = {
	id: string; // (provided) the id of the collection
	name: string; // (provided) the name of the collection
};
