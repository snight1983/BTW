// Copyright (c) 2010 Satoshi Nakamoto
// Copyright (c) 2009-2017 The Bitcoin Core developers
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

#include <amount.h>
#include <base58.h>
#include <chain.h>
#include <consensus/validation.h>
#include <core_io.h>
#include <httpserver.h>
#include <validation.h>
#include <net.h>
#include <policy/feerate.h>
#include <policy/fees.h>
#include <policy/policy.h>
#include <policy/rbf.h>
#include <rpc/safemode.h>
#include <rpc/server.h>
#include <rpc/util.h>
#include <script/sign.h>
#include <timedata.h>
#include <util.h>
#include <utilmoneystr.h>
#include <init.h>
#include <stdint.h>
#include <univalue.h>
#include <rpc/blockchain.h>
#include <utilstrencodings.h>

UniValue blockToDeltasJSON(const CBlock& block, const CBlockIndex* blockindex)
{
	UniValue result(UniValue::VOBJ);
	result.push_back(Pair("hash", block.GetHash(false).GetHex()));
	int confirmations = -1;
	// Only report confirmations if the block is on the main chain
	if (chainActive.Contains(blockindex)) {
		confirmations = chainActive.Height() - blockindex->nHeight + 1;
	} else {
		throw JSONRPCError(RPC_INVALID_ADDRESS_OR_KEY, "Block is an orphan");
	}
	result.push_back(Pair("confirmations", confirmations));
	result.push_back(Pair("size", (int)::GetSerializeSize(block, SER_NETWORK, PROTOCOL_VERSION)));
	result.push_back(Pair("height", blockindex->nHeight));
	//result.push_back(Pair("version", block.nVersion));
	uint32_t luVer = block.nVersion;
	result.push_back(Pair("version",(uint64_t)luVer));
	result.push_back(Pair("merkleroot", block.hashMerkleRoot.GetHex()));

	UniValue deltas(UniValue::VARR);

	for (unsigned int i = 0; i < block.vtx.size(); i++) {
		const CTransaction &tx = *(block.vtx[i]);
		const uint256 txhash = tx.GetHash();

		UniValue entry(UniValue::VOBJ);
		entry.push_back(Pair("txid", txhash.GetHex()));
		entry.push_back(Pair("index", (int)i));

		UniValue inputs(UniValue::VARR);

		if (!tx.IsCoinBase()) {

			for (size_t j = 0; j < tx.vin.size(); j++) {
				const CTxIn input = tx.vin[j];

				UniValue delta(UniValue::VOBJ);

				CSpentIndexValue spentInfo;
				CSpentIndexKey spentKey(input.prevout.hash, input.prevout.n);

				if (GetSpentIndex(spentKey, spentInfo)) {
					if (spentInfo.addressType == 1) {
						delta.push_back(Pair("address", CBitcoinAddress(CKeyID(spentInfo.addressHash)).ToString()));
					} else if (spentInfo.addressType == 2)  {
						delta.push_back(Pair("address", CBitcoinAddress(CScriptID(spentInfo.addressHash)).ToString()));
					} else {
						continue;
					}
					delta.push_back(Pair("satoshis", -1 * spentInfo.satoshis));
					delta.push_back(Pair("index", (int)j));
					delta.push_back(Pair("prevtxid", input.prevout.hash.GetHex()));
					delta.push_back(Pair("prevout", (int)input.prevout.n));

					inputs.push_back(delta);
				} else {
					throw JSONRPCError(RPC_INTERNAL_ERROR, "Spent information not available");
				}

			}
		}

		entry.push_back(Pair("inputs", inputs));

		UniValue outputs(UniValue::VARR);

		for (unsigned int k = 0; k < tx.vout.size(); k++) {
			const CTxOut &out = tx.vout[k];

			UniValue delta(UniValue::VOBJ);

			if (out.scriptPubKey.IsPayToScriptHash()) {
				std::vector<unsigned char> hashBytes(out.scriptPubKey.begin()+2, out.scriptPubKey.begin()+22);
				delta.push_back(Pair("address", CBitcoinAddress(CScriptID(uint160(hashBytes))).ToString()));

			} else if (out.scriptPubKey.IsPayToPublicKeyHash()) {
				std::vector<unsigned char> hashBytes(out.scriptPubKey.begin()+3, out.scriptPubKey.begin()+23);
				delta.push_back(Pair("address", CBitcoinAddress(CKeyID(uint160(hashBytes))).ToString()));
			} else {
				continue;
			}

			delta.push_back(Pair("satoshis", out.nValue));
			delta.push_back(Pair("index", (int)k));

			outputs.push_back(delta);
		}

		entry.push_back(Pair("outputs", outputs));
		deltas.push_back(entry);

	}
	result.push_back(Pair("deltas", deltas));
	result.push_back(Pair("time", block.GetBlockTime()));
	result.push_back(Pair("mediantime", (int64_t)blockindex->GetMedianTimePast()));
	result.push_back(Pair("nonce", (uint64_t)block.nNonce));
	result.push_back(Pair("bits", strprintf("%08x", block.nBits)));
	result.push_back(Pair("difficulty", GetDifficulty(blockindex)));
	result.push_back(Pair("chainwork", blockindex->nChainWork.GetHex()));

	if (blockindex->pprev)
		result.push_back(Pair("previousblockhash", blockindex->pprev->GetBlockHash().GetHex()));
	CBlockIndex *pnext = chainActive.Next(blockindex);
	if (pnext)
		result.push_back(Pair("nextblockhash", pnext->GetBlockHash().GetHex()));
	return result;
}

