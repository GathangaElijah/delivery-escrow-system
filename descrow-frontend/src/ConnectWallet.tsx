'use client'

import { useState } from 'react'
import { useConnect, useDisconnect, useAccount } from 'wagmi'
import { injected } from 'wagmi/connectors'

export default function WalletConnect() {
  const { connect, connectors, isPending } = useConnect()
  const { disconnect } = useDisconnect()
  const { address, isConnected } = useAccount()
  const [connecting, setConnecting] = useState(false)

  const handleConnect = async (connectorId: string) => {
    const connector = connectors.find((c) => c.id === connectorId)
    if (!connector) return
    setConnecting(true)
    try {
      await connect({ connector })
    } finally {
      setConnecting(false)
    }
  }

  return (
    <div>
      {isConnected ? (
        <div>
          <p>Connected: {address}</p>
          <button onClick={() => disconnect()}>Disconnect</button>
        </div>
      ) : (
        <div>
          {connectors.map((connector) => (
            <button
              key={connector.id}
              onClick={() => handleConnect(connector.id)}
              disabled={isPending || connecting}
            >
              {`Connect with ${connector.name}`}
            </button>
          ))}
        </div>
      )}
    </div>
  )
}
