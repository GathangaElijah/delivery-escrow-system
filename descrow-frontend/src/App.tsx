import React from 'react'

import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { WagmiProvider } from 'wagmi'
import { wagmiConfig } from './wagmi.config'

import WalletConnect from './ConnectWallet'

import './App.css'

const queryClient = new QueryClient;

function App() {

  return (
    <WagmiProvider config={wagmiConfig}>
      <QueryClientProvider client={queryClient}>
        <div style={{ padding: '2rem' }}>
          <h1>My DApp</h1>
          <WalletConnect />
        </div>
      </QueryClientProvider>
    </WagmiProvider>
  )
}

export default App
