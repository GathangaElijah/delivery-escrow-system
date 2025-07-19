# Ethereum to Polkadot Escrow Bridge

## Overview

This feature enables users to deposit funds from Ethereum-compatible wallets (e.g., MetaMask) into a Polkadot-based escrow system implemented using ink!. The Ethereum deposit acts as the funding mechanism, and once the business contract is fulfilled on the Polkadot side, the funds are settled accordingly.

This document describes the architecture and implementation plan for building this feature.

## Objective

Build a cross-chain deposit bridge that:

- Accepts deposits from Ethereum wallets (MetaMask).
- Detects and maps deposits to an internal escrow ID.
- Bridges the deposit event into the Polkadot-based escrow system.
- Interacts with an existing ink! contract for escrow lifecycle management.

## User Flow

1. User connects MetaMask on the frontend.
2. User deposits ETH or ERC20 tokens for a specific escrow transaction.
3. Ethereum transaction is confirmed.
4. Backend (written in Go) detects and verifies the deposit.
5. Backend interacts with the ink! smart contract on Polkadot to register the deposit and lock value.
6. Funds are held until escrow terms are fulfilled.

## Architecture

| Layer                  | Technology                                    | Role                                                                                               |
|------------------------|-----------------------------------------------|----------------------------------------------------------------------------------------------------|
| Frontend               | JavaScript + Web3.js or Ethers.js             | Wallet connection and deposit interface                                                            |
| Ethereum Deposit Layer | Solidity Smart Contract or Go wallet listener | Receives and logs deposit transactions; either through event emissions or direct wallet monitoring |
| Backend                | Go (with `go-ethereum`, JSON-RPC)             | Handles deposit logic, maps users to escrows, interacts with Polkadot                              |
| Polkadot Layer         | ink! smart contract                           | Escrow logic, settlement management                                                                |

> **Note:** Use a Solidity smart contract for explicit on-chain logging and traceability; use a Go wallet listener for reduced deployment complexity and fees. Choose based on your scalability and trust requirements.

## Smart Contract Option (ETH/Token Deposit)

```solidity
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract EscrowBridge {
    event Deposited(address indexed from, uint amount, string escrowId);

    function deposit(string memory escrowId) public payable {
        require(msg.value > 0, "No ETH sent");
        emit Deposited(msg.sender, msg.value, escrowId);
    }
}
