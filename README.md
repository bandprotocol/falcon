# Falcon

<div align="center">

![logo](logo.svg)

</div>

**Disclaimer:** This project is still in its early stages of development and is considered a prototype. Please refrain from using it in production.

`Falcon` is a CLI program designed for smart contract developers to relay data from the [Tunnel](https://github.com/bandprotocol/chain/tree/master/x/tunnel) on BandChain to their target blockchain seamlessly.

# Table Of Contents
- [Components](#components)
- [Flow](#flow)
- [Getting Started](#getting-started)

---
## Components 
### 1. BandChain Client
- Queries tunnel information.
- Queries packet data and EVM signatures for verification in the TSS verifier.

### 2. Chain Provider
Intermediary between the relayer application and the target blockchain. It encapsulates all logic required to interact with a specific type of blockchain, such as Ethereum Virtual Machine (EVM)-based chains, and abstracts the details of blockchain operations like transaction creation, execution, querying, and validation.
(At this moment, only EVM chains are supported.)
- Keys Manipulation
  - Manages key operations such as adding, deleting, listing, and securely storing them.
- Transaction Creation 
  - Constructs calldata for packet relaying.
  - Dynamically calculates gas fees for both EIP-1559 and Legacy transactions.
- Transaction Execution
  - Signs transactions using keys from the senders pool.
  - Broadcasts transactions to the blockchain. 
  - Verifies transaction success with a receipt check.

### 3. Scheduler
- Executes tasks for each `TunnelRelayer` periodically.
- Penalizes failing tasks with an exponential backoff.
- Synchronizes with BandChain periodically to ensure that Falcon keeps all tunnels aligned and up-to-date with BandChain.

### 4. Tunnel Relayer
- Fetches tunnel information and packet by `BandChain Client` and `Chain Provider`
- Handles the packet relaying process
  - Validates tunnel state (e.g., active status, latest sequence).
  - Fetches unrelayed packets from BandChain.
  - Submit packet to `Chain Provider` to continue on transaction process.


---
## Flow
### 1. Initialization
- Initializes Tunnel Relayers
  - Fetches tunnel information from BandChain and target chains.
  - Validates the provider and loads keys into the sender channel.
- Validates Passphrase 
  - Checks if the user-provided passphrase matches the stored hash to ensure security. 
- Starts Scheduler
  - Scheduler manages periodic relaying tasks for all initialized `TunnelRelayer` instances.

### 2. Starting Scheduler
The Scheduler starts execution ticker and tunnel synchronization ticker based on the given configuration, and handles penalized tasks resulting from consecutive relay packet failures.
- Periodic Execution
  - Asynchronously triggers each `TunnelRelayer` instances to check [(3.)](#3-checking-packets) and relay [(4.)](#4-relaying-packets) the packet and confirm transaction [(5.)](#5-confirming-transaction).
- Tunnels synchronization
  - Updates tunnel information from BandChain to ensure that Falcon keeps all tunnels aligned and up-to-date with BandChain.
- Handling penalty task
  - Retries failed tasks after waiting for a penalty duration, which is calculated using exponential backoff.

### 3. Checking Packets
- Checks Tunnel State
  - Queries tunnel information from BandChain and the target chain.
  - Validates if the target chain's tunnel is active and ensures sequence consistency.
- Determine Next Packet
  - Compares the LatestSequence on the BandChain and the target chain.
  - Identifies packets that have not yet been relayed.

### 4. Relaying Packets 
- Validate Connection 
  - Ensures the chain client is connected.
- Retry Logic 
  - Attempts to relay the packet up to `max_retry` times.
    If all attempts fail, logs the error and move this task to penalty task.
- #### 4.1 Create Transaction
  - Fetch the sender's nonce to ensure transaction order.
  - Calculate Gas Limit by using a pre-configured limit or dynamically estimate it based on transaction parameters.
  - Set Fees by GasPrice for legacy transactions or BaseFee + PriorityFee for EIP-1559 transactions.
- #### 4.2 Sign Transaction
  - Sign the transactionby the sender's private key to sign the transaction with the chosen signer 
    - EIP155Signer for legacy transactions.
    - LondonSigner for EIP-1559 transactions.
- #### 4.3 Broadcast Transaction 
  - Send the transaction to the target chain.
### 5. Confirming Transaction
- Fetch the transaction receipt and check its status and block confirmations to ensure success.
- Continuously check the transaction status within a specified duration. If the transaction remains unconfirmed for too long, mark it as unmined and handle it accordingly.
- If the transaction fails, adjust the gas by using `GasMultiplier` and retry to relay packet again

---
## Getting Started
**Note:** This guide assumes that the program runs on a Linux environment.
### 1. Node Installation

#### 1.1 Node Configuration
- make, gcc, g++ (can be obtained from the build-essential package on linux)
- wget, curl for downloading files

``` shell
# install required tools
sudo apt-get update && \
sudo apt-get upgrade -y && \
sudo apt-get install -y build-essential curl wget jq
```
- Install Go 1.22.3
```shell
# Install Go 1.22.3
wget https://go.dev/dl/go1.22.3.linux-amd64.tar.gz
tar xf go1.22.3.linux-amd64.tar.gz
sudo mv go /usr/local/go

# Set Go path to $PATH variable
echo "export PATH=$PATH:/usr/local/go/bin:~/go/bin" >> $HOME/.profile
source ~/.profile
```
Go binary should be at /usr/local/go/bin and any executable compiled by go install command should be at ~/go/bin

#### 1.2: Clone & Install Falcon
``` shell
cd ~
# Clone Falcon
git clone https://github.com/bandprotocol/falcon
cd falcon

# Install binaries to $GOPATH/bin
make install
```

### 2. Set passphrase 
In Falcon, the passphrase serves as an **encryption key** for securing sensitive data.
- The passphrase can be set in one of the following ways:
  - As an environment variable via a `.env` file.
    ``` shell
    PASSPHRASE=secret
    ```
  - Passed inline with commands that require it.
    ``` shell
    PASSPHRASE=secret falcon ...
    ```
- Some commands require the passphrase to:
  - Validate its correctness when performing operations like managing keys.
  - Initialize and encrypt configuration files for the program.
- If no passphrase is provided, it defaults to an empty string `""`.

### 3. Initialize the configuration directory/file
``` shell
falcon config init
```
- The passphrase will be init after this command, it will become your program's permanent passphrase and cannot be changed. Please check that you set the correct passphrase during this step.

- Default config directory: ~/.falcon/config

- By default, config will be initialized in format like this 
```toml
[global]
log_level = 'info'
checking_packet_interval = 60000000000
max_checking_packet_penalty_duration = 3600000000000
penalty_exponential_factor = 1.0

[bandchain]
rpc_endpoints = ['http://localhost:26657']
timeout = 3000000000

[target_chains]

```


To customize the config for relaying, you can use custom config file and use the `--file` flag when initializing the configuration.
```
falcon config init --file custom_config.toml
```

### 4. Configure target chains you want to relay
You need to create a chain configuration file to add it to the configuration. Currently, only EVM chains are supported.
<br/> Example:
``` toml
endpoints = ['http://localhost:8545']
chain_type = 'evm'
max_retry = 3
query_timeout = 3000000000
chain_id = 31337
tunnel_router_address = '0xDc64a140Aa3E981100a9becA4E685f962f0cF6C9'
private_key = ''
block_confirmation = 5
waiting_tx_duration = 3000000000
checking_tx_interval = 1000000000
gas_type = 'eip1559'
gas_multiplier = 1.1
execute_timeout = 3000000000
liveliness_checking_interval = 900000000000
```

The supported `gas_type` values are `legacy` and `eip1559`. Each type requires specific configuration fields.
- `legacy`
  - `max_gas_price` defines the maximum gas price.
  - If `max_gas_price` is not specified, it will be retrieved from the tunnel router.
- `eip1559`
  - `max_base_fee` defines the maximum base fee.
  - `max_priority_fee` defines the maximum priority fee.
   - If `max_priority_fee` is not defined, it will also be retrieved from the tunnel router

``` shell
falcon chains add testnet chain_config.toml
```

### 5. Check target chain's activeness
To relay packets to the target chain, you need to ensure that the tunnel on the target chain is active. This can be checked using
``` shell
falcon query tunnel <TUNNEL_ID>
```


### 6. Import OR create new keys to use when signing and relaying transactions.
>Please ensure that you are using the correct passphrase that was set during initialization for the `add`, `delete`, and `export` commands.
</br>
>
If you need to generate a new private key you can use the add subcommand.
``` shell
falcon keys add testnet testkey
```

There are 3 options for user to add key 
``` shell 
Choose how to add a key
> Private key (provide an existing private key)
  Mnemonic (recover from an existing mnemonic phrase)
  Generate new address (no private key or mnemonic needed)
```

If you already have a private key and want to retrive key from it, you can choose `Private key` option. 
``` shell
Enter your private key
>
```

If you already have a mnemonic and want to retrive key from it, you can choose `Mnemonic` option. 
``` shell
Enter your mnemonic
>
```

### 7. Check that the keys for the configured chains are funded

You can query the balance of each configured key by running:
``` shell
falcon q balance testkey
```
### 8. Start to relay packet
Starts all tunnels that `falcon query tunnels` can query
``` shell
falcon start
```
> NOTE: you can choose which tunnels do you want to relay
``` shell
falcon start 1 2 3
```