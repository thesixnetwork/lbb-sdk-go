# Welcome to LBB SDK Go! 🎉

## For Junior Developers

Welcome! This guide will help you get started with the LBB SDK Go for building certificate NFT applications. Don't worry if you're new to blockchain development—we'll guide you through everything step by step.

## 📖 What is This?

The LBB SDK (License-Based Blockchain SDK) helps you create and manage digital certificates as NFTs (Non-Fungible Tokens) on the blockchain. Think of it like creating digital diplomas, badges, or credentials that are:

- **Permanent**: Stored on blockchain forever
- **Verifiable**: Anyone can verify authenticity
- **Transferable**: Can be sent to others
- **Secure**: Protected by blockchain technology

## 🎯 Your Learning Path

We've created **9 separate example files** for you to learn each concept independently:

### Week 1: Basics (Files 1-3)
Start here to understand accounts and schemas.

1. **01_generate_wallet.go** - Learn how to create a wallet
2. **02_create_account.go** - Connect to the network
3. **03_deploy_schema.go** - Create a certificate template

### Week 2: Metadata & Contracts (Files 4-5)
Learn how to mint metadata and deploy smart contracts.

4. **04_mint_metadata.go** - Create certificate metadata
5. **05_deploy_contract.go** - Deploy your NFT contract

### Week 3: NFT Operations (Files 6-7)
Learn how to work with NFTs.

6. **06_mint_nft.go** - Create your first NFT
7. **07_transfer_nft.go** - Send NFT to someone

### Week 4: Advanced Features (Files 8-9)
Learn advanced features and queries.

8. **08_freeze_metadata.go** - Lock/unlock certificates
9. **09_query_nft.go** - Check NFT information

## 🚀 Getting Started (Day 1)

### Step 1: Setup Your Environment

```bash
# Navigate to the example directory
cd lbb-sdk-go/example

# Install dependencies
go mod download

# Verify Go is installed
go version  # Should show Go 1.21 or higher
```

### Step 2: Run Your First Example

```bash
# Generate a wallet
go run 01_generate_wallet.go
```

**What you'll see:**
- A 24-word mnemonic phrase
- Instructions to save it securely

**What you learned:**
- How wallet generation works
- Importance of mnemonic phrases

### Step 3: Understand the Code

Open `01_generate_wallet.go` in your editor and read through it. Notice:
- Clear comments explaining each step
- Simple, easy-to-follow structure
- Real-world examples

## 📚 How to Use These Examples

### Daily Practice Routine

**Day 1-3:** Run and read files 1-3
- Run each file
- Read the code carefully
- Take notes on what you don't understand

**Day 4-6:** Run and read files 4-5
- Compare with previous examples
- Notice patterns in the code

**Day 7-9:** Run and read files 6-7
- Practice modifying the code
- Try changing parameters

**Day 10-12:** Run and read files 8-9
- Experiment with different values
- Build confidence

### Learning Tips

1. **Read First**: Read the file before running it
2. **Run It**: Execute and see the output
3. **Modify**: Change small things and see what happens
4. **Document**: Write notes about what you learned
5. **Ask Questions**: Don't hesitate to ask for help

## 🎓 Key Concepts to Understand

### Blockchain Basics

**What is blockchain?**
- A database that everyone can see
- Records can't be changed once added
- Very secure and transparent

**What is a wallet?**
- Your identity on the blockchain
- Contains your private key (like a password)
- Has an address (like an email address)

**What is a transaction?**
- An action on the blockchain
- Costs a small fee (gas)
- Permanent once confirmed

### LBB SDK Concepts

**Schema** (Cosmos Layer)
- Template for your certificates
- Defines what data they contain
- Like a form template

**Metadata** (Cosmos Layer)
- Actual certificate data
- Instance of a schema
- Can be frozen/unfrozen

**NFT Contract** (EVM Layer)
- Smart contract that manages NFTs
- Like a digital certificate issuer
- Linked to a schema

**NFT Token** (EVM Layer)
- Individual certificate
- Has a unique ID
- Can be transferred

## 🛠️ Understanding the Two Layers

The LBB SDK uses TWO different blockchain layers:

### Cosmos Layer (Steps 3-4, 8)
**What it does:**
- Stores certificate templates (schemas)
- Stores certificate data (metadata)
- Manages data permissions

**When you use it:**
- Deploying schemas
- Minting metadata
- Freezing/unfreezing certificates

### EVM Layer (Steps 5-7, 9)
**What it does:**
- Manages NFT ownership
- Handles transfers
- Ethereum-compatible

**When you use it:**
- Deploying contracts
- Minting NFTs
- Transferring NFTs
- Querying ownership

**Important:** Both layers work together but are separate!

## 📝 Checklist for Each Example

Before running each example, check:

- [ ] I read the file header comments
- [ ] I understand what this example does
- [ ] I updated any required configuration
- [ ] I have the prerequisites ready
- [ ] I'm ready to take notes

