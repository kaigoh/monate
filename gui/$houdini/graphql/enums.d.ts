
type ValuesOf<T> = T[keyof T]
	
export declare const DedupeMatchMode: {
    readonly Variables: "Variables";
    readonly Operation: "Operation";
    readonly None: "None";
}

export type DedupeMatchMode$options = ValuesOf<typeof DedupeMatchMode>
 
export declare const InvoiceStatus: {
    readonly PENDING: "PENDING";
    readonly PAID: "PAID";
    readonly EXPIRED: "EXPIRED";
    readonly CANCELED: "CANCELED";
}

export type InvoiceStatus$options = ValuesOf<typeof InvoiceStatus>
 