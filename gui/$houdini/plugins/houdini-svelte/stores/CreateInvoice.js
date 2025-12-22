import artifact from '$houdini/artifacts/CreateInvoice'
import { MutationStore } from '../runtime/stores/mutation'

export class CreateInvoiceStore extends MutationStore {
	constructor() {
		super({
			artifact,
		})
	}
}
