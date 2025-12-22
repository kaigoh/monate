import { QueryStore } from '../runtime/stores/query'
import artifact from '$houdini/artifacts/Invoices'
import { initClient } from '$houdini/plugins/houdini-svelte/runtime/client'

export class InvoicesStore extends QueryStore {
	constructor() {
		super({
			artifact,
			storeName: "InvoicesStore",
			variables: true,
		})
	}
}

export async function load_Invoices(params) {
  await initClient()

	const store = new InvoicesStore()

	await store.fetch(params)

	return {
		Invoices: store,
	}
}
