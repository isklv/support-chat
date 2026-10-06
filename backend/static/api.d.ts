import type { SupportView } from './types';
export declare function getStoredToken(): string | null;
export declare function saveStoredToken(token: string): void;
export declare function getStoredContact(): {
    email: string;
    name: string;
};
export declare function saveStoredContact(email: string, name: string): void;
export declare class SupportClient {
    private baseUrl;
    constructor(apiUrl: string);
    private headers;
    fetchView(token?: string | null, userId?: string): Promise<SupportView>;
    sendMessage(form: FormData, token?: string | null, userId?: string): Promise<SupportView>;
    markSeen(token?: string | null, userId?: string): Promise<void>;
    checkUnread(token?: string | null, userId?: string): Promise<number>;
    attachmentUrl(messageId: string, token?: string | null): string;
}
