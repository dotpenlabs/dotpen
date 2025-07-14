import PocketBase from 'pocketbase';
import { PUBLIC_API_ENDPOINT } from '$env/static/public';

export const pb = new PocketBase(PUBLIC_API_ENDPOINT);
