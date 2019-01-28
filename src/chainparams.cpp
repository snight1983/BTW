// Copyright (c) 2010 Satoshi Nakamoto
// Copyright (c) 2009-2017 The Bitcoin Core developers
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

#include <chainparams.h>
#include <consensus/merkle.h>

#include <tinyformat.h>
#include <util.h>
#include <utilstrencodings.h>

#include <assert.h>
#include <memory>

#include <chainparamsseeds.h>


static CBlock CreateGenesisBlock(const char* pszTimestamp, const CScript& genesisOutputScript, uint32_t nTime, uint32_t nNonce, uint32_t nNonceLock, uint32_t nNonceEx, uint32_t nBits, int32_t nVersion, const CAmount& genesisReward)
{

    CMutableTransaction txNew;
    txNew.nVersion = 1;
    txNew.vin.resize(1);
    txNew.vout.resize(1);
    txNew.vin[0].scriptSig = CScript() << 0 << CScriptNum(nNonceEx) << std::vector<unsigned char>((const unsigned char*)pszTimestamp, (const unsigned char*)pszTimestamp + strlen(pszTimestamp));
	txNew.vout[0].nValue = genesisReward;
    txNew.vout[0].scriptPubKey = genesisOutputScript;

    CBlock genesis;
    genesis.nTime			= nTime;
	genesis.nHeight			= 0;
    genesis.nBits			= nBits;
    genesis.nNonce			= nNonce;
	genesis.nNonceLock_btcv	= nNonceLock,
    genesis.nVersion		= nVersion;
	genesis.hashSeed_btcv.SetNull();
	memcpy( genesis.hashSeed_btcv.begin(), "I love you, there is no purpose.", strlen("I love you, there is no purpose."));
    genesis.vtx.push_back(MakeTransactionRef(std::move(txNew)));


    genesis.hashPrevBlock.SetNull();
    genesis.hashMerkleRoot = BlockMerkleRoot(genesis);
    return genesis;
}

/**
 * Build the genesis block. Note that the output of its generation
 * transaction cannot be spent since it did not originally exist in the
 * database.
 *
 * CBlock(hash=000000000019d6, ver=1, hashPrevBlock=00000000000000, hashMerkleRoot=4a5e1e, nTime=1231006505, nBits=1d00ffff, nNonce=2083236893, vtx=1)
 *   CTransaction(hash=4a5e1e, ver=1, vin.size=1, vout.size=1, nLockTime=0)
 *     CTxIn(COutPoint(000000, -1), coinbase 04ffff001d0104455468652054696d65732030332f4a616e2f32303039204368616e63656c6c6f72206f6e206272696e6b206f66207365636f6e64206261696c6f757420666f722062616e6b73)
 *     CTxOut(nValue=50.00000000, scriptPubKey=0x5F1DF16B2B704C8A578D0B)
 *   vMerkleTree: 4a5e1e
 */
static CBlock CreateGenesisBlock(uint32_t nTime, uint32_t nNonce, uint32_t nNonceLock, uint32_t nNonceEx, uint32_t nBits, int32_t nVersion, const CAmount& genesisReward)
{
    const char* pszTimestamp = "I love you, there is no purpose. Just love you. YanZi YuXuan. October 20, 2018.";
	const CScript genesisOutputScript = CScript() << ParseHex("04dedbf36965ec429477ecfa9b466c9ed7102f36a06dc6f7be2a62cf67cb59678059eb24947aa29b3e0eefa61b5c73e48259d970ac351c6d2a92813bf6744306c9") << OP_CHECKSIG;
	return CreateGenesisBlock(pszTimestamp, genesisOutputScript, nTime, nNonce, nNonceLock, nNonceEx, nBits, nVersion, genesisReward);
}

void CChainParams::UpdateVersionBitsParameters(Consensus::DeploymentPos d, int64_t nStartTime, int64_t nTimeout)
{
    consensus.vDeployments[d].nStartTime = nStartTime;
    consensus.vDeployments[d].nTimeout = nTimeout;
}

/**
 * Main network
 */
