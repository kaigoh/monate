<script lang="ts">
  import type { PageData } from './$types';
  import StoreCard from '$lib/components/StoreCard.svelte';
  import InvoiceList from '$lib/components/InvoiceList.svelte';
  import { CheckCircle, RefreshCcw } from '@lucide/svelte';

  export let data: PageData;

  let stores = data.stores ?? [];
  let invoices = data.initialInvoices ?? [];
  let selectedStoreId: string | null = data.initialStoreId ?? null;
  let amountXmr = 0.1;
  let fiatAmount = 15;
  let currency = 'USD';
  let description = 'Coffee and snacks';
  let reference = '';
  let isSubmitting = false;
  let statusMessage = '';

  const formatAtomic = (xmr: number) => Math.round(xmr * 1_000_000_000_000);

  const fetchInvoices = async (storeId: string) => {
    const query = `#graphql\n      query Invoices($storeId: ID!) {\n        invoices(storeId: $storeId) {\n          id\n          description\n          reference\n          amountAtomic\n          fiatAmount\n          currency\n          status\n          complete\n          moneroPayAddress\n          createdAt\n          resolvedAt\n        }\n      }\n    `;
    const res = await fetch('/query', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ query, variables: { storeId } }),
    });
    const json = await res.json();
    if (json.errors) {
      throw new Error(json.errors[0]?.message ?? 'Failed to load invoices');
    }
    return json.data.invoices ?? [];
  };

  const handleStoreChange = async (storeId: string) => {
    selectedStoreId = storeId;
    statusMessage = '';
    invoices = await fetchInvoices(storeId);
  };

  const createInvoice = async () => {
    if (!selectedStoreId) return;
    isSubmitting = true;
    statusMessage = '';
    try {
      const mutation = `#graphql\n        mutation CreateInvoice($input: CreateInvoiceInput!) {\n          createInvoice(input: $input) {\n            id\n          }\n        }\n      `;
      const payload = {
        storeId: selectedStoreId,
        amountAtomic: formatAtomic(amountXmr),
        fiatAmount: fiatAmount,
        currency,
        description,
        reference,
      };
      const res = await fetch('/query', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          query: mutation,
          variables: { input: payload },
        }),
      });
      const json = await res.json();
      if (json.errors) {
        throw new Error(json.errors[0]?.message ?? 'Failed to create invoice');
      }
      statusMessage = `Invoice ${json.data.createInvoice.id} created`;
      invoices = await fetchInvoices(selectedStoreId);
    } catch (err) {
      statusMessage =
        err instanceof Error ? err.message : 'Unable to create invoice';
    } finally {
      isSubmitting = false;
    }
  };
</script>

<svelte:head>
  <title>Monate Dashboard</title>
</svelte:head>

<main class="page">
  <section class="hero">
    <div>
      <p class="eyebrow">Monate POS</p>
      <h1>Designable Monero checkouts for your community</h1>
      <p>
        Drive donations, manage campaigns, and follow MoneroPay invoices without
        leaving this dashboard.
      </p>
    </div>
  </section>

  <section class="grid">
    <div class="panel">
      <header>
        <h2>Your stores</h2>
        <p>Select a store to view invoices and create new payment links.</p>
      </header>
      <div class="store-grid">
        {#each stores as store}
          <StoreCard
            {store}
            active={store.id === selectedStoreId}
            on:click={() => handleStoreChange(store.id)}
          />
        {/each}
      </div>
    </div>

    <div class="panel">
      <header>
        <h2>Create invoice</h2>
        <p>Use theme presets or enter a custom amount/ref below.</p>
      </header>
      {#if !selectedStoreId}
        <p class="empty">Select a store to continue.</p>
      {:else}
        <form class="form" on:submit|preventDefault={createInvoice}>
          <label>
            Amount (XMR)
            <input
              type="number"
              step="0.01"
              min="0"
              bind:value={amountXmr}
              required
            />
          </label>
          <label>
            Fiat Amount
            <input
              type="number"
              step="0.01"
              min="0"
              bind:value={fiatAmount}
              required
            />
          </label>
          <label>
            Currency Code
            <input type="text" bind:value={currency} maxlength="8" required />
          </label>
          <label>
            Description
            <input
              type="text"
              bind:value={description}
              placeholder="What is this for?"
            />
          </label>
          <label>
            Reference
            <input
              type="text"
              bind:value={reference}
              placeholder="Optional internal note"
            />
          </label>
          <button type="submit" disabled={isSubmitting}>
            {#if isSubmitting}
              <span class="spinner-icon">
                <RefreshCcw />
              </span>
              Creating…
            {:else}
              <CheckCircle />
              Create invoice
            {/if}
          </button>
          {#if statusMessage}
            <p class="status-message">{statusMessage}</p>
          {/if}
        </form>
      {/if}
    </div>
  </section>

  <section class="panel">
    <header>
      <h2>Recent invoices</h2>
      <p>Live data pulled directly from the GraphQL API.</p>
    </header>
    <InvoiceList {invoices} />
  </section>
</main>

<style>
  .page {
    max-width: 1100px;
    margin: 0 auto;
    padding: 3rem 1.5rem 4rem;
    display: flex;
    flex-direction: column;
    gap: 2.5rem;
  }

  .hero {
    background: #fff;
    border-radius: 28px;
    padding: 2.5rem;
    box-shadow: 0 20px 60px rgba(10, 23, 49, 0.12);
  }

  .hero h1 {
    margin: 0.4rem 0 0.5rem;
    font-size: clamp(2.25rem, 4vw, 3rem);
  }

  .eyebrow {
    text-transform: uppercase;
    letter-spacing: 0.22em;
    color: rgba(9, 24, 41, 0.6);
    font-size: 0.85rem;
    margin: 0;
  }

  .grid {
    display: grid;
    gap: 1.5rem;
    grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  }

  .panel {
    background: #fff;
    border-radius: 24px;
    padding: 2rem;
    box-shadow: 0 20px 60px rgba(10, 23, 49, 0.08);
    display: flex;
    flex-direction: column;
    gap: 1.25rem;
  }

  .panel header h2 {
    margin: 0;
  }

  .store-grid {
    display: grid;
    gap: 1rem;
    grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  }

  .form {
    display: grid;
    gap: 1rem;
  }

  label {
    font-size: 0.9rem;
    display: flex;
    flex-direction: column;
    gap: 0.3rem;
    color: rgba(9, 24, 41, 0.72);
  }

  input {
    border-radius: 12px;
    border: 1px solid rgba(9, 24, 41, 0.15);
    padding: 0.75rem 1rem;
    font-size: 1rem;
  }

  button[type='submit'] {
    border: none;
    border-radius: 16px;
    background: linear-gradient(120deg, #ff7a18, #af002d 85%);
    color: #fff;
    font-size: 1rem;
    padding: 0.85rem 1rem;
    display: inline-flex;
    align-items: center;
    gap: 0.5rem;
  }

  button[disabled] {
    opacity: 0.6;
  }

  .spinner-icon {
    display: inline-flex;
    animation: spin 1.2s linear infinite;
  }

  @keyframes spin {
    from {
      transform: rotate(0deg);
    }
    to {
      transform: rotate(360deg);
    }
  }

  .status-message {
    margin: 0;
    color: #0f766e;
  }

  .empty {
    color: rgba(9, 24, 41, 0.5);
  }
</style>
