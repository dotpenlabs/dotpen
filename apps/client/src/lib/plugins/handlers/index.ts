import { handlers_app } from './app';
import { handlers_bookmarks } from './bookmarks';

export type HandlerFn = (plugin: any, data: any) => void;

export const handlers: Record<string, HandlerFn> = {
	...handlers_app,
	...handlers_bookmarks
};
