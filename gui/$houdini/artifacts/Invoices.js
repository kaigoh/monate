export default {
    "name": "Invoices",
    "kind": "HoudiniQuery",
    "hash": "0cea6fca4ec8d78027ec3bc75548bc201501d8153b817ee7ad8505f7ed449d9a",

    "raw": `query Invoices($storeId: ID!) {
  invoices(storeId: $storeId) {
    id
    description
    reference
    amountAtomic
    fiatAmount
    currency
    status
    complete
    coveredTotal
    moneroPayAddress
    createdAt
    resolvedAt
  }
}`,

    "rootType": "Query",
    "stripVariables": [],

    "selection": {
        "fields": {
            "invoices": {
                "type": "Invoice",
                "keyRaw": "invoices(storeId: $storeId)",

                "selection": {
                    "fields": {
                        "id": {
                            "type": "ID",
                            "keyRaw": "id",
                            "visible": true
                        },

                        "description": {
                            "type": "String",
                            "keyRaw": "description",
                            "nullable": true,
                            "visible": true
                        },

                        "reference": {
                            "type": "String",
                            "keyRaw": "reference",
                            "nullable": true,
                            "visible": true
                        },

                        "amountAtomic": {
                            "type": "Int64",
                            "keyRaw": "amountAtomic",
                            "visible": true
                        },

                        "fiatAmount": {
                            "type": "Float",
                            "keyRaw": "fiatAmount",
                            "visible": true
                        },

                        "currency": {
                            "type": "String",
                            "keyRaw": "currency",
                            "visible": true
                        },

                        "status": {
                            "type": "InvoiceStatus",
                            "keyRaw": "status",
                            "visible": true
                        },

                        "complete": {
                            "type": "Boolean",
                            "keyRaw": "complete",
                            "visible": true
                        },

                        "coveredTotal": {
                            "type": "Int64",
                            "keyRaw": "coveredTotal",
                            "visible": true
                        },

                        "moneroPayAddress": {
                            "type": "String",
                            "keyRaw": "moneroPayAddress",
                            "visible": true
                        },

                        "createdAt": {
                            "type": "DateTime",
                            "keyRaw": "createdAt",
                            "visible": true
                        },

                        "resolvedAt": {
                            "type": "DateTime",
                            "keyRaw": "resolvedAt",
                            "nullable": true,
                            "visible": true
                        }
                    }
                },

                "visible": true
            }
        }
    },

    "pluginData": {
        "houdini-svelte": {}
    },

    "input": {
        "fields": {
            "storeId": "ID"
        },

        "types": {},
        "defaults": {},
        "runtimeScalars": {}
    },

    "policy": "CacheOrNetwork",
    "partial": false
};

"HoudiniHash=963e78a8e3b295bb50f6dc39d74e31911dfdc7f01e8804489471c807440558c8";