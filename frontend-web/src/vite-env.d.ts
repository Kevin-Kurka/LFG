/// <reference types="vite/client" />

interface ImportMetaEnv {
  readonly VITE_API_BASE_URL: string;
  readonly VITE_WALLET_SERVICE_URL: string;
  readonly VITE_MARKET_SERVICE_URL: string;
  readonly VITE_ORDER_SERVICE_URL: string;
  readonly VITE_CREDIT_SERVICE_URL: string;
  readonly VITE_SPORTSBOOK_SERVICE_URL: string;
  readonly VITE_WS_URL: string;
}

interface ImportMeta {
  readonly env: ImportMetaEnv;
}
