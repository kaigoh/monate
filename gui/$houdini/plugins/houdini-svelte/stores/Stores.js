import { QueryStore } from '../runtime/stores/query'
import artifact from '$houdini/artifacts/Stores'
import { initClient } from '$houdini/plugins/houdini-svelte/runtime/client'

export class StoresStore extends QueryStore {
	constructor() {
		super({
			artifact,
			storeName: "StoresStore",
			variables: false,
		})
	}
}

export async function load_Stores(params) {
  await initClient()

	const store = new StoresStore()

	await store.fetch(params)

	return {
		Stores: store,
	}
}
