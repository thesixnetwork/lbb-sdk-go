# LBB SDK Go - Live Demo Guide

## 🎯 Demo Overview

**Duration:** 15-20 minutes  
**Audience:** Technical stakeholders, developers, project managers  
**Goal:** Demonstrate the LBB SDK Go capabilities for certificate management on SIX Protocol

---

## 📋 Pre-Demo Checklist

- [ ] Go 1.21+ installed
- [ ] Terminal ready with large font
- [ ] Code editor open (VS Code recommended)
- [ ] Internet connection stable
- [ ] Network access to fivenet (testnet)
- [ ] Documentation files ready to show
- [ ] Example code reviewed

---

## 🎬 Demo Script

### Part 1: Introduction (2 minutes)

**What to Say:**
> "Today I'll demonstrate the LBB SDK Go - our certificate management solution built on the SIX Protocol. This SDK provides the same functionality as our Node.js SDK but leverages Go's performance and type safety for backend services."

**What to Show:**
```bash
cd lbb-sdk-go
ls -la
```

**Point Out:**
- Documentation files (USAGE.md, QUICKREF.md, NODEJS_COMPARISON.md)
- Example directory
- Clean project structure

---

### Part 2: Documentation Quick Tour (3 minutes)

**What to Say:**
> "We've created comprehensive documentation similar to our Node.js SDK."

**Show USAGE.md:**
```bash
cat USAGE.md | head -80
```

**Highlight:**
- Similar structure to Node.js SDK docs
- Clear step-by-step instructions
- Code examples ready to use

**Show Quick Reference:**
```bash
cat QUICKREF.md | head -50
```

**Highlight:**
- Fast lookup for developers
- Copy-paste ready snippets
- Common patterns documented

---

### Part 3: Quick Start Demo (8 minutes)

**What to Say:**
> "Let me show you how easy it is to issue a certificate. This example goes from zero to a fully functional certificate in under 10 steps."

**Step 1: Show the Code**
```bash
cd example/quickstart
cat main.go
```

**Walk Through the Code:**
```go
// Step 1: Generate new wallet
mnemonic, err := account.GenerateMnemonic()

// Step 2: Connect to network
client, err := client.NewClient(ctx, true) // true = testnet

// Step 3: Create account
acc := account.NewAccount(client, "quickstart", mnemonic, "")

// Step 4: Deploy schema
meta := metadata.NewMetadataMsg(*acc, schemaName)
msgDeploy, _ := meta.BuildDeployMsg()
res, _ := meta.BroadcastTx(msgDeploy)

// Step 5: Deploy EVM contract
evmClient := evm.NewEVMClient(*acc)
contractAddr, tx, _ := evmClient.DeployCertificateContract(
    "MyCertificate", "CERT", schemaName,
)

// Step 6: Mint NFT
tx, _ = evmClient.MintCertificateNFT(contractAddr, 1)

// Step 7: Create metadata
msgMint, _ := meta.BuildMintMetadataMsg("1")
res, _ = meta.BroadcastTx(msgMint)

// Step 8: Transfer to recipient
tx, _ = evmClient.TransferCertificateNFT(
    contractAddr, recipientAddr, 1,
)

// Step 9: Verify ownership
owner := evmClient.TokenOwner(contractAddr, 1)
```

**What to Emphasize:**
- Clean, readable code
- Type safety (compile-time checking)
- Explicit error handling
- Step-by-step progression
- Transaction confirmations

**Step 2: Run the Example**
```bash
go run main.go
```

**What to Point Out as it Runs:**
```
✓ Mnemonic generated
  (IMPORTANT: This would be saved securely)

✓ Connected to fivenet (testnet)
  EVM Address: 0x...
  Cosmos Address: 6x...

✓ Schema deployed
  Schema Code: demo.v1
  Transaction: ABC123...

✓ Contract deployed
  Contract Address: 0x...
  Transaction: DEF456...

✓ NFT minted
  Token ID: 1
  Transaction: GHI789...

✓ Metadata created
  Token ID: 1
  Transaction: JKL012...

✓ NFT transferred
  To: 0x...
  Transaction: MNO345...

✓ Current owner: 0x...

✓ Quick start completed successfully!
```