bool getAddressFromIndex(const int &type, const uint160 &hash, std::string &address)
{
	if (type == 2) {
		address = CBitcoinAddress(CScriptID(hash)).ToString();
	} else if (type == 1) {
		address = CBitcoinAddress(CKeyID(hash)).ToString();
	} else {
		return false;
	}
	return true;
}

bool getAddressesFromParams(const UniValue& params, std::vector<std::pair<uint160, int> > &addresses)
{
	if (params[0].isStr()) {
		CBitcoinAddress address(params[0].get_str());
		uint160 hashBytes;
		int type = 0;
		if (!address.GetIndexKey(hashBytes, type)) {
			throw JSONRPCError(RPC_INVALID_ADDRESS_OR_KEY, "Invalid address");
		}
		addresses.push_back(std::make_pair(hashBytes, type));
	} else if (params[0].isObject()) {

		UniValue addressValues = find_value(params[0].get_obj(), "addresses");
		if (!addressValues.isArray()) {
			throw JSONRPCError(RPC_INVALID_ADDRESS_OR_KEY, "Addresses is expected to be an array");
		}

		std::vector<UniValue> values = addressValues.getValues();

		for (std::vector<UniValue>::iterator it = values.begin(); it != values.end(); ++it) {

			CBitcoinAddress address(it->get_str());
			uint160 hashBytes;
			int type = 0;
			if (!address.GetIndexKey(hashBytes, type)) {
				throw JSONRPCError(RPC_INVALID_ADDRESS_OR_KEY, "Invalid address");
			}
			addresses.push_back(std::make_pair(hashBytes, type));
		}
	} else {
		throw JSONRPCError(RPC_INVALID_ADDRESS_OR_KEY, "Invalid address");
	}

	return true;
}

bool heightSort(std::pair<CAddressUnspentKey, CAddressUnspentValue> a,
	std::pair<CAddressUnspentKey, CAddressUnspentValue> b) {
		return a.second.blockHeight < b.second.blockHeight;
}

bool timestampSort(std::pair<CMempoolAddressDeltaKey, CMempoolAddressDelta> a,
	std::pair<CMempoolAddressDeltaKey, CMempoolAddressDelta> b) {
		return a.second.time < b.second.time;
}


