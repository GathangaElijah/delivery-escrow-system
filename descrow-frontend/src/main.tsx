import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import './index.css'
import App from './App.tsx'

import { WagmiProvider } from 'wagmi'
import { wagmiConfig } from './wagmi.config'

createRoot(document.getElementById('root')!).render(
  <StrictMode>
     <WagmiProvider config={wagmiConfig}>
      <App />
    </WagmiProvider>
  </StrictMode>,
)
