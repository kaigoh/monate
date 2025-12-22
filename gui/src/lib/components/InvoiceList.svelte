<script lang="ts">
  export let invoices: Array<{
    id: string;
    description?: string | null;
    reference?: string | null;
    amountAtomic: number;
    fiatAmount: number;
    currency: string;
    status: string;
    complete: boolean;
    coveredTotal: number;
    moneroPayAddress: string;
    createdAt: string;
    resolvedAt?: string | null;
  }> = [];

  const formatDate = (input: string | null | undefined) => {
    if (!input) return '—';
    return new Date(input).toLocaleString();
  };

  const formatAmount = (atomic: number) => {
    return (atomic / 1e12).toFixed(4);
  };
</script>

<table class="invoice-table">
  <thead>
    <tr>
      <th>ID</th>
      <th>Description</th>
      <th>Status</th>
      <th>Payment</th>
      <th>Address</th>
      <th>Created</th>
      <th>Resolved</th>
    </tr>
  </thead>
  <tbody>
    {#if invoices.length === 0}
      <tr>
        <td colspan="7" class="empty">No invoices yet</td>
      </tr>
    {:else}
      {#each invoices as invoice}
        <tr>
          <td>{invoice.id.slice(0, 6)}…</td>
          <td>{invoice.description || invoice.reference || '—'}</td>
          <td>
            <span class={`status status-${invoice.status.toLowerCase()}`}>
              {invoice.status}
            </span>
          </td>
          <td>
            {formatAmount(invoice.amountAtomic)} XMR
            <small>{invoice.fiatAmount.toFixed(2)} {invoice.currency}</small>
          </td>
          <td>
            <a href={`monero:${invoice.moneroPayAddress}`} target="_blank" rel="noreferrer">
              {invoice.moneroPayAddress.slice(0, 10)}…
            </a>
          </td>
          <td>{formatDate(invoice.createdAt)}</td>
          <td>{formatDate(invoice.resolvedAt || null)}</td>
        </tr>
      {/each}
    {/if}
  </tbody>
</table>

<style>
  .invoice-table {
    width: 100%;
    border-collapse: collapse;
    font-size: 0.95rem;
  }

  .invoice-table th,
  .invoice-table td {
    padding: 0.75rem;
    text-align: left;
    border-bottom: 1px solid rgba(10, 23, 49, 0.08);
  }

  .invoice-table th {
    font-weight: 600;
    color: rgba(9, 24, 41, 0.7);
    text-transform: uppercase;
    font-size: 0.78rem;
    letter-spacing: 0.08em;
  }

  .invoice-table .empty {
    text-align: center;
    padding: 1.5rem;
    color: rgba(9, 24, 41, 0.5);
  }

  .status {
    border-radius: 999px;
    padding: 0.25rem 0.75rem;
    font-size: 0.75rem;
    font-weight: 600;
  }

  .status-paid {
    background: rgba(5, 150, 105, 0.1);
    color: #047857;
  }

  .status-pending {
    background: rgba(249, 115, 22, 0.12);
    color: #c2410c;
  }

  .status-expired,
  .status-canceled {
    background: rgba(239, 68, 68, 0.12);
    color: #b91c1c;
  }
</style>
