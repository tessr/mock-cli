package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"
)

const name = "poetic-ocean"

const manifest = `
name: %s # auto-generated
status: online
deployed_on: %s

l1:
  name: sepolia
  chain_id: 11155111
l2:
  deployment:
    op-deployer:
      version: op-deployer/v0.0.11
    l1-contracts:
      locator: https://storage.googleapis.com/oplabs-contract-artifacts/artifacts-v1-c3f2e2adbd52a93c2c08cab018cd637a4e203db53034e59c6c139c76b4297953.tar.gz
      version: 984bae9146398a2997ec13757bfe2438ca8f92eb
    l2-contracts:
      version: op-contracts/v1.7.0-beta.1+l2-contracts
    overrides:
      seconds_per_slot: 2
      fjord_time_offset: 0
      granite_time_offset: 0
      holocene_time_offset: 0
  components:
    op-node:
      version: op-node/v1.10.2
    op-geth:
      version: op-geth/v1.101411.4-rc.4
    op-reth:
      version: op-reth/v1.1.5
    op-proposer:
      version: op-proposer/v1.10.0-rc.2
    op-batcher:
      version: op-batcher/v1.10.0
    op-challenger:
      version: op-challenger/v1.3.1-rc.4
  bootnodes:
    - enode://e63f16f93b44938a43ba8304d9948da3304e47c7a14fb1ce9d75194f2bd6784665aab18a8b38dc7772cd9a7fa71dd51997cc6cd8af0fc8915bc2fddcb68dfc87@34.123.30.27:0?discport=30305
    - enode://524958a8fc12140491f59a78442d38d1f3cb96765136379dca2ce569f0d770de25931fc11ad1bcfe1af3d4db025d23c25a401e23c91f5f53ff281f25bd36e733@34.27.92.24:0?discport=30305
    - enode://7353cf20a38a59c67139ccceaf4fa94be3487f0ace88ca368c24026f0825ff6983493a6b7c6e4286e18394880f3f450987ca1f26d848b7057e7dfdf3eb6b71a5@34.28.88.99:0?discport=30305
  chains:
  - name: %s-0
    chain_id: %d
`

const inventory = `
chains:
  - name: "%s-0"
    nodes:
      - kind: "node"
        name: "sequencer-0"
        spec:
          kind: "sequencer"
          el:
            kind: "op-geth"
            spec:
              kind: "full"
          cl:
            kind: "op-node"

      - kind: "node"
        name: "sequencer-1"
        spec:
          kind: "sequencer"
          el:
            kind: "op-geth"
            spec:
              kind: "full"
          cl:
            kind: "op-node"

      - kind: "node"
        name: "sequencer-2"
        spec:
          kind: "sequencer"
          el:
            kind: "op-reth"
            spec:
              kind: "full"
          cl:
            kind: "op-node"

      - kind: "node"
        name: "rpc-0"
        spec:
          kind: "rpc"
          el:
            kind: "op-geth"
            spec:
              kind: "archive"
          cl:
            kind: "op-node"

      - kind: "node"
        name: "rpc-1"
        spec:
          kind: "rpc"
          el:
            kind: "op-reth"
            spec:
              kind: "archive"
          cl:
            kind: "op-node"

      - kind: "node"
        name: "snapsync-0"
        spec:
          kind: "snapsync"
          el:
            kind: "op-geth"
            spec:
              kind: "full"
          cl:
            kind: "op-node"

      - kind: "node"
        name: "snapsync-1"
        spec:
          kind: "snapsync"
          el:
            kind: "op-reth"
            spec:
              kind: "full"
          cl:
            kind: "op-node"

    services:

      - kind: "proposer"
        name: "op-proposer"
        spec:
          kind: "op-proposer"
        deps:
          nodes:
            - "rpc-0"

      - kind: "batcher"
        name: "op-batcher"
        spec:
          kind: "op-batcher"
        deps:
          nodes:
            - "sequencer-0"
            - "sequencer-1"
            - "sequencer-2"

      - kind: "challenger"
        name: "op-challenger"
        spec:
          kind: "op-challenger"
        deps:
          nodes:
            - "rpc-0"

      - kind: "proxyd"
        name: "proxyd-public"
        spec:
          kind: "public"
        deps:
          nodes:
            - "rpc-0"
            - "rpc-1"

      - kind: "peer-mgmt-service"
        name: "peer-mgmt-service"
        spec:
          kind: "peer-mgmt-service"
        deps:
          nodes:
            - "sequencer-0"
            - "sequencer-1"
            - "sequencer-2"
            - "rpc-0"
            - "rpc-1"
            - "snapsync-0"
            - "snapsync-1"

      - kind: "op-conductor-mon"
        name: "op-conductor-mon"
        spec:
          kind: "op-conductor-mon"
        deps:
          nodes:
            - "sequencer-0"
            - "sequencer-1"
            - "sequencer-2"

      - kind: "op-dispute-mon"
        name: "op-dispute-mon"
        spec:
          kind: "op-dispute-mon"
        deps:
          nodes:
            - "rpc-0"
            - "rpc-1"
`