/**
 * What makes a good checkpoint block?
 * + Is surrounded by blocks with reasonable timestamps
 *   (no blocks before with a timestamp after, none after with
 *    timestamp before)
 * + Contains no strange transactions
 */

class CMainParams : public CChainParams {
public:
    CMainParams() {
        strNetworkID = "main";
        consensus.nSubsidyHalvingInterval = 210000;
        consensus.powLimit = uint256S("0000ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff");
		consensus.nPowTargetTimespan = 24 * 60 * 60;
		consensus.nPowTargetSpacing = 5 * 60;
		consensus.fPowAllowMinDifficultyBlocks = false;
        consensus.fPowNoRetargeting = false;
        consensus.nRuleChangeActivationThreshold = 274;
		consensus.nMinerConfirmationWindow = 288; // nPowTargetTimespan / nPowTargetSpacing

		consensus.vDeployments[Consensus::DEPLOYMENT_TESTDUMMY].bit			= 28;
		consensus.vDeployments[Consensus::DEPLOYMENT_TESTDUMMY].nStartTime	= Consensus::BIP9Deployment::ALWAYS_ACTIVE;
		consensus.vDeployments[Consensus::DEPLOYMENT_TESTDUMMY].nTimeout	= Consensus::BIP9Deployment::NO_TIMEOUT;

		consensus.vDeployments[Consensus::DEPLOYMENT_CSV].bit				= 0;
		consensus.vDeployments[Consensus::DEPLOYMENT_CSV].nStartTime		= Consensus::BIP9Deployment::ALWAYS_ACTIVE;
		consensus.vDeployments[Consensus::DEPLOYMENT_CSV].nTimeout			= Consensus::BIP9Deployment::NO_TIMEOUT;

		consensus.vDeployments[Consensus::DEPLOYMENT_SEGWIT].bit			= 1;
		consensus.vDeployments[Consensus::DEPLOYMENT_SEGWIT].nStartTime		= Consensus::BIP9Deployment::ALWAYS_ACTIVE;
		consensus.vDeployments[Consensus::DEPLOYMENT_SEGWIT].nTimeout		= Consensus::BIP9Deployment::NO_TIMEOUT;

// 		
// 		consensus.vDeployments[Consensus::DEPLOYMENT_TESTDUMMY].nStartTime = 1540036500;
// 		consensus.vDeployments[Consensus::DEPLOYMENT_TESTDUMMY].nTimeout = 1555804500;
//  
//          // Deployment of BIP68, BIP112, and BIP113.
// 		
// 		consensus.vDeployments[Consensus::DEPLOYMENT_CSV].nStartTime = 1540036500;
// 		consensus.vDeployments[Consensus::DEPLOYMENT_CSV].nTimeout = 1555804500;
//  
//          // Deployment of SegWit (BIP141, BIP143, and BIP147)
// 		
// 		consensus.vDeployments[Consensus::DEPLOYMENT_SEGWIT].nStartTime = 1540036500;
// 		consensus.vDeployments[Consensus::DEPLOYMENT_SEGWIT].nTimeout = 1555804500;

        // The best chain should have at least this much work.
		consensus.nMinimumChainWork = uint256S("0x000000000000000000000000000000000000000000000000000000021a4984fa");

        // By default assume that the signatures in ancestors of this block are valid.
        consensus.defaultAssumeValid = uint256S("0x00");

        /**
         * The message start string is designed to be unlikely to occur in normal data.
         * The characters are rarely used upper ASCII, not valid as UTF-8, and produce
         * a large 32-bit integer with any alignment.
         */
		pchMessageStart[0] = 0xe7;
		pchMessageStart[1] = 0xac;
		pchMessageStart[2] = 0xa2;
		pchMessageStart[3] = 0xc7;
		nDefaultPort = 8339;
		nPruneAfterHeight = 1;



		genesis = CreateGenesisBlock(1540036500, 7362240, 2100566, 417, 0x1f00ffff, 1, 50 * COIN * COIN_SCALE);
		consensus.hashGenesisBlock = genesis.GetHash();
		assert(consensus.hashGenesisBlock == uint256S("0x0000c7765285f588370f77e5cdd9171af52cdc7464921a181ef67db48b5a7a67"));
		assert(genesis.hashMerkleRoot == uint256S("0xbd2ba9c0210229e3b4d6e7025d1de6d3f990c77ca189f470c7dfd80d2b108e81"));


        // Note that of those which support the service bits prefix, most only support a subset of
        // possible options.
        // This is fine at runtime as we'll fall back to using them as a oneshot if they dont support the
        // service bits we want, but we should get them updated to support all service bits wanted by any
        // release ASAP to avoid it where possible.
        vSeeds.emplace_back("dnsseed1.bitvip.org"); 
        vSeeds.emplace_back("dnsseed2.bitvip.org"); 
        vSeeds.emplace_back("dnsseed3.bitvip.org"); 
        vSeeds.emplace_back("dnsseed4.bitvip.org"); 
        vSeeds.emplace_back("dnsseed5.bitvip.org"); 
        vSeeds.emplace_back("dnsseed6.bitvip.org");
		vSeeds.emplace_back("dnsseed7.bitvip.org");
		vSeeds.emplace_back("dnsseed8.bitvip.org");
		vSeeds.emplace_back("dnsseed9.bitvip.org");
		vSeeds.emplace_back("dnsseed10.bitvip.org");
        //vSeeds.emplace_back("seed.bitcoin.sipa.be"); // Pieter Wuille, only supports x1, x5, x9, and xd
        //vSeeds.emplace_back("dnsseed.bluematt.me"); // Matt Corallo, only supports x9
		base58Prefixes[PUBKEY_ADDRESS] = std::vector<unsigned char>(1,70); 
		base58Prefixes[SCRIPT_ADDRESS] = std::vector<unsigned char>(1,132);
        base58Prefixes[SECRET_KEY] =     std::vector<unsigned char>(1,128);
		base58Prefixes[B_PUBKEY_ADDRESS] = std::vector<unsigned char>(1,75); 
		base58Prefixes[B_SCRIPT_ADDRESS] = std::vector<unsigned char>(1,135); 
        base58Prefixes[EXT_PUBLIC_KEY] = {0x04, 0x88, 0xB2, 0x1E};
        base58Prefixes[EXT_SECRET_KEY] = {0x04, 0x88, 0xAD, 0xE4};


		bech32_hrp = "vip";

        vFixedSeeds = std::vector<SeedSpec6>(pnSeed6_main, pnSeed6_main + ARRAYLEN(pnSeed6_main));

        fDefaultConsistencyChecks = false;
        fRequireStandard = true;
        fMineBlocksOnDemand = false;

        checkpointData = {
            {
                { 13589, uint256S("0x0000057bdd743f64f9e00b917a36d99ccd5c12b18763992bfc2acc5de70e8fbf")},
            }
        };

		chainTxData = ChainTxData{
				0,
				0,
				0 
		};
    }
};