UniValue getblockhashes(const JSONRPCRequest& request)
{
	if (request.fHelp || request.params.size() < 2)
		throw std::runtime_error(
		"getblockhashes timestamp\n"
		"\nReturns array of hashes of blocks within the timestamp range provided.\n"
		"\nArguments:\n"
		"1. high         (numeric, required) The newer block timestamp\n"
		"2. low          (numeric, required) The older block timestamp\n"
		"3. options      (string, required) A json object\n"
		"    {\n"
		"      \"noOrphans\":true   (boolean) will only include blocks on the main chain\n"
		"      \"logicalTimes\":true   (boolean) will include logical timestamps with hashes\n"
		"    }\n"
		"\nResult:\n"
		"[\n"
		"  \"hash\"         (string) The block hash\n"
		"]\n"
		"[\n"
		"  {\n"
		"    \"blockhash\": (string) The block hash\n"
		"    \"logicalts\": (numeric) The logical timestamp\n"
		"  }\n"
		"]\n"
		"\nExamples:\n"
		+ HelpExampleCli("getblockhashes", "1231614698 1231024505")
		+ HelpExampleRpc("getblockhashes", "1231614698, 1231024505")
		+ HelpExampleCli("getblockhashes", "1231614698 1231024505 '{\"noOrphans\":false, \"logicalTimes\":true}'")
		);

	unsigned int high = request.params[0].get_int();
	unsigned int low = request.params[1].get_int();
	bool fActiveOnly = false;
	bool fLogicalTS = false;

	if (request.params.size() > 2) {
		if (request.params[2].isObject()) {
			UniValue noOrphans = find_value(request.params[2].get_obj(), "noOrphans");
			UniValue returnLogical = find_value(request.params[2].get_obj(), "logicalTimes");

			if (noOrphans.isBool())
				fActiveOnly = noOrphans.get_bool();

			if (returnLogical.isBool())
				fLogicalTS = returnLogical.get_bool();
		}
	}

	std::vector<std::pair<uint256, unsigned int> > blockHashes;

	if (fActiveOnly)
		LOCK(cs_main);

	if (!GetTimestampIndex(high, low, fActiveOnly, blockHashes)) {
		throw JSONRPCError(RPC_INVALID_ADDRESS_OR_KEY, "No information available for block hashes");
	}

	UniValue result(UniValue::VARR);

	for (std::vector<std::pair<uint256, unsigned int> >::const_iterator it=blockHashes.begin(); it!=blockHashes.end(); it++) {
		if (fLogicalTS) {
			UniValue item(UniValue::VOBJ);
			item.push_back(Pair("blockhash", it->first.GetHex()));
			item.push_back(Pair("logicalts", (int)it->second));
			result.push_back(item);
		} else {
			result.push_back(it->first.GetHex());
		}
	}

	return result;
}

