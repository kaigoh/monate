import type { Record } from "./public/record";
import { Stores$result, Stores$input } from "$houdini/artifacts/Stores";
import { StoresStore } from "../plugins/houdini-svelte/stores/Stores";
import { Invoices$result, Invoices$input } from "$houdini/artifacts/Invoices";
import { InvoicesStore } from "../plugins/houdini-svelte/stores/Invoices";

export declare type CacheTypeDef = {
    types: {
        ThemeSettings: {
            idFields: never;
            fields: {
                primaryColor: {
                    type: string;
                    args: never;
                };
                accentColor: {
                    type: string;
                    args: never;
                };
                backgroundUrl: {
                    type: string | null;
                    args: never;
                };
                logoUrl: {
                    type: string | null;
                    args: never;
                };
                headline: {
                    type: string | null;
                    args: never;
                };
                customCopy: {
                    type: string | null;
                    args: never;
                };
                showFiatAmount: {
                    type: boolean;
                    args: never;
                };
                showTicker: {
                    type: boolean;
                    args: never;
                };
                presetAmounts: {
                    type: (number)[];
                    args: never;
                };
            };
            fragments: [];
        };
        Store: {
            idFields: {
                id: string;
            };
            fields: {
                id: {
                    type: string;
                    args: never;
                };
                name: {
                    type: string;
                    args: never;
                };
                slug: {
                    type: string;
                    args: never;
                };
                publicUrl: {
                    type: string;
                    args: never;
                };
                theme: {
                    type: Record<CacheTypeDef, "ThemeSettings">;
                    args: never;
                };
                createdAt: {
                    type: string;
                    args: never;
                };
                updatedAt: {
                    type: string;
                    args: never;
                };
            };
            fragments: [];
        };
        Invoice: {
            idFields: {
                id: string;
            };
            fields: {
                id: {
                    type: string;
                    args: never;
                };
                storeId: {
                    type: string;
                    args: never;
                };
                store: {
                    type: Record<CacheTypeDef, "Store">;
                    args: never;
                };
                description: {
                    type: string | null;
                    args: never;
                };
                reference: {
                    type: string | null;
                    args: never;
                };
                amountAtomic: {
                    type: number;
                    args: never;
                };
                expectedAmount: {
                    type: number;
                    args: never;
                };
                currency: {
                    type: string;
                    args: never;
                };
                fiatAmount: {
                    type: number;
                    args: never;
                };
                moneroPayAddress: {
                    type: string;
                    args: never;
                };
                status: {
                    type: InvoiceStatus;
                    args: never;
                };
                complete: {
                    type: boolean;
                    args: never;
                };
                coveredTotal: {
                    type: number;
                    args: never;
                };
                coveredUnlocked: {
                    type: number;
                    args: never;
                };
                manualCheckCount: {
                    type: number;
                    args: never;
                };
                createdAt: {
                    type: string;
                    args: never;
                };
                updatedAt: {
                    type: string;
                    args: never;
                };
                resolvedAt: {
                    type: string | null;
                    args: never;
                };
            };
            fragments: [];
        };
        __ROOT__: {
            idFields: {};
            fields: {
                stores: {
                    type: (Record<CacheTypeDef, "Store">)[];
                    args: never;
                };
                store: {
                    type: Record<CacheTypeDef, "Store"> | null;
                    args: {
                        id: string | number;
                    };
                };
                invoices: {
                    type: (Record<CacheTypeDef, "Invoice">)[];
                    args: {
                        storeId: string | number;
                    };
                };
                invoice: {
                    type: Record<CacheTypeDef, "Invoice"> | null;
                    args: {
                        id: string | number;
                    };
                };
            };
            fragments: [];
        };
    };
    lists: {};
    queries: [[InvoicesStore, Invoices$result, Invoices$input], [StoresStore, Stores$result, Stores$input]];
};