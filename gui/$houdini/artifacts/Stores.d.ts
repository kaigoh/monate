export type Stores = {
    readonly "input": Stores$input;
    readonly "result": Stores$result | undefined;
};

export type Stores$result = {
    readonly stores: ({
        readonly id: string;
        readonly name: string;
        readonly slug: string;
        readonly publicUrl: string;
        readonly theme: {
            readonly primaryColor: string;
            readonly accentColor: string;
            readonly headline: string | null;
            readonly customCopy: string | null;
            readonly showFiatAmount: boolean;
            readonly showTicker: boolean;
            readonly presetAmounts: (number)[];
            readonly backgroundUrl: string | null;
            readonly logoUrl: string | null;
        };
    })[];
};

export type Stores$input = null;

export type Stores$artifact = {
    "name": "Stores";
    "kind": "HoudiniQuery";
    "hash": "69fb3630673c7d3673c1b143327a00390ffb98fa8ed97139598b0ba281f43bac";
    "raw": `query Stores {
  stores {
    id
    name
    slug
    publicUrl
    theme {
      primaryColor
      accentColor
      headline
      customCopy
      showFiatAmount
      showTicker
      presetAmounts
      backgroundUrl
      logoUrl
    }
  }
}`;
    "rootType": "Query";
    "stripVariables": [];
    "selection": {
        "fields": {
            "stores": {
                "type": "Store";
                "keyRaw": "stores";
                "selection": {
                    "fields": {
                        "id": {
                            "type": "ID";
                            "keyRaw": "id";
                            "visible": true;
                        };
                        "name": {
                            "type": "String";
                            "keyRaw": "name";
                            "visible": true;
                        };
                        "slug": {
                            "type": "String";
                            "keyRaw": "slug";
                            "visible": true;
                        };
                        "publicUrl": {
                            "type": "String";
                            "keyRaw": "publicUrl";
                            "visible": true;
                        };
                        "theme": {
                            "type": "ThemeSettings";
                            "keyRaw": "theme";
                            "selection": {
                                "fields": {
                                    "primaryColor": {
                                        "type": "String";
                                        "keyRaw": "primaryColor";
                                        "visible": true;
                                    };
                                    "accentColor": {
                                        "type": "String";
                                        "keyRaw": "accentColor";
                                        "visible": true;
                                    };
                                    "headline": {
                                        "type": "String";
                                        "keyRaw": "headline";
                                        "nullable": true;
                                        "visible": true;
                                    };
                                    "customCopy": {
                                        "type": "String";
                                        "keyRaw": "customCopy";
                                        "nullable": true;
                                        "visible": true;
                                    };
                                    "showFiatAmount": {
                                        "type": "Boolean";
                                        "keyRaw": "showFiatAmount";
                                        "visible": true;
                                    };
                                    "showTicker": {
                                        "type": "Boolean";
                                        "keyRaw": "showTicker";
                                        "visible": true;
                                    };
                                    "presetAmounts": {
                                        "type": "Int";
                                        "keyRaw": "presetAmounts";
                                        "visible": true;
                                    };
                                    "backgroundUrl": {
                                        "type": "String";
                                        "keyRaw": "backgroundUrl";
                                        "nullable": true;
                                        "visible": true;
                                    };
                                    "logoUrl": {
                                        "type": "String";
                                        "keyRaw": "logoUrl";
                                        "nullable": true;
                                        "visible": true;
                                    };
                                };
                            };
                            "visible": true;
                        };
                    };
                };
                "visible": true;
            };
        };
    };
    "pluginData": {
        "houdini-svelte": {};
    };
    "policy": "CacheOrNetwork";
    "partial": false;
};