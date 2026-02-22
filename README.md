
# MEV Arbitrage Scanner (Go)

A minimal, event-driven arbitrage scanner built in Go for Uniswap V2-style AMMs.

This project monitors two liquidity pools in real-time, maintains in-memory reserve state, and simulates deterministic arbitrage opportunities between them.

It is designed for educational and research purposes to understand MEV mechanics, AMM math, and on-chain event-driven systems.

---

## 📌 What This Project Does

* Connects to an Ethereum-compatible chain via WebSocket RPC
* Subscribes to `Sync` events from two Uniswap V2 pools
* Maintains live reserve state in memory
* Simulates swap paths using exact constant product AMM math
* Detects theoretical arbitrage opportunities
* Runs entirely event-driven (no polling loop)

---

## 🧠 Core Concept

Uniswap V2 pools follow the invariant:

```
x * y = k
```

Swap output formula (0.3% fee):

```
amountOut =
(amountIn * 997 * reserveOut) /
(reserveIn * 1000 + amountIn * 997)
```

If two pools quote different effective prices for the same pair, a round-trip swap may yield a profit.

This scanner detects:

```
Token0 -> Token1 (Pool A)
Token1 -> Token0 (Pool B)
```

If:

```
outB > amountIn
```

an arbitrage opportunity is reported.

---

## 🏗 Architecture

```
WebSocket RPC
      ↓
Sync Event Listener (Pool A, Pool B)
      ↓
In-memory Reserve State
      ↓
Arbitrage Simulation Engine
      ↓
Opportunity Logger
```

Key design decisions:

* Exact integer math using `big.Int`
* ABI-based decoding
* Non-blocking event trigger
* No floating point arithmetic
* No polling loops

---

## 📁 Project Structure

```
mev-scanner/
    cmd/main.go
    internal/
        config/
        pool/
        math/
        strategy/
```

---

## 🚀 How To Run

From project root:

```
go mod tidy
go run ./cmd
```

Make sure to:

* Provide a valid WebSocket RPC URL
* Set two valid Uniswap V2 pair addresses
* Ensure both pools trade the same token pair

---

## 🎯 Learning Goals

This project helps understand:

* AMM pricing mechanics
* Event-driven blockchain systems
* Deterministic simulation
* Reserve-based arbitrage
* MEV fundamentals