const state = `

{
	"version": 1,
	"create2Salt": "0x15daff19dec42a5ac45712b7563022efcac2030e98b8bc10cfab6c7ee57f484f",
	"appliedIntent": {
	  "configType": "standard-overrides",
	  "l1ChainID": 11155111,
	  "superchainRoles": {
		"proxyAdminOwner": "0x1eb2ffc903729a0f03966b917003800b145f56e2",
		"protocolVersionsOwner": "0x79add5713b383daa0a138d3c4780c7a1804a8090",
		"guardian": "0x7a50f00e8d05b95f98fe38d8bee366a7324dcf7e"
	  },
	  "fundDevAccounts": false,
	  "useInterop": false,
	  "l1ContractsLocator": "https://storage.googleapis.com/oplabs-contract-artifacts/artifacts-v1-c3f2e2adbd52a93c2c08cab018cd637a4e203db53034e59c6c139c76b4297953.tar.gz",
	  "l2ContractsLocator": "tag://op-contracts/v1.7.0-beta.1+l2-contracts",
	  "chains": [
		{
		  "id": "0x00000000000000000000000000000000000000000000000000000a25406f3e60",
		  "baseFeeVaultRecipient": "0x100f829718b5be38013cc7b29c5c62a08d00f1ff",
		  "l1FeeVaultRecipient": "0xbaeaf33e883068937ab4a50871f2fd52e241013a",
		  "sequencerFeeVaultRecipient": "0xd0d5d18f0ebb07b7d728b14aae014eeda814d6bd",
		  "eip1559DenominatorCanyon": 250,
		  "eip1559Denominator": 50,
		  "eip1559Elasticity": 6,
		  "roles": {
			"l1ProxyAdminOwner": "0xdf5a644aed1b5d6ce0da2add778bc5f39d97ac88",
			"l2ProxyAdminOwner": "0xc40445cd88dda2a410f86f6ef8e00fd52d8381fd",
			"systemConfigOwner": "0xb32296e6929f2507db8153a64b036d175ac6e89e",
			"unsafeBlockSigner": "0xa53526b516df4eee3791734ce85311569e0ead78",
			"batcher": "0x8680d36811420359093fd321ed386a6e76be2af3",
			"proposer": "0x41b3b204099771adf857f826015703a1030b6675",
			"challenger": "0x7b51a480daee699ca3a4f68f9aaa434452112ef7"
		  },
		  "deployOverrides": null,
		  "dangerousAltDAConfig": {
			"useAltDA": false,
			"daCommitmentType": "",
			"daChallengeWindow": 0,
			"daResolveWindow": 0,
			"daBondSize": 0,
			"daResolverRefundPercentage": 0
		  },
		  "dangerousAdditionalDisputeGames": null
		}
	  ],
	  "globalDeployOverrides": null
	},
	"superchainDeployment": {
	  "proxyAdminAddress": "0x4afeef31b14cd42945e15dbaf3464ccf93094653",
	  "superchainConfigProxyAddress": "0xa2e5be63f9d8b63a79cf910e9f69a115eff864c1",
	  "superchainConfigImplAddress": "0x755443d9a4b52648077972c61f9cb1aaa106f112",
	  "protocolVersionsProxyAddress": "0xc146bab396f3f826b536e911c026cea67b8d0970",
	  "protocolVersionsImplAddress": "0x6025148d8f5989b58e3a14a19eae5d64fb85cb04"
	},
	"implementationsDeployment": {
	  "opcmAddress": "0x02c4ede5296cc42909ca9ae821282ff5b5521705",
	  "delayedWETHImplAddress": "0xfe6aa167f0a92669954572565a7593a3b250f62c",
	  "optimismPortalImplAddress": "0x9f215ce8a399caeb629347c64666f88c5d8d7386",
	  "preimageOracleSingletonAddress": "0xae225b2accc35e8be115c585d1485a80859f96d4",
	  "mipsSingletonAddress": "0xc46ffa9de7ffe3cc551a5467bf9d5f30ac1799a2",
	  "systemConfigImplAddress": "0x356db1ba2736d43142382509a16f8f3b9422c559",
	  "l1CrossDomainMessengerImplAddress": "0x0a3fe7895d439493072c1670c971e2caf3f85c0a",
	  "l1ERC721BridgeImplAddress": "0x73461e861e9cca7de74ec103cd35ef2d14dbad8c",
	  "l1StandardBridgeImplAddress": "0xf708a07a1a7e6d70e2637e20cd0ccee50c82b89a",
	  "optimismMintableERC20FactoryImplAddress": "0xedb2637c3f195147f14784bbca4eca485ce5a680",
	  "disputeGameFactoryImplAddress": "0x7e62dfcf3348b738496d196306baa1b745ee5576"
	},
	"opChainDeployments": [
	  {
		"id": "0x00000000000000000000000000000000000000000000000000000a25406f3e60",
		"proxyAdminAddress": "0x42c19f797e51e2acbfcb2421aaf87f4b49b87843",
		"addressManagerAddress": "0xcc4cb98cbf7915b5b91348a3fe0c86b5358f036c",
		"l1ERC721BridgeProxyAddress": "0x9d5e7030332c720a29878f090589adf6269865e0",
		"systemConfigProxyAddress": "0x5ba3ab3ce0a91f7a6684b2b0e2246bf4f04bebc1",
		"optimismMintableERC20FactoryProxyAddress": "0xc75c6a9ac523265aac7e195b79002e0b7461e677",
		"l1StandardBridgeProxyAddress": "0xf5b780fb6ac553fe77531cd8a8e9708bf71501b5",
		"l1CrossDomainMessengerProxyAddress": "0xc456c4d07aceda9cfd8284e183c6bac4299eb7f3",
		"optimismPortalProxyAddress": "0xd796da8b8c116b293faf8769c84f0cd7eac9a5cb",
		"disputeGameFactoryProxyAddress": "0xd9801c8f2eae608942f3d7db02dc17869109f226",
		"anchorStateRegistryProxyAddress": "0x04dda72536f612e8ad4aec8e09126644e946b5de",
		"anchorStateRegistryImplAddress": "0x0000000000000000000000000000000000000000",
		"faultDisputeGameAddress": "0x0000000000000000000000000000000000000000",
		"permissionedDisputeGameAddress": "0xac739a7da8561b8f08d2b3f085e8c7576f72a3bd",
		"delayedWETHPermissionedGameProxyAddress": "0xfa40c8012e6f7668d227914c5d20a7ee127913af",
		"delayedWETHPermissionlessGameProxyAddress": "0x0000000000000000000000000000000000000000",
		"dataAvailabilityChallengeProxyAddress": "0x0000000000000000000000000000000000000000",
		"dataAvailabilityChallengeImplAddress": "0x0000000000000000000000000000000000000000",
		"additionalDisputeGames": null,
		"startBlock": {
		  "parentHash": "0xbfbf7e85c93e031b97fad589175d509631672c62f76c4b12280614cce4031ff9",
		  "sha3Uncles": "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
		  "miner": "0x3826539cbd8d68dcf119e80b994557b4278cec9f",
		  "stateRoot": "0x5d6473e3a88df722a3b15051e6b86494bca5efa5097d0136ed421a75337072cf",
		  "transactionsRoot": "0xd1877fd6c3c6b679e5efbe5b3f3851e164dbe305f3ff994a1505511239a33c04",
		  "receiptsRoot": "0x45884f147a6617f73e16b1c9e7e62743e94abb95dc47f47e19c4c7d1b4123769",
		  "logsBloom": "0x117000a49108570198948993a68808340864ad4c2207daa6181a915972aa0796a08e1dce18380019a882204105bf1ca8519831208aaa741056071ac45fb58d026322500a2704040f48b66b4b1208c03514c48a2bb0602a70fe42548cb131ba122dfe6a4eaa2253d18725e0655424eae0d87a0ec661440cc0034d9519208b0e00e260710399d26648e842025b7b388e151ca005213409169c2d15404a1c000d792e8c6ad6d088f09e49c9704208a040562042f66857d29482849593d23a0a685b1b54ad9246088a00a12b14b128e2a008caf61224a181c9729123d13a0c6061b390540038658c0005c9241311fc1b020591ab52fc4d50826a8582313548a85228",
		  "difficulty": "0x0",
		  "number": "0x727172",
		  "gasLimit": "0x2255100",
		  "gasUsed": "0xf6b306",
		  "timestamp": "0x67884564",
		  "extraData": "0x",
		  "mixHash": "0x402f3299bb65035f01ef0a5d094d3dccd784afe9f7c5112ebd93a1b72bd38ff3",
		  "nonce": "0x0000000000000000",
		  "baseFeePerGas": "0x4ce311",
		  "withdrawalsRoot": "0x77645d5d04666c1f0318365dc837c6c679c9ba9bb9d69f7ec0bb6576eb7a3acd",
		  "blobGasUsed": "0x80000",
		  "excessBlobGas": "0x20000",
		  "parentBeaconBlockRoot": "0x7018baa7f975635c98217fda911f3ea910d90c5dece4c4858cb7a28c2625a97a",
		  "requestsRoot": null,
		  "hash": "0xd84d7e6e3de812c7e0305d52971dc7488acaa2b2611ecc5e222e6bfc350d1940"
		}
	  }
	],
	"l1StateDump": null,
	"DeploymentCalldata": null
  }

`

