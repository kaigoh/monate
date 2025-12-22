// @ts-nocheck
import type { PageLoad } from './$types';
import { load_Stores, load_Invoices } from '$houdini';

export const prerender = false;

export const load = async (event: Parameters<PageLoad>[0]) => {
  const { Stores } = await load_Stores({ event });
  const storeList = Stores.data?.stores ?? [];

  let initialInvoices = [];
  const firstStoreId = storeList[0]?.id ?? null;
  if (firstStoreId) {
    const { Invoices } = await load_Invoices({
      event,
      variables: { storeId: firstStoreId }
    });
    initialInvoices = Invoices.data?.invoices ?? [];
  }

  return {
    stores: storeList,
    initialInvoices,
    initialStoreId: firstStoreId
  };
};