UniValue getblockdeltas(const JSONRPCRequest& request)
{
	if (request.fHelp || request.params.size() != 1)
		throw std::runtime_error("");

	std::string strHash = request.params[0].get_str();
	uint256 hash(uint256S(strHash));

	if (mapBlockIndex.count(hash) == 0)
		throw JSONRPCError(RPC_INVALID_ADDRESS_OR_KEY, "Block not found");

	CBlock block;
	CBlockIndex* pblockindex = mapBlockIndex[hash];

	if (fHavePruned && !(pblockindex->nStatus & BLOCK_HAVE_DATA) && pblockindex->nTx > 0)
		throw JSONRPCError(RPC_INTERNAL_ERROR, "Block not available (pruned data)");

	if(!ReadBlockFromDisk(block, pblockindex, Params().GetConsensus()))
		throw JSONRPCError(RPC_INTERNAL_ERROR, "Can't read block from disk");

	return blockToDeltasJSON(block, pblockindex);
}
////////////////////////////
UniValue getaddressmempool(const JSONRPCRequest& request)
{
	if (request.fHelp || request.params.size() != 1)
		throw std::runtime_error(
		"getaddressmempool\n"
		"\nReturns all mempool deltas for an address (requires addressindex to be enabled).\n"
		"\nArguments:\n"
		"{\n"
		"  \"addresses\"\n"
		"    [\n"
		"      \"address\"  (string) The base58check encoded address\n"
		"      ,...\n"
		"    ]\n"
		"}\n"
		"\nResult:\n"
		"[\n"
		"  {\n"
		"    \"address\"  (string) The base58check encoded address\n"
		"    \"txid\"  (string) The related txid\n"
		"    \"index\"  (number) The related input or output index\n"
		"    \"satoshis\"  (number) The difference of satoshis\n"
		"    \"timestamp\"  (number) The time the transaction entered the mempool (seconds)\n"
		"    \"prevtxid\"  (string) The previous txid (if spending)\n"
		"    \"prevout\"  (string) The previous transaction output index (if spending)\n"
		"  }\n"
		"]\n"
		"\nExamples:\n"
		+ HelpExampleCli("getaddressmempool", "'{\"addresses\": [\"12c6DSiU4Rq3P4ZxziKxzrL5LmMBrzjrJX\"]}'")
		+ HelpExampleRpc("getaddressmempool", "{\"addresses\": [\"12c6DSiU4Rq3P4ZxziKxzrL5LmMBrzjrJX\"]}")
		);

	std::vector<std::pair<uint160, int> > addresses;

	if (!getAddressesFromParams(request.params, addresses)) {
		throw JSONRPCError(RPC_INVALID_ADDRESS_OR_KEY, "Invalid address");
	}

	std::vector<std::pair<CMempoolAddressDeltaKey, CMempoolAddressDelta> > indexes;

	if (!mempool.getAddressIndex(addresses, indexes)) {
		throw JSONRPCError(RPC_INVALID_ADDRESS_OR_KEY, "No information available for address");
	}

	std::sort(indexes.begin(), indexes.end(), timestampSort);

	UniValue result(UniValue::VARR);

	for (std::vector<std::pair<CMempoolAddressDeltaKey, CMempoolAddressDelta> >::iterator it = indexes.begin();
		it != indexes.end(); it++) {

			std::string address;
			if (!getAddressFromIndex(it->first.type, it->first.addressBytes, address)) {
				throw JSONRPCError(RPC_INVALID_ADDRESS_OR_KEY, "Unknown address type");
			}

			UniValue delta(UniValue::VOBJ);
			delta.push_back(Pair("address", address));
			delta.push_back(Pair("txid", it->first.txhash.GetHex()));
			delta.push_back(Pair("index", (int)it->first.index));
			delta.push_back(Pair("satoshis", it->second.amount));
			delta.push_back(Pair("timestamp", it->second.time));
			if (it->second.amount < 0) {
				delta.push_back(Pair("prevtxid", it->second.prevhash.GetHex()));
				delta.push_back(Pair("prevout", (int)it->second.prevout));
			}
			result.push_back(delta);
	}

	return result;
}