After running each example:

- [ ] I saw the expected output
- [ ] I understand what happened
- [ ] I saved important values (addresses, hashes)
- [ ] I can explain it to someone else
- [ ] I'm ready for the next example

## 🐛 Common Mistakes (and How to Fix Them)

### Mistake 1: Forgot to Update Contract Address
**Error:** `"Contract address 0x0000..."`

**Fix:** Copy the contract address from step 5 output and paste it into the file

### Mistake 2: Using Wrong Schema Name
**Error:** `"Schema not found"`

**Fix:** Make sure the schema name matches exactly (case-sensitive)

### Mistake 3: Running Out of Order
**Error:** Various errors about missing dependencies

**Fix:** Follow the numbered order: 01 → 02 → 03 → ... → 09

### Mistake 4: No Testnet Tokens
**Error:** `"Insufficient funds"`

**Fix:** Use the provided test mnemonic (`account.TestMnemonic`)

## 🎯 Mini Projects to Try

After completing all examples, try these mini projects:

### Project 1: Create Your Own Certificate
- Use your own schema name
- Mint 5 different metadata instances
- Create NFTs for each

### Project 2: Certificate Transfer System
- Mint an NFT
- Transfer to another address
- Verify the transfer worked

### Project 3: Frozen Certificate
- Mint metadata
- Freeze it
- Try to unfreeze
- Document the process

## 📖 Additional Resources

### When You're Stuck

1. **Check TUTORIAL.md** - Detailed explanations
2. **Check QUICK_REFERENCE.md** - Code snippets
3. **Re-read the example comments** - Often has the answer
4. **Check error messages** - They tell you what's wrong

### Learning More

- **Go Programming**: [go.dev/tour](https://go.dev/tour)
- **Blockchain Basics**: Read TUTORIAL.md introduction
- **Cosmos SDK**: Understanding the Cosmos layer
- **Ethereum/EVM**: Understanding smart contracts

## 🤝 Getting Help

### Before Asking for Help

1. Read the error message carefully
2. Check the TUTORIAL.md troubleshooting section
3. Review the example comments
4. Try to Google the error

### How to Ask for Help

**Good Question:**
```
"I'm running 05_deploy_contract.go and getting error:
'schema not found'

I already ran 03_deploy_schema.go and got:
Schema Code: myorg.lbbv01
Transaction: ABC123...

My schema name in both files is 'myorg.lbbv01'

What am I missing?"
```

**Not Helpful:**
```
"It doesn't work, please help"
```

## 🎓 Success Criteria

You're ready to move on when you can:

- [ ] Run all 9 examples successfully
- [ ] Explain what each example does
- [ ] Modify example code confidently
- [ ] Deploy your own schema
- [ ] Create and transfer your own NFTs
- [ ] Debug common errors independently

## 🌟 Next Steps After Completion

Once you've mastered these examples:

1. **Build a Real Application**
   - Design your own certificate system
   - Implement it using the SDK
   - Test thoroughly on testnet

2. **Explore Advanced Features**
   - Custom metadata fields
   - Batch operations
   - Integration with frontend

3. **Deploy to Mainnet**
   - When you're confident
   - With real tokens
   - For production use

## 💪 Motivation

Remember:
- **Every expert was once a beginner**
- **Take your time** - Understanding is more important than speed
- **Mistakes are learning opportunities**
- **Ask questions** - There are no stupid questions
- **Practice daily** - Consistency is key

## 📅 Recommended 4-Week Schedule

### Week 1: Foundations
- **Monday**: Setup environment, run example 01
- **Tuesday**: Deep dive into example 01, run example 02
- **Wednesday**: Study example 02, run example 03
- **Thursday**: Understand example 03
- **Friday**: Review Week 1, take notes

### Week 2: Building Blocks
- **Monday**: Run and study example 04
- **Tuesday**: Deep dive example 04, run example 05
- **Wednesday**: Study example 05 thoroughly
- **Thursday**: Practice deploying schemas and contracts
- **Friday**: Review Week 2, mini project

### Week 3: NFT Operations
- **Monday**: Run and study example 06
- **Tuesday**: Deep dive example 06, run example 07
- **Wednesday**: Study example 07 thoroughly
- **Thursday**: Practice minting and transferring
- **Friday**: Review Week 3, create test certificates

### Week 4: Advanced & Review
- **Monday**: Run and study example 08
- **Tuesday**: Deep dive example 08, run example 09
- **Wednesday**: Study example 09, review all examples
- **Thursday**: Mini project: Build complete workflow
- **Friday**: Final review, document learnings

## 🎉 You're Ready!

Start with `01_generate_wallet.go` and work your way through. Take your time, read carefully, and don't hesitate to ask questions.

**Good luck on your blockchain development journey!** 🚀

---

**Questions?** Check TUTORIAL.md or ask your senior developer.

**Remember:** The goal is understanding, not speed. Take your time! ⏰