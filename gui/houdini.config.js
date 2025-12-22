/// <references types="houdini-svelte">

/** @type {import('houdini').ConfigFile} */
const config = {
	runtimeDir: '$houdini',
	schemaPath: '../graph/schema.graphqls',
	apiUrl: '/query',
	scalars: {
		DateTime: {
			type: 'string',
			marshal(val) {
				return val;
			},
			unmarshal(val) {
				return val;
			}
		},
		Int64: {
			type: 'number',
			marshal(val) {
				return val;
			},
			unmarshal(val) {
				return val;
			}
		}
	},
	plugins: {
		'houdini-svelte': {
			client: './src/lib/graphql/client'
		}
	}
};

export default config;