UniValue getaddressutxos(const JSONRPCRequest& request)
{
	if (request.fHelp || request.params.size() != 1)
		throw std::runtime_error(
		"getaddressutxos\n"
		"\nReturns all unspent outputs for an address (requires addressindex to be enabled).\n"
		"\nArguments:\n"
		"{\n"
		"  \"addresses\"\n"
		"    [\n"
		"      \"address\"  (string) The base58check encoded address\n"
		"      ,...\n"
		"    ],\n"
		"  \"chainInfo\"  (boolean) Include chain info with results\n"
		"}\n"
		"\nResult\n"
		"[\n"
		"  {\n"
		"    \"address\"  (string) The address base58check encoded\n"
		"    \"txid\"  (string) The output txid\n"
		"    \"height\"  (number) The block height\n"
		"    \"outputIndex\"  (number) The output index\n"
		"    \"script\"  (strin) The script hex encoded\n"
		"    \"satoshis\"  (number) The number of satoshis of the output\n"
		"  }\n"
		"]\n"
		"\nExamples:\n"
		+ HelpExampleCli("getaddressutxos", "'{\"addresses\": [\"12c6DSiU4Rq3P4ZxziKxzrL5LmMBrzjrJX\"]}'")
		+ HelpExampleRpc("getaddressutxos", "{\"addresses\": [\"12c6DSiU4Rq3P4ZxziKxzrL5LmMBrzjrJX\"]}")
		);

	bool includeChainInfo = false;
	if (request.params[0].isObject()) {
		UniValue chainInfo = find_value(request.params[0].get_obj(), "chainInfo");
		if (chainInfo.isBool()) {
			includeChainInfo = chainInfo.get_bool();
		}
	}

	std::vector<std::pair<uint160, int> > addresses;

	if (!getAddressesFromParams(request.params, addresses)) {
		throw JSONRPCError(RPC_INVALID_ADDRESS_OR_KEY, "Invalid address");
	}

	std::vector<std::pair<CAddressUnspentKey, CAddressUnspentValue> > unspentOutputs;

	for (std::vector<std::pair<uint160, int> >::iterator it = addresses.begin(); it != addresses.end(); it++) {
		if (!GetAddressUnspent((*it).first, (*it).second, unspentOutputs)) {
			throw JSONRPCError(RPC_INVALID_ADDRESS_OR_KEY, "No information available for address");
		}
	}

	std::sort(unspentOutputs.begin(), unspentOutputs.end(), heightSort);

	UniValue utxos(UniValue::VARR);

	for (std::vector<std::pair<CAddressUnspentKey, CAddressUnspentValue> >::const_iterator it=unspentOutputs.begin(); it!=unspentOutputs.end(); it++) {
		UniValue output(UniValue::VOBJ);
		std::string address;
		if (!getAddressFromIndex(it->first.type, it->first.hashBytes, address)) {
			throw JSONRPCError(RPC_INVALID_ADDRESS_OR_KEY, "Unknown address type");
		}

		output.push_back(Pair("address", address));
		output.push_back(Pair("txid", it->first.txhash.GetHex()));
		output.push_back(Pair("outputIndex", (int)it->first.index));
		output.push_back(Pair("script", HexStr(it->second.script.begin(), it->second.script.end())));
		output.push_back(Pair("satoshis", it->second.satoshis));
		output.push_back(Pair("height", it->second.blockHeight));
		utxos.push_back(output);
	}

	if (includeChainInfo) {
		UniValue result(UniValue::VOBJ);
		result.push_back(Pair("utxos", utxos));

		LOCK(cs_main);
		result.push_back(Pair("hash", chainActive.Tip()->GetBlockHash().GetHex()));
		result.push_back(Pair("height", (int)chainActive.Height()));
		return result;
	} else {
		return utxos;
	}
}