/**
 * Testnet (v3)
 */
class CTestNetParams : public CChainParams {
public:
    CTestNetParams() {
        strNetworkID = "test";
        consensus.nSubsidyHalvingInterval = 210000;
        consensus.powLimit = uint256S("0000ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff");
        consensus.nPowTargetTimespan = 14 * 24 * 60 * 60; // two weeks
        consensus.nPowTargetSpacing = 10 * 60;
        consensus.fPowAllowMinDifficultyBlocks = true;
        consensus.fPowNoRetargeting = false;
        consensus.nRuleChangeActivationThreshold = 1512; // 75% for testchains
        consensus.nMinerConfirmationWindow = 2016; // nPowTargetTimespan / nPowTargetSpacing
        consensus.vDeployments[Consensus::DEPLOYMENT_TESTDUMMY].bit = 28;
        consensus.vDeployments[Consensus::DEPLOYMENT_TESTDUMMY].nStartTime = 1540036500; 
        consensus.vDeployments[Consensus::DEPLOYMENT_TESTDUMMY].nTimeout = 1555804500;

        // Deployment of BIP68, BIP112, and BIP113.
        consensus.vDeployments[Consensus::DEPLOYMENT_CSV].bit = 0;
        consensus.vDeployments[Consensus::DEPLOYMENT_CSV].nStartTime = 1540036500;
        consensus.vDeployments[Consensus::DEPLOYMENT_CSV].nTimeout = 1555804500; 

        // Deployment of SegWit (BIP141, BIP143, and BIP147)
        consensus.vDeployments[Consensus::DEPLOYMENT_SEGWIT].bit = 1;
        consensus.vDeployments[Consensus::DEPLOYMENT_SEGWIT].nStartTime = 1540036500;
        consensus.vDeployments[Consensus::DEPLOYMENT_SEGWIT].nTimeout = 1555804500;

        // The best chain should have at least this much work.
        consensus.nMinimumChainWork = uint256S("0x00000000000000000000000000000000000000000000002830dab7f76dbb7d63");

        // By default assume that the signatures in ancestors of this block are valid.
        consensus.defaultAssumeValid = uint256S("0x0000000002e9e7b00e1f6dc5123a04aad68dd0f0968d8c7aa45f6640795c37b1"); //1135275

        pchMessageStart[0] = 0x0b;
        pchMessageStart[1] = 0x11;
        pchMessageStart[2] = 0x09;
        pchMessageStart[3] = 0x07;
        nDefaultPort = 18339;
        nPruneAfterHeight = 1000;


		genesis = CreateGenesisBlock(1540036500, 5535949, 306879, 1259, 0x1f00ffff, 1, 50 * COIN * COIN_SCALE);
		consensus.hashGenesisBlock = genesis.GetHash();
		assert(consensus.hashGenesisBlock == uint256S("0x0000b4f638971d518e0046f4cc315800474914211b1b735405eb035609d3fa30"));
		assert(genesis.hashMerkleRoot == uint256S("0x88248b090eed6be8eca12fcad971b9458eeede97497579d856066e48f2cbf1d7"));

        vFixedSeeds.clear();
        vSeeds.clear();
        // nodes with support for servicebits filtering should be at the top


		base58Prefixes[PUBKEY_ADDRESS] = std::vector<unsigned char>(1,70); 
		base58Prefixes[SCRIPT_ADDRESS] = std::vector<unsigned char>(1,132);
		base58Prefixes[SECRET_KEY] =     std::vector<unsigned char>(1,128);
        base58Prefixes[EXT_PUBLIC_KEY] = {0x04, 0x35, 0x87, 0xCF};
        base58Prefixes[EXT_SECRET_KEY] = {0x04, 0x35, 0x83, 0x94};

        bech32_hrp = "tb";

        vFixedSeeds = std::vector<SeedSpec6>(pnSeed6_test, pnSeed6_test + ARRAYLEN(pnSeed6_test));

        fDefaultConsistencyChecks = false;
        fRequireStandard = false;
        fMineBlocksOnDemand = false;


        checkpointData = {
            {

            }
        };

        chainTxData = ChainTxData{
            0,
            0,
            0
        };

    }
};

