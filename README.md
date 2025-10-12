# Project Name

# Delivery Escrow System (DES)
![Supply Chain Diagram](/DES/DesSupplychain.jpg)

A blockchain-based escrow system for secure e-commerce,freelancing and other service providers transactions, and prevent clients from fraud and scams.

![License](https://img.shields.io/github/license/yourusername/yourproject)
![Build](https://img.shields.io/github/actions/workflow/status/yourusername/yourproject/build.yml)
![Version](https://img.shields.io/badge/version-1.0-blue)

---

## 📌 Table of Contents

- [Overview](#-overview)
- [Problem Statement](#-problem-statement)
- [Solution](#-solution)
- [Tech Stack](#-tech-stack)
- [Getting Started](#-getting-started)
- [Usage](#-usage)
- [API Reference](#-api-reference) _(optional)_
- [Contributing](#-contributing)
- [License](#-license)
- [Author](#-author)

---

## 📖 Overview

The **Descrow Payment Protocol** is a blockchain-based escrow mechanism that ensures secure and transparent online transactions between buyers and sellers. It eliminates the risk of scams by only releasing funds when the buyer confirms successful delivery. The protocol is designed to provide decentralized trust, removing the need for intermediaries in peer-to-peer commerce.

---

## Problem Statement
Online transactions often suffer from **scams and fraud**, particularly through **seller impersonation traps**—where unsuspecting buyers send payments to fake sellers who disappear after receiving the funds. Traditional online payment systems lack a verifiable way to ensure that products are delivered before funds are released, leading to widespread mistrust in digital marketplaces.

---

## Solution
Descrow introduces a **smart contract-based payment protocol** that safely holds the buyer’s payment in escrow until the product’s delivery is confirmed. Both buyer and seller are onboarded and verified on-chain. Once the buyer confirms receipt of goods (and product authenticity is verified), the funds are automatically released to the seller.

---
## ⚙️ Tech Stack

| Layer          | Technology                |
| -------------- | ------------------------- |
| Smart Contract | Rust, ink!                |
| Backend        | Rust          
| Frontend       | React, Typescript
| Blockchain     | Polkadot (Substrate)      |
| Other Tools    | Docker, Git, curl         |

---

## 🛠️ Installation

## Project Structure

## Workflow
1. Buyer connects their crypto wallet to Descrow.
2. Buyer creates a purchase order and stakes payment into the escrow.
3. Seller accepts the order and dispatches the product.
4. Product tracking and seal verification (if applicable) are logged to ensure authenticity.
5. Buyer confirms receipt → funds are released to the seller.
6. If there’s a complaint → the transaction enters conflict resolution.

## Conflict Resolution
If the buyer reports a complaint (e.g., undelivered or tampered goods), the system triggers a **dispute resolution process**:
- The funds remain locked in escrow.
- Both parties submit evidence (e.g., delivery logs, product hashes, or proof of shipment).
- A **trusted arbitrator or decentralized resolution protocol** reviews the claim.
- Based on the verdict, funds are either refunded to the buyer or released to the seller.


## Security Measures
- **Immutable Escrow Contracts**: Funds are handled by smart contracts only.
- **Double Verification**: Delivery and authenticity checks before releasing funds.
- **Tamper Detection**: Optional QR seal verification for physical goods.
- **No Central Custodian**: All logic is handled on-chain.

## Future Enhancements
- **Cross-Network Wallet Support** (MetaMask, Lisk, Stellar, etc.)
- **Automated Dispute Arbitration** using decentralized governance.
- **AI-Powered Fraud Detection** to prevent impersonation or fake listings.

---

## License

This project is licensed under the GPL License - see the [LICENSE](./LICENSE) file for details.

## Contact

For questions or support, please contact [Elijah Gathanga](elyg3672@gmail.com). <br>
© 2025 Descrow Labs. All rights reserved.


