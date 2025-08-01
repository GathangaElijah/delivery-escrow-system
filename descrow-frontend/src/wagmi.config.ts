import { http, createConfig } from 'wagmi'
import { mainnet, sepolia } from 'wagmi/chains'
import { walletConnect, injected } from 'wagmi/connectors'

import { projectId } from './constants'


export const wagmiConfig = createConfig({
  chains: [mainnet, sepolia],
  connectors: [
    injected(), // e.g., MetaMask
    walletConnect({
      projectId: projectId, // from WalletConnect Cloud
    }),
  ],
  transports: {
    [mainnet.id]: http(),
    [sepolia.id]: http(),
  },
  ssr: false,
})
