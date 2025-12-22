<script lang="ts">
  export let store: {
    id: string;
    name: string;
    slug: string;
    publicUrl: string;
    theme: {
      primaryColor: string;
      accentColor: string;
      headline?: string | null;
      customCopy?: string | null;
      presetAmounts: number[];
      showFiatAmount: boolean;
      showTicker: boolean;
      logoUrl?: string | null;
      backgroundUrl?: string | null;
    };
  };
  export let active = false;
</script>

<button
  type="button"
  aria-pressed={active}
  class={`store-card ${active ? 'active' : ''}`}
>
  <div class="store-card__header" style={`background:${store.theme.primaryColor}`}>
    {#if store.theme.logoUrl}
      <img src={store.theme.logoUrl} alt={`${store.name} logo`} />
    {:else}
      <span>{store.name.substring(0, 2).toUpperCase()}</span>
    {/if}
  </div>
  <div class="store-card__body">
    <p class="store-card__slug">{store.slug}</p>
    <h3>{store.name}</h3>
    {#if store.theme.headline}
      <p class="store-card__headline">{store.theme.headline}</p>
    {/if}
    <p class="store-card__url">{store.publicUrl}</p>
    <div class="store-card__presets">
      {#each store.theme.presetAmounts.slice(0, 3) as amount}
        <span>{amount} XMR</span>
      {/each}
    </div>
  </div>
</button>

<style>
  .store-card {
    border: 2px solid transparent;
    border-radius: 18px;
    padding: 0;
    display: flex;
    flex-direction: column;
    background: #fff;
    box-shadow: 0 10px 35px rgba(10, 23, 49, 0.08);
    transition: all 0.2s ease;
    text-align: left;
  }

  .store-card:hover {
    transform: translateY(-2px);
    box-shadow: 0 15px 40px rgba(10, 23, 49, 0.12);
  }

  .store-card.active {
    border-color: #ff7a18;
  }

  .store-card__header {
    min-height: 120px;
    border-radius: 16px 16px 0 0;
    display: flex;
    align-items: center;
    justify-content: center;
    color: #fff;
    font-weight: 600;
    font-size: 1.5rem;
  }

  .store-card__header img {
    width: 72px;
    height: 72px;
    border-radius: 50%;
    object-fit: cover;
    border: 3px solid rgba(255, 255, 255, 0.6);
  }

  .store-card__body {
    padding: 1.25rem 1.5rem 1.5rem;
    display: flex;
    flex-direction: column;
    gap: 0.35rem;
  }

  .store-card__slug {
    font-size: 0.85rem;
    text-transform: uppercase;
    letter-spacing: 0.16em;
    color: rgba(9, 24, 41, 0.5);
    margin: 0;
  }

  .store-card__headline {
    color: rgba(9, 24, 41, 0.7);
    margin: 0;
  }

  .store-card__url {
    color: #3e4a59;
    font-size: 0.9rem;
    margin: 0.35rem 0;
    word-break: break-all;
  }

  .store-card__presets {
    display: flex;
    flex-wrap: wrap;
    gap: 0.4rem;
  }

  .store-card__presets span {
    background: rgba(6, 18, 33, 0.06);
    border-radius: 999px;
    padding: 0.15rem 0.65rem;
    font-size: 0.8rem;
  }
</style>