UniValue getaddressdeltas(const JSONRPCRequest& request)
{
	if (request.fHelp || request.params.size() != 1 || !request.params[0].isObject())
		throw std::runtime_error(
		"getaddressdeltas\n"
		"\nReturns all changes for an address (requires addressindex to be enabled).\n"
		"\nArguments:\n"
		"{\n"
		"  \"addresses\"\n"
		"    [\n"
		"      \"address\"  (string) The base58check encoded address\n"
		"      ,...\n"
		"    ]\n"
		"  \"start\" (number) The start block height\n"
		"  \"end\" (number) The end block height\n"
		"  \"chainInfo\" (boolean) Include chain info in results, only applies if start and end specified\n"
		"}\n"
		"\nResult:\n"
		"[\n"
		"  {\n"
		"    \"satoshis\"  (number) The difference of satoshis\n"
		"    \"txid\"  (string) The related txid\n"
		"    \"index\"  (number) The related input or output index\n"
		"    \"height\"  (number) The block height\n"
		"    \"address\"  (string) The base58check encoded address\n"
		"  }\n"
		"]\n"
		"\nExamples:\n"
		+ HelpExampleCli("getaddressdeltas", "'{\"addresses\": [\"12c6DSiU4Rq3P4ZxziKxzrL5LmMBrzjrJX\"]}'")
		+ HelpExampleRpc("getaddressdeltas", "{\"addresses\": [\"12c6DSiU4Rq3P4ZxziKxzrL5LmMBrzjrJX\"]}")
		);


	UniValue startValue = find_value(request.params[0].get_obj(), "start");
	UniValue endValue = find_value(request.params[0].get_obj(), "end");

	UniValue chainInfo = find_value(request.params[0].get_obj(), "chainInfo");
	bool includeChainInfo = false;
	if (chainInfo.isBool()) {
		includeChainInfo = chainInfo.get_bool();
	}

	int start = 0;
	int end = 0;

	if (startValue.isNum() && endValue.isNum()) {
		start = startValue.get_int();
		end = endValue.get_int();
		if (start <= 0 || end <= 0) {
			throw JSONRPCError(RPC_INVALID_ADDRESS_OR_KEY, "Start and end is expected to be greater than zero");
		}
		if (end < start) {
			throw JSONRPCError(RPC_INVALID_ADDRESS_OR_KEY, "End value is expected to be greater than start");
		}
	}

	std::vector<std::pair<uint160, int> > addresses;

	if (!getAddressesFromParams(request.params, addresses)) {
		throw JSONRPCError(RPC_INVALID_ADDRESS_OR_KEY, "Invalid address");
	}

	std::vector<std::pair<CAddressIndexKey, CAmount> > addressIndex;

	for (std::vector<std::pair<uint160, int> >::iterator it = addresses.begin(); it != addresses.end(); it++) {
		if (start > 0 && end > 0) {
			if (!GetAddressIndex((*it).first, (*it).second, addressIndex, start, end)) {
				throw JSONRPCError(RPC_INVALID_ADDRESS_OR_KEY, "No information available for address");
			}
		} else {
			if (!GetAddressIndex((*it).first, (*it).second, addressIndex)) {
				throw JSONRPCError(RPC_INVALID_ADDRESS_OR_KEY, "No information available for address");
			}
		}
	}

	UniValue deltas(UniValue::VARR);

	for (std::vector<std::pair<CAddressIndexKey, CAmount> >::const_iterator it=addressIndex.begin(); it!=addressIndex.end(); it++) {
		std::string address;
		if (!getAddressFromIndex(it->first.type, it->first.hashBytes, address)) {
			throw JSONRPCError(RPC_INVALID_ADDRESS_OR_KEY, "Unknown address type");
		}

		UniValue delta(UniValue::VOBJ);
		delta.push_back(Pair("satoshis", it->second));
		delta.push_back(Pair("txid", it->first.txhash.GetHex()));
		delta.push_back(Pair("index", (int)it->first.index));
		delta.push_back(Pair("blockindex", (int)it->first.txindex));
		delta.push_back(Pair("height", it->first.blockHeight));
		delta.push_back(Pair("address", address));
		deltas.push_back(delta);
	}

	UniValue result(UniValue::VOBJ);

	if (includeChainInfo && start > 0 && end > 0) {
		LOCK(cs_main);

		if (start > chainActive.Height() || end > chainActive.Height()) {
			throw JSONRPCError(RPC_INVALID_ADDRESS_OR_KEY, "Start or end is outside chain range");
		}

		CBlockIndex* startIndex = chainActive[start];
		CBlockIndex* endIndex = chainActive[end];

		UniValue startInfo(UniValue::VOBJ);
		UniValue endInfo(UniValue::VOBJ);

		startInfo.push_back(Pair("hash", startIndex->GetBlockHash().GetHex()));
		startInfo.push_back(Pair("height", start));

		endInfo.push_back(Pair("hash", endIndex->GetBlockHash().GetHex()));
		endInfo.push_back(Pair("height", end));

		result.push_back(Pair("deltas", deltas));
		result.push_back(Pair("start", startInfo));
		result.push_back(Pair("end", endInfo));

		return result;
	} else {
		return deltas;
	}
}

