import { CreateInvoiceStore } from "../plugins/houdini-svelte/stores/CreateInvoice";
import { StoresStore } from "../plugins/houdini-svelte/stores/Stores";
import { InvoicesStore } from "../plugins/houdini-svelte/stores/Invoices";
import type { Cache as InternalCache } from "./cache/cache";
import type { CacheTypeDef } from "./generated";
import { Cache } from "./public";
export * from "./client";
export * from "./lib";

export function graphql(
    str: "mutation CreateInvoice($input: CreateInvoiceInput!) {\n  createInvoice(input: $input) {\n    id\n    moneroPayAddress\n    status\n    amountAtomic\n    fiatAmount\n    currency\n    createdAt\n  }\n}\n"
): CreateInvoiceStore;

export function graphql(
    str: "query Stores {\n  stores {\n    id\n    name\n    slug\n    publicUrl\n    theme {\n      primaryColor\n      accentColor\n      headline\n      customCopy\n      showFiatAmount\n      showTicker\n      presetAmounts\n      backgroundUrl\n      logoUrl\n    }\n  }\n}\n"
): StoresStore;

export function graphql(
    str: "query Invoices($storeId: ID!) {\n  invoices(storeId: $storeId) {\n    id\n    description\n    reference\n    amountAtomic\n    fiatAmount\n    currency\n    status\n    complete\n    coveredTotal\n    moneroPayAddress\n    createdAt\n    resolvedAt\n  }\n}\n"
): InvoicesStore;

export declare function graphql<_Payload, _Result = _Payload>(str: TemplateStringsArray): _Result;
export declare const cache: Cache<CacheTypeDef>;
export declare function getCache(): InternalCache;