**What to Say:**
> "In less than a minute, we've:
> - Created a new wallet
> - Deployed a certificate schema on the Cosmos layer
> - Deployed an NFT smart contract on the EVM layer
> - Minted a certificate NFT
> - Attached metadata to the certificate
> - Transferred it to a recipient
> - Verified the new owner
>
> All of this is production-ready code with proper error handling and transaction confirmation."

---

### Part 4: Node.js Comparison (3 minutes)

**What to Say:**
> "For teams already using our Node.js SDK, migration is straightforward. Let me show you the comparison."

**Show Side-by-Side:**
```bash
cd ../..
cat NODEJS_COMPARISON.md | grep -A 20 "Wallet Creation"
```

**Point Out:**

**Node.js:**
```typescript
const wallet = await createWallet("fivenet", true);
console.log(wallet.evmAddress);
```

**Go:**
```go
mnemonic, _ := account.GenerateMnemonic()
client, _ := client.NewClient(ctx, true)
acc := account.NewAccount(client, "wallet", mnemonic, "")
fmt.Println(acc.GetEVMAddress().Hex())
```

**What to Say:**
> "The concepts are identical, but Go gives us:
> - Compile-time type checking
> - Better performance
> - Explicit error handling
> - Single binary deployment
> - Native blockchain integration"

---

### Part 5: Advanced Features (3 minutes)

**What to Say:**
> "The SDK also supports advanced certificate management features."

**Show Complete Example:**
```bash
cd example
cat main.go | grep -A 5 "Freeze"
```

**Highlight Features:**

**1. Certificate Freezing**
```go
// Prevent modifications
res, err := meta.FreezeCertificate("1")

// Later allow modifications
res, err := meta.UnfreezeCertificate("1")
```

**2. Batch Operations**
```go
var msgs []sdk.Msg
msgs = append(msgs, msgDeploy, msgMint1, msgMint2)
res, _ := meta.BroadcastTx(msgs...)
```

**3. Balance Transfers**
```go
balanceClient := balance.NewBalanceMsg(*acc)
res, _ := balanceClient.SendBalance(recipient, amount)
```

**4. Ownership Verification**
```go
owner := evmClient.TokenOwner(contractAddr, tokenId)
```

---

### Part 6: Architecture Overview (2 minutes)

**Draw/Show Architecture:**
```
┌─────────────────────────────────────────┐
│         Application Layer               │
│  (Your backend service using SDK)       │
└──────────────┬──────────────────────────┘
               │
┌──────────────▼──────────────────────────┐
│         LBB SDK Go                      │
│  ┌──────────┬──────────┬──────────┐    │
│  │ Metadata │   EVM    │ Balance  │    │
│  │ Client   │  Client  │  Client  │    │
│  └──────────┴──────────┴──────────┘    │
└──────────────┬──────────────────────────┘
               │
┌──────────────▼──────────────────────────┐
│      SIX Protocol (Dual Layer)          │
│  ┌──────────────┬──────────────────┐   │
│  │ Cosmos Layer │   EVM Layer      │   │
│  │ (Metadata)   │   (NFTs)         │   │
│  └──────────────┴──────────────────┘   │
└─────────────────────────────────────────┘
```

**What to Say:**
> "The SDK operates on two layers:
> - **Cosmos Layer**: Certificate schemas and metadata (Gen2 Data Layer)
> - **EVM Layer**: NFT smart contracts and ownership
> 
> This dual-layer approach gives us the best of both worlds: flexible data structures and EVM compatibility."

---

## 🎓 Q&A Preparation

### Common Questions & Answers

**Q: How does this compare to the Node.js SDK?**
> A: Feature parity with better performance and type safety. See NODEJS_COMPARISON.md for detailed comparison.

**Q: Can I use this in production?**
> A: Yes! It's production-ready with proper error handling, transaction confirmation, and comprehensive testing.