UniValue getaddressbalance(const JSONRPCRequest& request)
{
	if (request.fHelp || request.params.size() != 1)
		throw std::runtime_error(
		"getaddressbalance\n"
		"\nReturns the balance for an address(es) (requires addressindex to be enabled).\n"
		"\nArguments:\n"
		"{\n"
		"  \"addresses\"\n"
		"    [\n"
		"      \"address\"  (string) The base58check encoded address\n"
		"      ,...\n"
		"    ]\n"
		"}\n"
		"\nResult:\n"
		"{\n"
		"  \"balance\"  (string) The current balance in satoshis\n"
		"  \"received\"  (string) The total number of satoshis received (including change)\n"
		"}\n"
		"\nExamples:\n"
		+ HelpExampleCli("getaddressbalance", "'{\"addresses\": [\"12c6DSiU4Rq3P4ZxziKxzrL5LmMBrzjrJX\"]}'")
		+ HelpExampleRpc("getaddressbalance", "{\"addresses\": [\"12c6DSiU4Rq3P4ZxziKxzrL5LmMBrzjrJX\"]}")
		);

	std::vector<std::pair<uint160, int> > addresses;

	if (!getAddressesFromParams(request.params, addresses)) {
		throw JSONRPCError(RPC_INVALID_ADDRESS_OR_KEY, "Invalid address");
	}

	std::vector<std::pair<CAddressIndexKey, CAmount> > addressIndex;

	for (std::vector<std::pair<uint160, int> >::iterator it = addresses.begin(); it != addresses.end(); it++) {
		if (!GetAddressIndex((*it).first, (*it).second, addressIndex)) {
			throw JSONRPCError(RPC_INVALID_ADDRESS_OR_KEY, "No information available for address");
		}
	}

	CAmount balance = 0;
	CAmount received = 0;

	for (std::vector<std::pair<CAddressIndexKey, CAmount> >::const_iterator it=addressIndex.begin(); it!=addressIndex.end(); it++) {
		if (it->second > 0) {
			received += it->second;
		}
		balance += it->second;
	}

	UniValue result(UniValue::VOBJ);
	result.push_back(Pair("balance", balance));
	result.push_back(Pair("received", received));

	return result;

}

UniValue getaddresstxids(const JSONRPCRequest& request)
{
	if (request.fHelp || request.params.size() != 1)
		throw std::runtime_error(
		"getaddresstxids\n"
		"\nReturns the txids for an address(es) (requires addressindex to be enabled).\n"
		"\nArguments:\n"
		"{\n"
		"  \"addresses\"\n"
		"    [\n"
		"      \"address\"  (string) The base58check encoded address\n"
		"      ,...\n"
		"    ]\n"
		"  \"start\" (number) The start block height\n"
		"  \"end\" (number) The end block height\n"
		"}\n"
		"\nResult:\n"
		"[\n"
		"  \"transactionid\"  (string) The transaction id\n"
		"  ,...\n"
		"]\n"
		"\nExamples:\n"
		+ HelpExampleCli("getaddresstxids", "'{\"addresses\": [\"12c6DSiU4Rq3P4ZxziKxzrL5LmMBrzjrJX\"]}'")
		+ HelpExampleRpc("getaddresstxids", "{\"addresses\": [\"12c6DSiU4Rq3P4ZxziKxzrL5LmMBrzjrJX\"]}")
		);

	std::vector<std::pair<uint160, int> > addresses;

	if (!getAddressesFromParams(request.params, addresses)) {
		throw JSONRPCError(RPC_INVALID_ADDRESS_OR_KEY, "Invalid address");
	}

	int start = 0;
	int end = 0;
	if (request.params[0].isObject()) {
		UniValue startValue = find_value(request.params[0].get_obj(), "start");
		UniValue endValue = find_value(request.params[0].get_obj(), "end");
		if (startValue.isNum() && endValue.isNum()) {
			start = startValue.get_int();
			end = endValue.get_int();
		}
	}

	std::vector<std::pair<CAddressIndexKey, CAmount> > addressIndex;

	for (std::vector<std::pair<uint160, int> >::iterator it = addresses.begin(); it != addresses.end(); it++) {
		if (start > 0 && end > 0) {
			if (!GetAddressIndex((*it).first, (*it).second, addressIndex, start, end)) {
				throw JSONRPCError(RPC_INVALID_ADDRESS_OR_KEY, "No information available for address");
			}
		} else {
			if (!GetAddressIndex((*it).first, (*it).second, addressIndex)) {
				throw JSONRPCError(RPC_INVALID_ADDRESS_OR_KEY, "No information available for address");
			}
		}
	}

	std::set<std::pair<int, std::string> > txids;
	UniValue result(UniValue::VARR);

	for (std::vector<std::pair<CAddressIndexKey, CAmount> >::const_iterator it=addressIndex.begin(); it!=addressIndex.end(); it++) {
		int height = it->first.blockHeight;
		std::string txid = it->first.txhash.GetHex();

		if (addresses.size() > 1) {
			txids.insert(std::make_pair(height, txid));
		} else {
			if (txids.insert(std::make_pair(height, txid)).second) {
				result.push_back(txid);
			}
		}
	}

	if (addresses.size() > 1) {
		for (std::set<std::pair<int, std::string> >::const_iterator it=txids.begin(); it!=txids.end(); it++) {
			result.push_back(it->second);
		}
	}

	return result;

}