/**
 * Regression test
 */
class CRegTestParams : public CChainParams {
public:
    CRegTestParams() {
        strNetworkID = "regtest";
        consensus.nSubsidyHalvingInterval = 150;
        consensus.powLimit = uint256S("7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff");
        consensus.nPowTargetTimespan = 14 * 24 * 60 * 60; // two weeks
        consensus.nPowTargetSpacing = 10 * 60;
        consensus.fPowAllowMinDifficultyBlocks = true;
        consensus.fPowNoRetargeting = true;
        consensus.nRuleChangeActivationThreshold = 108; // 75% for testchains
        consensus.nMinerConfirmationWindow = 144; // Faster than normal for regtest (144 instead of 2016)
        consensus.vDeployments[Consensus::DEPLOYMENT_TESTDUMMY].bit = 28;
        consensus.vDeployments[Consensus::DEPLOYMENT_TESTDUMMY].nStartTime = 0;
        consensus.vDeployments[Consensus::DEPLOYMENT_TESTDUMMY].nTimeout = Consensus::BIP9Deployment::NO_TIMEOUT;
        consensus.vDeployments[Consensus::DEPLOYMENT_CSV].bit = 0;
        consensus.vDeployments[Consensus::DEPLOYMENT_CSV].nStartTime = 0;
        consensus.vDeployments[Consensus::DEPLOYMENT_CSV].nTimeout = Consensus::BIP9Deployment::NO_TIMEOUT;
        consensus.vDeployments[Consensus::DEPLOYMENT_SEGWIT].bit = 1;
        consensus.vDeployments[Consensus::DEPLOYMENT_SEGWIT].nStartTime = Consensus::BIP9Deployment::ALWAYS_ACTIVE;
        consensus.vDeployments[Consensus::DEPLOYMENT_SEGWIT].nTimeout = Consensus::BIP9Deployment::NO_TIMEOUT;

        // The best chain should have at least this much work.
        consensus.nMinimumChainWork = uint256S("0x00");

        // By default assume that the signatures in ancestors of this block are valid.
        consensus.defaultAssumeValid = uint256S("0x00");

        pchMessageStart[0] = 0xfa;
        pchMessageStart[1] = 0xbf;
        pchMessageStart[2] = 0xb5;
        pchMessageStart[3] = 0xda;
        nDefaultPort = 18444;
        nPruneAfterHeight = 1000;

		genesis = CreateGenesisBlock(1540036500, 7213126, 4764699, 2122, 0x1f00ffff, 1, 50 * COIN * COIN_SCALE);
		consensus.hashGenesisBlock = genesis.GetHash();
		assert(consensus.hashGenesisBlock == uint256S("0x0000982940b9080b280c61b96b706f29a6a45fe6b98538059d75a4e686825d98"));
		assert(genesis.hashMerkleRoot == uint256S("0x40b61291c5e59b82c1a5c5f1fe104c7ee1253dbebe11bd50cb455dfd28968178"));

        vFixedSeeds.clear(); //!< Regtest mode doesn't have any fixed seeds.
        vSeeds.clear();      //!< Regtest mode doesn't have any DNS seeds.

        fDefaultConsistencyChecks = true;
        fRequireStandard = false;
        fMineBlocksOnDemand = true;

        checkpointData = {
            {

            }
        };

        chainTxData = ChainTxData{
            0,
            0,
            0
        };

		base58Prefixes[PUBKEY_ADDRESS] = std::vector<unsigned char>(1,70); 
		base58Prefixes[SCRIPT_ADDRESS] = std::vector<unsigned char>(1,132);
		base58Prefixes[SECRET_KEY] =     std::vector<unsigned char>(1,128);
        base58Prefixes[EXT_PUBLIC_KEY] = {0x04, 0x35, 0x87, 0xCF};
        base58Prefixes[EXT_SECRET_KEY] = {0x04, 0x35, 0x83, 0x94};

        bech32_hrp = "bcrt";
    }
};

static std::unique_ptr<CChainParams> globalChainParams;

const CChainParams &Params() {
    assert(globalChainParams);
    return *globalChainParams;
}

std::unique_ptr<CChainParams> CreateChainParams(const std::string& chain)
{
    if (chain == CBaseChainParams::MAIN)
        return std::unique_ptr<CChainParams>(new CMainParams());
    else if (chain == CBaseChainParams::TESTNET)
        return std::unique_ptr<CChainParams>(new CTestNetParams());
    else if (chain == CBaseChainParams::REGTEST)
        return std::unique_ptr<CChainParams>(new CRegTestParams());
    throw std::runtime_error(strprintf("%s: Unknown chain %s.", __func__, chain));
}

void SelectParams(const std::string& network)
{
    SelectBaseParams(network);
    globalChainParams = CreateChainParams(network);
}

void UpdateVersionBitsParameters(Consensus::DeploymentPos d, int64_t nStartTime, int64_t nTimeout)
{
    globalChainParams->UpdateVersionBitsParameters(d, nStartTime, nTimeout);
}