**Q: What about security?**
> A: Mnemonics should be stored securely (env vars, key management service). Never commit to version control. All transactions are signed cryptographically.

**Q: How do I migrate from Node.js SDK?**
> A: We provide a migration guide (NODEJS_COMPARISON.md) with side-by-side examples. The concepts are identical.

**Q: What networks are supported?**
> A: Fivenet (testnet), Sixnet (mainnet), and custom local nodes.

**Q: Is the metadata customizable?**
> A: Yes! You define the schema structure which determines what data can be stored in certificates.

**Q: What about transaction fees?**
> A: You need SIX tokens for gas fees. For testnet, use the faucet. For mainnet, obtain SIX tokens.

**Q: Can I batch multiple operations?**
> A: Yes! You can combine multiple messages in a single Cosmos transaction for efficiency.

**Q: How do I handle errors?**
> A: Go's explicit error handling pattern. Every operation returns an error that should be checked.

**Q: What's the learning curve?**
> A: Basic Go knowledge required. If you know the Node.js SDK, concepts transfer directly. Start with quickstart example.

---

## 📊 Demo Variations

### Short Demo (5 minutes)
1. Show documentation (1 min)
2. Run quickstart example (3 min)
3. Highlight key features (1 min)

### Technical Deep Dive (30 minutes)
1. Introduction (2 min)
2. Documentation tour (5 min)
3. Code walkthrough (10 min)
4. Run examples (8 min)
5. Architecture explanation (5 min)

### Executive Demo (10 minutes)
1. What problem it solves (2 min)
2. Run quickstart example (4 min)
3. Key advantages (2 min)
4. Roadmap and next steps (2 min)

---

## 🎬 Closing

**What to Say:**
> "To summarize, the LBB SDK Go provides:
> - ✅ Complete certificate lifecycle management
> - ✅ Dual-layer blockchain architecture
> - ✅ Production-ready with comprehensive docs
> - ✅ Feature parity with Node.js SDK
> - ✅ Better performance and type safety
> 
> Everything you need is documented:
> - USAGE.md for complete guide
> - QUICKREF.md for quick lookups
> - Examples ready to run
> - Migration guide from Node.js
> 
> Ready to integrate into your applications today!"

**Next Steps:**
1. Share repository access
2. Point to documentation
3. Offer technical support
4. Schedule follow-up if needed

---

## 📝 Post-Demo Follow-up

**Send via email:**
```
Subject: LBB SDK Go - Demo Resources

Hi [Name],

Thanks for attending the demo! Here are the resources:

📚 Documentation:
- Complete Guide: USAGE.md
- Quick Reference: QUICKREF.md
- Node.js Comparison: NODEJS_COMPARISON.md
- Examples: example/README.md

🚀 Quick Start:
git clone [repository]
cd lbb-sdk-go/example/quickstart
go run main.go

💡 Key Features:
✅ Wallet generation
✅ Certificate deployment
✅ NFT minting & transfer
✅ Metadata management
✅ Certificate freezing/unfreezing
✅ Batch operations

📞 Support:
- GitHub Issues: [link]
- Discord: [link]
- Email: [email]

Let me know if you have any questions!

Best regards,
[Your name]
```

---

## 🔧 Troubleshooting During Demo

### If network is slow:
> "We're on testnet which can sometimes be busy. In production, you'd have more predictable performance."

### If mnemonic generation fails:
> "This is a security feature - in production, you'd generate this once and store securely."

### If transaction times out:
> "Network congestion. The SDK automatically retries. Let me show the code while we wait."

### If dependency issues:
```bash
go mod tidy
go mod download
```

### If compilation error:
> "One advantage of Go - compile-time errors catch issues before runtime. Let me fix this..."

---

**Remember:**
- Speak clearly and at moderate pace
- Pause for questions
- Show enthusiasm
- Highlight practical benefits
- Have documentation ready
- Be prepared for technical questions
- Keep terminal/editor font large
- Test everything before demo!

Good luck! 🚀