UniValue getspentinfo(const JSONRPCRequest& request)
{
	if (request.fHelp || request.params.size() != 1 || !request.params[0].isObject())
		throw std::runtime_error(
		"getspentinfo\n"
		"\nReturns the txid and index where an output is spent.\n"
		"\nArguments:\n"
		"{\n"
		"  \"txid\" (string) The hex string of the txid\n"
		"  \"index\" (number) The start block height\n"
		"}\n"
		"\nResult:\n"
		"{\n"
		"  \"txid\"  (string) The transaction id\n"
		"  \"index\"  (number) The spending input index\n"
		"  ,...\n"
		"}\n"
		"\nExamples:\n"
		+ HelpExampleCli("getspentinfo", "'{\"txid\": \"0437cd7f8525ceed2324359c2d0ba26006d92d856a9c20fa0241106ee5a597c9\", \"index\": 0}'")
		+ HelpExampleRpc("getspentinfo", "{\"txid\": \"0437cd7f8525ceed2324359c2d0ba26006d92d856a9c20fa0241106ee5a597c9\", \"index\": 0}")
		);

	UniValue txidValue = find_value(request.params[0].get_obj(), "txid");
	UniValue indexValue = find_value(request.params[0].get_obj(), "index");

	if (!txidValue.isStr() || !indexValue.isNum()) {
		throw JSONRPCError(RPC_INVALID_ADDRESS_OR_KEY, "Invalid txid or index");
	}

	uint256 txid = ParseHashV(txidValue, "txid");
	int outputIndex = indexValue.get_int();

	CSpentIndexKey key(txid, outputIndex);
	CSpentIndexValue value;

	if (!GetSpentIndex(key, value)) {
		throw JSONRPCError(RPC_INVALID_ADDRESS_OR_KEY, "Unable to get spent info");
	}

	UniValue obj(UniValue::VOBJ);
	obj.push_back(Pair("txid", value.txid.GetHex()));
	obj.push_back(Pair("index", (int)value.inputIndex));
	obj.push_back(Pair("height", value.blockHeight));
	return obj;
}

static const CRPCCommand commands[] =
{ //  category              name                        actor (function)           argNames
	//  --------------------- ------------------------    -----------------------  ----------
	{ "blockchain",         "getblockdeltas",         &getblockdeltas,         {} },
	{ "blockchain",         "getblockhashes",         &getblockhashes,         {} },
	{ "blockchain",         "getspentinfo",           &getspentinfo,           {"txid","index"} },
	/* Address index */
	{ "addressindex",       "getaddressmempool",      &getaddressmempool,      {"addresses"} },
	{ "addressindex",       "getaddressutxos",        &getaddressutxos,        {"addresses", "chainInfo"} },
	{ "addressindex",       "getaddressdeltas",       &getaddressdeltas,       {"addresses","start","end","chainInfo"} },
	{ "addressindex",       "getaddresstxids",        &getaddresstxids,        {"addresses"} },
	{ "addressindex",       "getaddressbalance",      &getaddressbalance,      {"addresses"} },
};

void RegisterAddressRPCCommands(CRPCTable &t)
{
	for (unsigned int vcidx = 0; vcidx < ARRAYLEN(commands); vcidx++)
		t.appendCommand(commands[vcidx].name, &commands[vcidx]);
}