const (
	cReset  = "\033[0m"
	cRed    = "\033[31m"
	cGreen  = "\033[32m"
	cYellow = "\033[33m"
	cBlue   = "\033[34m"
	cPurple = "\033[35m"
	cCyan   = "\033[36m"
	cGray   = "\033[37m"
	cWhite  = "\033[97m"
)

func main() {
	// log.SetFlags(0)

	if len(os.Args) < 2 {
		log.Println("usage: netchef generate [] []")
		os.Exit(1)
	}

	if len(os.Args) == 2 && os.Args[1] == "generate" {
		// generate a default network
		err := generate()
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
	}

	os.Exit(0)
}

func generate() error {
	chainID := rand.Intn(10000000000)
	log.Printf("🧑🏻‍🍳  Cooking up a devnet with %sdefault state%s and %srandomly generated chain name and ID%s... ", cGreen, cReset, cGreen, cReset)
	time.Sleep(time.Second)
	log.Printf("🍤  Generated %snetwork_name=%s %schain_id=%d... %s", cGreen, name, cBlue, chainID, cReset)
	time.Sleep(time.Second)

	// write manifest file
	log.Printf("🥄  Cooking up a manifest file...\n")
	b := []byte(fmt.Sprintf(manifest, name, time.Now().Format("2006-01-02"), name, chainID))
	err := os.WriteFile("manifest.yaml", b, 0644)
	if err != nil {
		return err
	}

	time.Sleep(time.Second)

	// write inventory file
	log.Println("🥄  Cooking up an inventory file with all default services... ")
	b = []byte(fmt.Sprintf(inventory, name))
	err = os.WriteFile("inventory.yaml", b, 0644)
	if err != nil {
		return err
	}
	time.Sleep(time.Second)

	// write state.json
	log.Println("🚀  Calling into op-deployer to generate L1 contracts and genesis state...")
	time.Sleep(2 * time.Second)
	log.Println("🥄  Writing the state.json with the genesis created by op-deployer... ")
	b = []byte(state)
	err = os.WriteFile("state.json", b, 0644)
	if err != nil {
		return err
	}
	time.Sleep(2 * time.Second)

	// launch nodes
	log.Printf("🥄  Cooking up an %sop-geth %ssequencer %snode...", cBlue, cGreen, cReset)
	time.Sleep(time.Second)
	log.Printf("🥄  Cooking up an %sop-reth %ssequencer %snode...", cBlue, cGreen, cReset)
	time.Sleep(time.Second)
	log.Printf("🥄  Cooking up an %sop-geth %srpc %snode...", cBlue, cGreen, cReset)
	time.Sleep(time.Second)
	log.Printf("🥄  Cooking up an %sop-reth %srpc %snode...", cBlue, cGreen, cReset)
	time.Sleep(2 * time.Second)

	// launch services
	log.Printf("🥄  Cooking up a %sproposer %sservice...", cBlue, cReset)
	time.Sleep(time.Second)
	log.Printf("🥄  Cooking up a %sbatcher %sservice...", cBlue, cReset)
	time.Sleep(time.Second)
	log.Printf("🥄  Cooking up a %schallenger %sservice...", cBlue, cReset)
	time.Sleep(time.Second)
	log.Printf("🥄  Cooking up a %sproxyd %sservice...", cBlue, cReset)
	time.Sleep(time.Second)
	log.Printf("🥄  Cooking up an %sop-conductor %sservice...", cBlue, cReset)
	time.Sleep(2 * time.Second)

	// buen apetito!
	log.Printf("🧑🏻‍🍳  Your network is now available at %shttps://%s-0.optimism.io%s", cBlue, name, cReset)
	log.Printf("🧑🏻‍🍳  Buen apetito!")

	return nil
}
