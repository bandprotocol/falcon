# Falcon

<div align="center">

![logo](logo.svg)

</div>

**Disclaimer:** This project is still in its early stages of development and is considered a prototype. Please refrain from using it in production.

`Falcon` is a CLI program designed for smart contract developers from any blockchains to easily request and retrieve data from BandChain.

---
## Components 
### 1. BandChain Client
- Queries tunnel information.
- Retrieves packet data and EVM signatures for verification in the TSS verifier.

### 2. Chain Provider
intermediary between the relayer application and the target blockchain. It encapsulates all logic required to interact with a specific type of blockchain, such as Ethereum Virtual Machine (EVM)-based chains, and abstracts the details of blockchain operations like transaction creation, execution, querying, and validation.
(At this moment, only EVM chains are supported.)
- Keys Management
  - sads
- Transaction Creation 
  - Constructs calldata for packet relaying using ABI methods.
  - Dynamically calculates gas fees for both EIP-1559 and Legacy transactions.
- Transaction Execution
  - Signs transactions using keys from the senders pool.
  - Broadcasts transactions to the blockchain. 
  - Verifies transaction success with a receipt check.

- Balance Querying 
  - Fetches account balances for monitoring key usage.

### 3. Scheduler
- Executes tasks for each `TunnelRelayer` periodically.
- Exponential backoff for error scenarios.
- Synchornize latest tunnel from the Bandchain periodically.

### 4. Tunnel Relayer
- Prerequisite handles the packet relaying process
  - Validates tunnel state (e.g., active status, latest sequence).
  - Fetches unrelayed packets from BandChain.
  - Submit packet to `Chain Provider` to continue on transaction process.


---
## Flow
### 1. Initialization
- Initialize Tunnel Relayers
   - Retrieves tunnel information from BandChain and target contract
   - Validates the provider and loads keys into the FreeSenders channel.
- Validate Passphrase 
  - Checks if the user-provided passphrase matches the stored hash to ensure security. 
- Start Scheduler
  - Scheduler manages periodic relaying tasks for all initialized `TunnelRelayer` instances.

### 2. Starting Scheduler
The Scheduler starts execution time and tunnel synchronization tickers based on the given configuration and handles penalized tasks resulting from consecutive relay packet failures.
- Periodic Execution
  - asynchronously trigger each `TunnelRelayer` instances to check (3.) and relay (4.) the packet
- Tunnels synchronization
  - Updates tunnel information from BandChain to ensure relayers operate on the latest tunnel metadata.
- Handling penalty task
  - Retries failed tasks after waiting for a penalty duration, which is calculated using exponential backoff.

### 3. Checking Packets
- Check Tunnel State
  - Queries tunnel information from BandChain and the target chain.
  - Validates if the target chain's tunnel is active and ensures sequence consistency.
- Determine Next Packet
  - Compares the LatestSequence on the BandChain and the target chain.
  - Identifies packets that have not yet been relayed.

### 4. Relaying Packets 
- Validate Connection 
  - Ensures the chain client is connected.
- Retry Logic 
  - Attempts to relay the packet up to `max_retry times.
    If all attempts fail, logs the error and move this task to penalty task.
#### 4.1 Create Transaction
#### 4.2 Sign Transaction
#### 4.3







