import type { HandlerFn } from '../index';
import { reload } from './reload';
import { toast_info, toast_error, toast_success, toast_warning } from './toasts';

export const handlers_app: Record<string, HandlerFn> = {
	reload,
	toast_info,
	toast_error,
	toast_success,
	toast_warning
};
