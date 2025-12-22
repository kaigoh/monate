export default {
    "name": "CreateInvoice",
    "kind": "HoudiniMutation",
    "hash": "eb64334b4727f6a957f6cf301c1e8c98c4bd639e1c5c458b69bd9ce204fda464",

    "raw": `mutation CreateInvoice($input: CreateInvoiceInput!) {
  createInvoice(input: $input) {
    id
    moneroPayAddress
    status
    amountAtomic
    fiatAmount
    currency
    createdAt
  }
}`,

    "rootType": "Mutation",
    "stripVariables": [],

    "selection": {
        "fields": {
            "createInvoice": {
                "type": "Invoice",
                "keyRaw": "createInvoice(input: $input)",

                "selection": {
                    "fields": {
                        "id": {
                            "type": "ID",
                            "keyRaw": "id",
                            "visible": true
                        },

                        "moneroPayAddress": {
                            "type": "String",
                            "keyRaw": "moneroPayAddress",
                            "visible": true
                        },

                        "status": {
                            "type": "InvoiceStatus",
                            "keyRaw": "status",
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

                        "createdAt": {
                            "type": "DateTime",
                            "keyRaw": "createdAt",
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
            "input": "CreateInvoiceInput"
        },

        "types": {
            "CreateInvoiceInput": {
                "storeId": "ID",
                "amountAtomic": "Int64",
                "fiatAmount": "Float",
                "currency": "String",
                "description": "String",
                "reference": "String"
            }
        },

        "defaults": {},
        "runtimeScalars": {}
    }
};

"HoudiniHash=c6c9dde030212785858dc0abf65dcf57c0b5fd0bc5104cc19c1e1bd4a4322b8a";